package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuzi/cdk-recharge-system/internal/db"
	"github.com/tuzi/cdk-recharge-system/internal/gptcard"
)

// 允许通过后台读写的解绑卡配置项（白名单，避免任意写 site_settings）。
var cardDeleteSettingKeys = []string{
	"card_delete_enabled",
	"card_delete_delay_seconds",
	"card_delete_renew_check_attempts",
	"card_delete_max_attempts",
	"card_delete_rate_per_minute",
	"card_delete_concurrency",
	"card_delete_setdefault_method",
	"card_delete_setdefault_path",
	"card_delete_setdefault_body",
	"card_delete_since",
}

// AdminListCardDeleteJobs GET /api/v1/admin/card-delete/jobs?limit=&action=
func AdminListCardDeleteJobs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	jobs, err := db.ListCardDeleteJobs(limit, c.Query("action"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]gin.H, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, gin.H{
			"id":                 j.ID,
			"cdk_id":             j.CDKID,
			"order_id":           j.OrderID,
			"plan":               j.Plan,
			"account_email":      j.AccountEmail,
			"card_last_four":     j.CardLastFour,
			"pm_id":              j.PMID,
			"pm_brand":           j.PMBrand,
			"cards_before":       j.CardsBefore,
			"action":             j.Action,
			"reason":             j.Reason,
			"http_status":        j.HTTPStatus,
			"attempts":           j.Attempts,
			"setdefault_used":    j.SetDefaultUsed,
			"session_state":      j.SessionState,
			"session_checked_at": j.SessionCheckedAt,
			"next_run_at":        j.NextRunAt,
			"created_at":         j.CreatedAt,
			"updated_at":         j.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"jobs": out, "enabled": gptcard.Enabled()})
}

// AdminRetryCardDeleteJob POST /api/v1/admin/card-delete/jobs/:id/retry
func AdminRetryCardDeleteJob(c *gin.Context) {
	id, err := strconv.ParseInt(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := gptcard.RetryJob(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, _ := c.Get("username")
	username, _ := u.(string)
	db.WriteAudit(username, "card_delete_retry", "job="+strconv.FormatInt(id, 10), c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"ok": true, "id": id})
}

// AdminCheckCardDeleteSession POST /api/v1/admin/card-delete/session-check
// body: { cdk_id }
func AdminCheckCardDeleteSession(c *gin.Context) {
	var body struct {
		CDKID int64 `json:"cdk_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.CDKID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cdk_id required"})
		return
	}
	state, err := gptcard.CheckSession(body.CDKID)
	if err != nil && state == "unknown" {
		c.JSON(http.StatusOK, gin.H{"state": "unknown", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"state": state})
}

// AdminCardDeleteMaterial GET /api/v1/admin/card-delete/material?cdk_id=&reveal=
// [CD-3] 取出该卡密绑定的会话材料：默认只返回摘要，reveal=1 才返回原文并写审计。
func AdminCardDeleteMaterial(c *gin.Context) {
	id, err := strconv.ParseInt(strings.TrimSpace(c.Query("cdk_id")), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cdk_id required"})
		return
	}
	raw, err := gptcard.SessionMaterial(id)
	if err != nil || strings.TrimSpace(raw) == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "该卡密未绑定 Session"})
		return
	}
	email, plan, accountID := sessionMaterialView(raw)
	resp := gin.H{
		"has_material": true,
		"bytes":        len(raw),
		"email":        maskSessionEmail(email),
		"plan":         plan,
		"account_id":   accountID,
	}
	if c.Query("reveal") == "1" {
		u, _ := c.Get("username")
		username, _ := u.(string)
		db.WriteAudit(username, "card_delete_material",
			fmt.Sprintf("cdk_id=%d bytes=%d", id, len(raw)), c.ClientIP())
		resp["material"] = raw
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, resp)
}

// sessionMaterialView 从 session JSON 中抽展示字段（不校验、不解密）。
func sessionMaterialView(raw string) (email, plan, accountID string) {
	var m map[string]any
	if json.Unmarshal([]byte(raw), &m) != nil {
		return "", "", ""
	}
	if u, ok := m["user"].(map[string]any); ok {
		email, _ = u["email"].(string)
	}
	if a, ok := m["account"].(map[string]any); ok {
		plan, _ = a["planType"].(string)
		accountID, _ = a["id"].(string)
	}
	return
}

// maskSessionEmail 脱敏：ab***@domain
func maskSessionEmail(e string) string {
	i := strings.IndexByte(e, '@')
	if i <= 2 {
		return e
	}
	return e[:2] + "***" + e[i:]
}

// AdminGetCardDeleteSettings GET /api/v1/admin/card-delete/settings
func AdminGetCardDeleteSettings(c *gin.Context) {
	out := gin.H{}
	for _, k := range cardDeleteSettingKeys {
		v, _ := db.GetSetting(k)
		out[k] = v
	}
	out["card_delete_setdefault_default"] = gin.H{
		"method": "POST",
		"path":   "/backend-api/payments/payment_method/default",
		"body":   `{"account_id":"{account_id}","payment_method_id":"{pm_id}"}`,
	}
	c.JSON(http.StatusOK, out)
}

// AdminPutCardDeleteSettings PUT /api/v1/admin/card-delete/settings
func AdminPutCardDeleteSettings(c *gin.Context) {
	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	for _, k := range cardDeleteSettingKeys {
		v, ok := body[k]
		if !ok {
			continue
		}
		if err := db.SetSetting(k, strings.TrimSpace(v)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	// 首次开启时记录生效时间：只处理此后的新订单，不回补历史单（session 多半已过期）
	if en, _ := db.GetSetting("card_delete_enabled"); strings.TrimSpace(en) == "1" {
		if cur, _ := db.GetSetting("card_delete_since"); strings.TrimSpace(cur) == "" {
			_ = db.SetSetting("card_delete_since", time.Now().UTC().Format("2006-01-02 15:04:05"))
		}
	}
	u, _ := c.Get("username")
	username, _ := u.(string)
	db.WriteAudit(username, "card_delete_settings", "updated", c.ClientIP())
	AdminGetCardDeleteSettings(c)
}
