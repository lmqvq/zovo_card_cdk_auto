package gptcard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

const (
	defaultUA    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	maxBodyBytes = 2 << 20
	sessionChunk = 3933
)

// Card 是一张已绑定的支付方式。
type Card struct {
	ID        string
	Brand     string
	Last4     string
	ExpMonth  int
	ExpYear   int
	Created   int64
	IsDefault bool
}

// Label 返回 "mastercard •••• 5794" 形式的展示串。
func (c Card) Label() string {
	brand := strings.TrimSpace(c.Brand)
	if brand == "" {
		brand = "card"
	}
	if c.Last4 == "" {
		return brand
	}
	return brand + " •••• " + c.Last4
}

// CardList 是列卡结果。
type CardList struct {
	Cards     []Card
	DefaultID string
}

// APIError 是 chatgpt.com 返回的非 200 响应。
type APIError struct {
	Status int
	Body   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("chatgpt HTTP %d: %s", e.Status, e.Body)
}

// Unauthorized 表示凭据失效（401/403）。
func (e *APIError) Unauthorized() bool { return e.Status == 401 || e.Status == 403 }

// Client 复用 tls-client（Chrome 指纹）。chatgpt.com 在 Cloudflare 后，net/http 会被拦。
type Client struct {
	httpc tls_client.HttpClient
	base  string
}

// NewClient 创建一个带浏览器指纹的客户端。
func NewClient() (*Client, error) {
	httpc, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(),
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_131),
		tls_client.WithForceHttp1(),
		tls_client.WithNotFollowRedirects(),
	)
	if err != nil {
		return nil, fmt.Errorf("初始化 tls-client 失败: %w", err)
	}
	return &Client{httpc: httpc, base: "https://chatgpt.com"}, nil
}

// SetBaseURL 只给测试用（默认 https://chatgpt.com）。
func (c *Client) SetBaseURL(u string) {
	if s := strings.TrimRight(strings.TrimSpace(u), "/"); s != "" {
		c.base = s
	}
}

// SetDefaultOptions 是"设为默认"的请求模板，后台可覆盖。
type SetDefaultOptions struct {
	Method string
	Path   string
	Body   string
}

// DefaultSetDefaultOptions 是抓包核对过的默认值。
func DefaultSetDefaultOptions() SetDefaultOptions {
	return SetDefaultOptions{
		Method: "POST",
		Path:   "/backend-api/payments/payment_method/default",
		Body:   `{"account_id":"{account_id}","payment_method_id":"{pm_id}"}`,
	}
}

func (c *Client) do(method, rawURL string, headers map[string]string, body []byte, cookies []*fhttp.Cookie) (int, []byte, error) {
	var rdr io.Reader
	if len(body) > 0 {
		rdr = bytes.NewReader(body)
	}
	req, err := fhttp.NewRequest(method, rawURL, rdr)
	if err != nil {
		return 0, nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	for _, ck := range cookies {
		req.AddCookie(ck)
	}
	resp, err := c.httpc.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	return resp.StatusCode, b, err
}

func (c *Client) headers(s *Session) map[string]string {
	h := map[string]string{
		"Accept":        "*/*",
		"Authorization": "Bearer " + s.AccessToken,
		"User-Agent":    defaultUA,
		"Origin":        "https://chatgpt.com",
		"Referer":       "https://chatgpt.com/",
	}
	if s.AccountID != "" {
		h["chatgpt-account-id"] = s.AccountID
	}
	return h
}

// List 列出账号已绑定的支付方式。
func (c *Client) List(s *Session) (*CardList, error) {
	if s.AccountID == "" {
		return nil, fmt.Errorf("缺少 account_id")
	}
	u := c.base + "/backend-api/payments/payment_methods?account_id=" + url.QueryEscape(s.AccountID)
	st, raw, err := c.do(fhttp.MethodGet, u, c.headers(s), nil, nil)
	if err != nil {
		return nil, err
	}
	if st != 200 {
		return nil, &APIError{Status: st, Body: snippet(raw)}
	}
	var data struct {
		PaymentMethods []map[string]any `json:"payment_methods"`
		DefaultID      string           `json:"default_payment_method_id"`
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("列卡返回非 JSON")
	}
	out := &CardList{DefaultID: data.DefaultID}
	for _, m := range data.PaymentMethods {
		out.Cards = append(out.Cards, parseCard(m, data.DefaultID))
	}
	return out, nil
}

// Delete 删除一张支付方式，成功返回 success=true。
func (c *Client) Delete(s *Session, pmID string) (int, bool, error) {
	u := c.base + "/backend-api/payments/payment_method/" + url.PathEscape(pmID) + "?account_id=" + url.QueryEscape(s.AccountID)
	st, raw, err := c.do(fhttp.MethodDelete, u, c.headers(s), nil, nil)
	if err != nil {
		return 0, false, err
	}
	ok := false
	var resp map[string]any
	if json.Unmarshal(raw, &resp) == nil {
		if b, isBool := resp["success"].(bool); isBool {
			ok = b
		}
	}
	return st, ok, nil
}

// SetDefault 把某张卡设为默认（删除默认卡的前置动作）。
func (c *Client) SetDefault(s *Session, pmID string, opt SetDefaultOptions) (int, error) {
	method := strings.TrimSpace(opt.Method)
	if method == "" {
		method = "POST"
	}
	path := strings.TrimSpace(opt.Path)
	if path == "" {
		path = "/backend-api/payments/payment_method/default"
	}
	body := strings.ReplaceAll(opt.Body, "{pm_id}", pmID)
	body = strings.ReplaceAll(body, "{account_id}", s.AccountID)
	h := c.headers(s)
	h["Content-Type"] = "application/json"
	st, raw, err := c.do(strings.ToUpper(method), c.base+path, h, []byte(body), nil)
	if err != nil {
		return 0, err
	}
	if st != 200 && st != 204 {
		return st, &APIError{Status: st, Body: snippet(raw)}
	}
	return st, nil
}

// Refresh 用 sessionToken 换一个新的 accessToken。
func (c *Client) Refresh(s *Session) error {
	if s.SessionToken == "" {
		return fmt.Errorf("缺少 sessionToken")
	}
	headers := map[string]string{
		"Accept":     "application/json",
		"User-Agent": defaultUA,
		"Referer":    "https://chatgpt.com/",
	}
	st, raw, err := c.do(fhttp.MethodGet, c.base+"/api/auth/session", headers, nil, sessionTokenCookies(s.SessionToken))
	if err != nil {
		return err
	}
	if st != 200 {
		return fmt.Errorf("刷新 session 失败 HTTP %d", st)
	}
	var data map[string]any
	if json.Unmarshal(raw, &data) != nil {
		return fmt.Errorf("刷新 session 返回非 JSON")
	}
	tok := firstString(data["accessToken"], data["access_token"])
	if tok == "" {
		return fmt.Errorf("刷新 session 未返回 accessToken")
	}
	s.AccessToken = tok
	if s.AccountID == "" {
		s.AccountID = accountIDFromJWT(tok)
	}
	return nil
}

func sessionTokenCookies(token string) []*fhttp.Cookie {
	const base = "__Secure-next-auth.session-token"
	if len(token) <= sessionChunk {
		return []*fhttp.Cookie{{Name: base, Value: token}}
	}
	var out []*fhttp.Cookie
	for i, idx := 0, 0; i < len(token); i, idx = i+sessionChunk, idx+1 {
		end := i + sessionChunk
		if end > len(token) {
			end = len(token)
		}
		out = append(out, &fhttp.Cookie{Name: fmt.Sprintf("%s.%d", base, idx), Value: token[i:end]})
	}
	return out
}

func parseCard(m map[string]any, defaultID string) Card {
	card, _ := m["card"].(map[string]any)
	c := Card{
		ID:    str(m["id"]),
		Brand: firstString(card["brand"], m["brand"]),
		Last4: firstString(card["last4"], m["last4"]),
	}
	c.ExpMonth = intOf(card["exp_month"])
	c.ExpYear = intOf(card["exp_year"])
	c.Created = int64Of(m["created"])
	c.IsDefault = c.ID != "" && c.ID == defaultID
	return c
}

func str(v any) string {
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

func intOf(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case json.Number:
		n, _ := t.Int64()
		return int(n)
	}
	return 0
}

func int64Of(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case int:
		return int64(t)
	case json.Number:
		n, _ := t.Int64()
		return n
	}
	return 0
}

func snippet(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 200 {
		s = s[:200] + "…"
	}
	return s
}
