package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuzi/cdk-recharge-system/internal/db"
	"github.com/tuzi/cdk-recharge-system/internal/pro20x"
)

// 自建码（local）兑换链路：preview → preflight → redeem → result。
//
// 会话状态全部落在 pro20x_orders 表：
//   preview   → 插入一条 status=preview
//   preflight → 校验 session 后 status=preflight
//   redeem    → 提交 Pro20x 任务后 status=running，写 task_id
//   result    → 轮询任务并回写 status/stage
//
// 前端不需要改：返回结构与 zovo 的公开兑换接口保持同形。

const (
	localStatusPreview   = "preview"
	localStatusPreflight = "preflight"
	localStatusRunning   = "running"
)

// localCDKEnabled 自建码通道是否可用（开关 + API Key）。
func localCDKEnabled() bool {
	if !pro20x.Enabled() {
		return false
	}
	return strings.TrimSpace(pro20x.LoadConfig().APIKey) != ""
}

// isLocalCDKCode 判断卡密是否为本站自建码。
func isLocalCDKCode(code string) bool {
	if strings.TrimSpace(code) == "" {
		return false
	}
	_, ok, _ := db.LookupLocalCDK(code)
	return ok
}

// isLocalToken 判断兑换令牌是否属于自建码会话。
func isLocalToken(token string) bool {
	if strings.TrimSpace(token) == "" {
		return false
	}
	o, err := db.GetPro20xOrderByToken(token)
	return err == nil && o != nil
}

func newRedeemToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("RT-%d", time.Now().UnixNano())
	}
	return "RT-" + hex.EncodeToString(b)
}

// localPaymentSource 读自建码通道的支付方式与卡台 ID。
func localPaymentSource() (source, platformID string) {
	source, _ = db.GetSetting("pro20x_payment_source")
	source = strings.TrimSpace(source)
	if source == "" {
		source = "auto_card"
	}
	platformID, _ = db.GetSetting("pro20x_card_platform_id")
	return source, strings.TrimSpace(platformID)
}

// handleLocalPreview 自建码预览。
func handleLocalPreview(c *gin.Context, code string) {
	if !localCDKEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "自建卡密通道未启用"})
		return
	}
	cd, ok, err := db.LookupLocalCDK(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询卡密失败"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "卡密不存在"})
		return
	}
	switch strings.ToLower(strings.TrimSpace(cd.Status)) {
	case "unused", "":
		// 可兑换
	case "consumed":
		c.JSON(http.StatusBadRequest, gin.H{"error": "卡密已使用"})
		return
	case "disabled":
		c.JSON(http.StatusBadRequest, gin.H{"error": "卡密已禁用"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "卡密状态异常: " + cd.Status})
		return
	}
	if cd.ExpiresAt != "" {
		if exp, err := time.Parse("2006-01-02 15:04:05", cd.ExpiresAt); err == nil && time.Now().After(exp) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "卡密已过期"})
			return
		}
	}
	if _, mapped := pro20x.MapPlan(cd.Plan); !mapped {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该套餐暂不支持自建码通道: " + cd.Plan})
		return
	}

	// 复用已有未消费的预览会话，避免重复点预览堆积记录
	if old, err := db.GetPro20xOrderByCDK(code); err == nil && old != nil &&
		(old.Status == localStatusPreview || old.Status == localStatusPreflight) {
		c.JSON(http.StatusOK, gin.H{
			"redemption_token": old.RedeemToken,
			"plan":             cd.Plan,
			"payment_country":  cd.Country,
			"credential_mode":  "session",
			"channel":          "local",
		})
		return
	}

	token := newRedeemToken()
	_, err = db.InsertPro20xOrder(db.Pro20xOrder{
		RedeemToken:   token,
		CDKCode:       cd.Code,
		Plan:          cd.Plan,
		Country:       cd.Country,
		ClientOrderID: "FB-" + strings.TrimPrefix(token, "RT-"),
		Status:        localStatusPreview,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建兑换会话失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"redemption_token": token,
		"plan":             cd.Plan,
		"payment_country":  cd.Country,
		"credential_mode":  "session",
		"channel":          "local",
	})
}

// handleLocalPreflight 自建码预检：校验 session、确认账号可购买。
func handleLocalPreflight(c *gin.Context, body map[string]any) {
	if !localCDKEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "自建卡密通道未启用"})
		return
	}
	token := str(body["redemption_token"])
	order, err := db.GetPro20xOrderByToken(token)
	if err != nil || order == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "兑换会话不存在，请重新预览卡密"})
		return
	}
	session := extractCredentialSession(body["credential"])
	if session == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供完整 Session JSON（含 sessionToken）"})
		return
	}

	cli := pro20x.NewFromSettings()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 90*time.Second)
	defer cancel()
	res, err := cli.AccountCheck(ctx, session)
	if err != nil {
		log.Printf("[local-cdk] account-check failed order=%d: %v", order.ID, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "账号检测失败: " + err.Error()})
		return
	}
	var acc *pro20x.AccountInfo
	for i := range res.Accounts {
		if res.Accounts[i].AccountID == res.SelectedAccountID {
			acc = &res.Accounts[i]
			break
		}
	}
	if acc == nil && len(res.Accounts) > 0 {
		acc = &res.Accounts[0]
	}
	if acc == nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "未识别到账号信息，请重新导出 Session"})
		return
	}
	if acc.HasActiveSubscription {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("账号已有有效订阅（%s），无需重复开通", acc.PlanLabel),
			"code":  "already_subscribed",
		})
		return
	}

	// 存 session 供后续兑换/对账使用
	_ = db.BindCDKSession(order.CDKCode, token, session)

	order.Status = localStatusPreflight
	order.AccountEmail = res.Email
	_ = db.UpdatePro20xOrder(*order)

	c.JSON(http.StatusOK, gin.H{
		"preflight_token": token,
		"account": gin.H{
			"plan_type":               strings.ToLower(acc.Tier),
			"plan_label":              acc.PlanLabel,
			"has_active_subscription": acc.HasActiveSubscription,
			"billing_currency":        acc.BillingCurrency,
			"email":                   res.Email,
		},
		"channel": "local",
	})
}

// handleLocalRedeem 自建码兑换：提交 Pro20x 任务。
func handleLocalRedeem(c *gin.Context, body map[string]any) {
	if !localCDKEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "自建卡密通道未启用"})
		return
	}
	token := str(body["redemption_token"])
	preflightToken := str(body["preflight_token"])
	if token == "" {
		token = preflightToken
	}
	order, err := db.GetPro20xOrderByToken(token)
	if err != nil || order == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "兑换会话不存在，请重新预览卡密"})
		return
	}
	if order.Status != localStatusPreflight {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先完成账号检测后再兑换"})
		return
	}
	cd, ok, _ := db.LookupLocalCDK(order.CDKCode)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "卡密不存在"})
		return
	}
	if strings.ToLower(cd.Status) == "consumed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "卡密已使用"})
		return
	}
	session, _ := db.GetSessionByCDK(order.CDKCode)
	if strings.TrimSpace(session) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话已失效，请重新预检"})
		return
	}

	source, platformID := localPaymentSource()
	req, err := pro20x.BuildCreateRequest(pro20x.BuildOptions{
		ClientOrderID:   order.ClientOrderID,
		CDKPlan:         order.Plan,
		Country:         order.Country,
		FallbackCountry: pro20x.DefaultCountry(),
		AuthSession:     session,
		PaymentSource:   source,
		CardPlatformID:  platformID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cli := pro20x.NewFromSettings()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()
	task, err := cli.CreateTask(ctx, *req)
	if err != nil {
		log.Printf("[local-cdk] create task failed order=%d: %v", order.ID, err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "提交充值任务失败: " + err.Error()})
		return
	}

	order.TaskID = task.ID
	order.Status = localStatusRunning
	order.Stage = task.Status
	_ = db.UpdatePro20xOrder(*order)
	_ = db.MarkLocalCDKUsed(order.CDKCode, order.AccountEmail)

	c.JSON(http.StatusOK, gin.H{
		"id":      task.ID,
		"status":  task.Status,
		"stage":   task.Status,
		"channel": "local",
	})
}

// handleLocalResult 自建码进度查询。
func handleLocalResult(c *gin.Context, token string) {
	order, err := db.GetPro20xOrderByToken(token)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "兑换记录不存在"})
		return
	}
	status, stage := mapLocalStatus(order.Status, order.Stage)
	events := []gin.H{
		{"step": "accept", "message": "任务已提交", "created_at": order.CreatedAt},
	}
	if order.TaskID != "" {
		if st, sg, ok := refreshLocalOrder(c, order); ok {
			status, stage = st, sg
		}
		events = append(events, gin.H{"step": "pay", "message": "支付处理中", "created_at": order.UpdatedAt})
	}
	if status == "completed" {
		events = append(events, gin.H{"step": "done", "message": "开通成功", "created_at": order.UpdatedAt})
	}
	c.JSON(http.StatusOK, gin.H{
		"order": gin.H{
			"status":        status,
			"stage":         stage,
			"account_email": order.AccountEmail,
			"task_id":       order.TaskID,
			"channel":       "local",
		},
		"events": events,
	})
}

// refreshLocalOrder 主动拉一次任务状态并回写；返回映射后的 status/stage。
func refreshLocalOrder(c *gin.Context, order *db.Pro20xOrder) (string, string, bool) {
	cli := pro20x.NewFromSettings()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()
	task, err := cli.GetTask(ctx, order.TaskID)
	if err != nil {
		return "", "", false
	}
	order.Status = task.Status
	order.Stage = task.Stage
	order.ErrorCode = task.ErrorCode
	order.CleanupStatus = extractCleanupStatus(task.Cleanup)
	_ = db.UpdatePro20xOrder(*order)
	st, sg := mapLocalStatus(task.Status, task.Stage)
	return st, sg, true
}

func extractCleanupStatus(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var m struct {
		Status string `json:"status"`
	}
	if json.Unmarshal(raw, &m) == nil {
		return strings.TrimSpace(m.Status)
	}
	return ""
}

// mapLocalStatus 把 Pro20x 任务状态映射成前端认识的 status/stage。
func mapLocalStatus(taskStatus, taskStage string) (string, string) {
	switch strings.ToLower(strings.TrimSpace(taskStatus)) {
	case "queued":
		return "running", "await_card"
	case "running":
		return "running", "payment"
	case "captcha_required":
		return "review", "captcha"
	case "succeeded", "success":
		return "completed", "done"
	case "failed":
		return "failed", "failed"
	case "manual_review":
		return "review", "review"
	case "cancelled":
		return "cancelled", "cancelled"
	case localStatusPreview, localStatusPreflight:
		return "pending", "preflight"
	}
	if taskStage != "" {
		return "running", taskStage
	}
	return "pending", ""
}
