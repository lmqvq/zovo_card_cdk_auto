package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tuzi/cdk-recharge-system/internal/db"
	"github.com/tuzi/cdk-recharge-system/internal/pro20x"
)

// Pro20x 通道配置键（后台可改，改完无需发版）。
var pro20xSettingKeys = []string{
	"pro20x_enabled",
	"pro20x_api_base",
	"pro20x_api_key",
	"pro20x_payment_source",
	"pro20x_card_platform_id",
	"pro20x_default_country",
}

// AdminGetPro20xSettings GET /api/v1/admin/pro20x/settings
func AdminGetPro20xSettings(c *gin.Context) {
	out := gin.H{}
	for _, k := range pro20xSettingKeys {
		v, _ := db.GetSetting(k)
		if k == "pro20x_api_key" {
			out[k+"_configured"] = strings.TrimSpace(v) != ""
			out[k+"_hint"] = maskSecret(v)
			continue
		}
		out[k] = v
	}
	out["channel_ready"] = pro20x.Enabled() && strings.TrimSpace(pro20x.LoadConfig().APIKey) != ""
	c.JSON(http.StatusOK, out)
}

// AdminPutPro20xSettings PUT /api/v1/admin/pro20x/settings
func AdminPutPro20xSettings(c *gin.Context) {
	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	for _, k := range pro20xSettingKeys {
		v, ok := body[k]
		if !ok {
			continue
		}
		// API Key 传空表示不修改（避免误清空）
		if k == "pro20x_api_key" && strings.TrimSpace(v) == "" {
			continue
		}
		if err := db.SetSetting(k, strings.TrimSpace(v)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	u, _ := c.Get("username")
	username, _ := u.(string)
	db.WriteAudit(username, "pro20x_settings", "updated", c.ClientIP())
	AdminGetPro20xSettings(c)
}

// AdminTestPro20x GET /api/v1/admin/pro20x/test
// 拉一次 /config，确认 Key 可用并回显套餐/地区/积分。
func AdminTestPro20x(c *gin.Context) {
	cli := pro20x.NewFromSettings()
	cfg, err := cli.GetConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"ok": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":           true,
		"auto_card":    cfg.SupportsAutoCard(),
		"task_costs":   cfg.PointsBilling.TaskCosts,
		"plan_regions": cfg.PlanRegions,
	})
}
