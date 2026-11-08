// Package gptcard 用 ChatGPT session 管理账号已绑定的支付方式（payment method）。
//
// 只做三件事：列出 payment_methods、把某张卡设为默认、删除某张卡。
// 不碰订阅、不碰 checkout、不碰账号资料，也不调用卡台的任何写接口。
package gptcard

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Session 是调用 chatgpt.com/backend-api/payments/* 所需的凭据。
type Session struct {
	// Raw 是原始 session JSON（或 sessionToken / accessToken），刷新时复用。
	Raw string
	// AccessToken 用于 Authorization: Bearer（约 1 小时有效）。
	AccessToken string
	// SessionToken 用于 /api/auth/session 换新 accessToken（月级有效）。
	SessionToken string
	// AccountID 是 chatgpt_account_id，既作 query 参数也作 chatgpt-account-id 头。
	AccountID string
}

// ParseSession 接受完整 session JSON、裸 sessionToken（JWE）或裸 accessToken。
func ParseSession(raw string) (*Session, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return nil, fmt.Errorf("session 为空")
	}
	s := &Session{Raw: v}
	switch {
	case strings.HasPrefix(v, "{"):
		var data map[string]any
		if err := json.Unmarshal([]byte(v), &data); err != nil {
			return nil, fmt.Errorf("session JSON 解析失败")
		}
		s.AccessToken = firstString(data["accessToken"], data["access_token"])
		s.SessionToken = firstString(data["sessionToken"], data["session_token"])
		if acc, ok := data["account"].(map[string]any); ok {
			s.AccountID = firstString(acc["id"], acc["account_id"])
		}
	case strings.Count(v, ".") >= 4:
		// JWE sessionToken 是 5 段
		s.SessionToken = v
	default:
		s.AccessToken = v
	}
	if s.AccountID == "" && s.AccessToken != "" {
		s.AccountID = accountIDFromJWT(s.AccessToken)
	}
	if s.AccessToken == "" && s.SessionToken == "" {
		return nil, fmt.Errorf("未找到 accessToken/sessionToken")
	}
	return s, nil
}

// accountIDFromJWT 从 accessToken 的 payload 里取 chatgpt_account_id。
func accountIDFromJWT(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ""
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}
	var claims struct {
		Auth struct {
			AccountID string `json:"chatgpt_account_id"`
		} `json:"https://api.openai.com/auth"`
	}
	if json.Unmarshal(payload, &claims) != nil {
		return ""
	}
	return strings.TrimSpace(claims.Auth.AccountID)
}

func firstString(vals ...any) string {
	for _, v := range vals {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}
