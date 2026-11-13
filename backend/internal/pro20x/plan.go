package pro20x

import (
	"strings"

	"github.com/tuzi/cdk-recharge-system/internal/db"
)

// 后台配置键
const (
	SettingEnabled        = "pro20x_enabled"
	SettingAPIBase        = "pro20x_api_base"
	SettingAPIKey         = "pro20x_api_key"
	SettingDefaultCountry = "pro20x_default_country"
)

// 默认充值地区（CDK 没记地区时使用）
const defaultCountry = "PH"

// supportedCountries Pro20x 当前支持的充值地区（以 /config 为准，这里做白名单兜底）。
var supportedCountries = map[string]bool{
	"PH": true, "EG": true, "JP": true, "US": true, "CL": true, "SG": true,
	"TZ": true, "SE": true,
}

// Enabled 读取 Pro20x 通道总开关。
func Enabled() bool {
	v, _ := db.GetSetting(SettingEnabled)
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "on", "yes":
		return true
	}
	return false
}

// DefaultCountry 返回 CDK 没记地区时使用的充值地区。
func DefaultCountry() string {
	v, _ := db.GetSetting(SettingDefaultCountry)
	if c := NormalizeCountry(v, ""); c != "" {
		return c
	}
	return defaultCountry
}

// NormalizeCountry 校验并标准化地区代码；不支持时返回 fallback。
func NormalizeCountry(code, fallback string) string {
	c := strings.ToUpper(strings.TrimSpace(code))
	if supportedCountries[c] {
		return c
	}
	return fallback
}

// PlanMapping 是 CDK 套餐到 Pro20x 请求参数的映射。
type PlanMapping struct {
	PlanID          string
	CreditsQuantity int // 仅 credits 用
}

// MapPlan 把本站 CDK 套餐映射成 Pro20x 的 planId。
//
// 套餐名以 cardplatform_cdk_codes.plan 为准：
// plus / pro_5x / pro_20x / credit250 / credit500 / credit1000。
func MapPlan(cdkPlan string) (PlanMapping, bool) {
	p := strings.ToLower(strings.TrimSpace(cdkPlan))
	switch p {
	case "plus":
		return PlanMapping{PlanID: "plus-monthly"}, true
	case "pro_5x", "pro5x", "5x":
		return PlanMapping{PlanID: "pro5x-monthly"}, true
	case "pro_20x", "pro20x", "20x":
		return PlanMapping{PlanID: "pro20x-direct-monthly"}, true
	case "promax", "pro_max", "50x":
		return PlanMapping{PlanID: "promax-monthly"}, true
	case "go":
		return PlanMapping{PlanID: "go-monthly"}, true
	case "credit250":
		return PlanMapping{PlanID: "credits", CreditsQuantity: 250}, true
	case "credit500":
		return PlanMapping{PlanID: "credits", CreditsQuantity: 500}, true
	case "credit1000":
		return PlanMapping{PlanID: "credits", CreditsQuantity: 1000}, true
	}
	return PlanMapping{}, false
}

// BuildOptions 是组装 Pro20x 充值请求的参数。
type BuildOptions struct {
	ClientOrderID   string
	CDKPlan         string
	Country         string
	FallbackCountry string // 通常传 DefaultCountry()
	AuthSession     any
	PaymentSource   string // auto_card / wallet / card
	CardPlatformID  string
}

// BuildCreateRequest 组装 Pro20x 充值请求。
//
// paymentSource: "auto_card"（走绑定的 zovo 卡台）或 "wallet"。
// CardPlatformID 为空时交给 Pro20x 按用卡策略选卡台。
func BuildCreateRequest(opt BuildOptions) (*CreateTaskRequest, error) {
	m, ok := MapPlan(opt.CDKPlan)
	if !ok {
		return nil, &APIError{Code: "plan_not_supported", Message: "CDK 套餐未映射到 Pro20x: " + opt.CDKPlan}
	}
	fallback := NormalizeCountry(opt.FallbackCountry, defaultCountry)
	c := NormalizeCountry(opt.Country, fallback)
	yes := true
	req := &CreateTaskRequest{
		ClientOrderID:     opt.ClientOrderID,
		PlanID:            m.PlanID,
		CheckoutCountry:   c,
		AuthSession:       opt.AuthSession,
		PaymentSource:     opt.PaymentSource,
		CardPlatformID:    opt.CardPlatformID,
		CancelAutoRenew:   &yes,
		DeletePaymentCard: &yes, // 关键：支付完成后由 Pro20x 删卡并核验
	}
	if m.CreditsQuantity > 0 {
		req.CreditsQuantity = m.CreditsQuantity
	}
	if m.PlanID == "plus-monthly" {
		req.DetectPlusPromo = &yes
	}
	return req, nil
}
