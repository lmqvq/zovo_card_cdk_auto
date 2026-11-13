package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuzi/cdk-recharge-system/internal/db"
	"github.com/tuzi/cdk-recharge-system/internal/pro20x"
)

// AdminGenerateLocalCDKs POST /api/v1/admin/cdk/local/generate
// 生成自建卡密（走 Pro20x 通道）。完整码只在本次响应返回一次。
func AdminGenerateLocalCDKs(c *gin.Context) {
	var body struct {
		Plan       string `json:"plan"`
		Count      int    `json:"count"`
		Country    string `json:"country"`
		ValidDays  int    `json:"valid_days"`
		PriceMinor int64  `json:"price_minor"`
		AgentID    string `json:"agent_id"`
		Batch      string `json:"batch"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	plan := strings.ToLower(strings.TrimSpace(body.Plan))
	if plan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan required"})
		return
	}
	if _, ok := pro20x.MapPlan(plan); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该套餐未映射到 Pro20x 通道: " + plan})
		return
	}
	if body.Count <= 0 || body.Count > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "count 需在 1~500 之间"})
		return
	}
	country := pro20x.NormalizeCountry(body.Country, pro20x.DefaultCountry())
	batch := strings.TrimSpace(body.Batch)
	if batch == "" {
		batch = "B" + nowCompact()
	}

	codes, err := db.GenerateLocalCDKs(plan, country, batch, strings.TrimSpace(body.AgentID),
		body.Count, body.PriceMinor, body.ValidDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u, _ := c.Get("username")
	username, _ := u.(string)
	db.WriteAudit(username, "local_cdk_generate",
		"batch="+batch+" plan="+plan+" count="+strconv.Itoa(len(codes)), c.ClientIP())

	out := make([]gin.H, 0, len(codes))
	for _, cd := range codes {
		out = append(out, gin.H{
			"code":        cd.Code,
			"plan":        cd.Plan,
			"country":     cd.Country,
			"batch_no":    cd.BatchNo,
			"agent_id":    cd.AgentID,
			"price_minor": cd.PriceMinor,
			"expires_at":  cd.ExpiresAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"batch":  batch,
		"count":  len(codes),
		"codes":  out,
		"notice": "完整卡密只显示这一次，请立即保存或导出",
	})
}

// AdminListLocalCDKs GET /api/v1/admin/cdk/local/list?limit=&status=&batch=
func AdminListLocalCDKs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := db.ListLocalCDKs(limit, c.Query("status"), c.Query("batch"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	total, _ := db.CountLocalCDKs("")
	unused, _ := db.CountLocalCDKs("unused")
	out := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		out = append(out, gin.H{
			"code":        r.Code,
			"code_prefix": r.Prefix,
			"plan":        r.Plan,
			"country":     r.Country,
			"status":      r.Status,
			"batch_no":    r.BatchNo,
			"agent_id":    r.AgentID,
			"price_minor": r.PriceMinor,
			"bound_email": r.BoundEmail,
			"expires_at":  r.ExpiresAt,
			"created_at":  r.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"list":   out,
		"total":  total,
		"unused": unused,
		"source": "local",
	})
}

// maskCode 列表只回前 7 位 + 后 4 位，避免完整码被随手截屏。
func maskCode(code string) string {
	if len(code) <= 12 {
		return code
	}
	return code[:7] + "****" + code[len(code)-4:]
}

// AdminUpdateLocalCDKStatus PUT /api/v1/admin/cdk/local/status
// body: { code, status }  status: unused / disabled / consumed
func AdminUpdateLocalCDKStatus(c *gin.Context) {
	var body struct {
		Code   string `json:"code"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Code) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code required"})
		return
	}
	st := strings.ToLower(strings.TrimSpace(body.Status))
	switch st {
	case "unused", "disabled", "consumed":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "status 只支持 unused / disabled / consumed"})
		return
	}
	if err := db.UpdateLocalCDKStatus(body.Code, st); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	u, _ := c.Get("username")
	username, _ := u.(string)
	db.WriteAudit(username, "local_cdk_status", "code="+maskCode(body.Code)+" status="+st, c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"ok": true, "code": body.Code, "status": st})
}

func nowCompact() string {
	return strings.ReplaceAll(strings.ReplaceAll(
		time.Now().Format("2006-01-02 15:04:05"), "-", ""), ":", "")
}
