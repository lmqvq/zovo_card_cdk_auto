package handler

import (
	"net/http"
	"strconv"
	"strings"

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
	u, _ := c.Get("username")
	username, _ := u.(string)
	db.WriteAudit(username, "card_delete_settings", "updated", c.ClientIP())
	AdminGetCardDeleteSettings(c)
}
