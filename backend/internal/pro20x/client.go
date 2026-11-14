// Package pro20x 封装 Pro20x Automate 开放 API（V4）。
//
// 文档：pro20x-automate-api-v4.md（本地） / https://pro20x.com/automate?tab=api
//
// 只做客户端：配置、账号查询、提交任务、轮询任务、取消任务。
// 不碰兑换主流程，是否启用由上层决定。
package pro20x

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tuzi/cdk-recharge-system/internal/db"
)

const (
	defaultBase = "https://pro20x.com/api/v3/automate"
	maxBody     = 4 << 20
)

// Config 从 site_settings 读取（优先）或环境变量。
type Config struct {
	BaseURL string
	APIKey  string
}

// LoadConfig 读取配置。
func LoadConfig() Config {
	base, _ := db.GetSetting("pro20x_api_base")
	key, _ := db.GetSetting("pro20x_api_key")
	if strings.TrimSpace(base) == "" {
		base = defaultBase
	}
	return Config{
		BaseURL: strings.TrimRight(strings.TrimSpace(base), "/"),
		APIKey:  strings.TrimSpace(key),
	}
}

// Client 是 Pro20x Automate 客户端。
type Client struct {
	cfg    Config
	client *http.Client
}

// New 创建客户端。
func New(cfg Config) *Client {
	return &Client{
		cfg: cfg,
		client: &http.Client{
			Timeout: 40 * time.Second,
		},
	}
}

// NewFromSettings 按后台配置创建客户端。
func NewFromSettings() *Client {
	return New(LoadConfig())
}

// APIError 是 Pro20x 返回的业务/HTTP 错误。
type APIError struct {
	HTTPStatus int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("pro20x: http=%d code=%s msg=%s", e.HTTPStatus, e.Code, e.Message)
	}
	return fmt.Sprintf("pro20x: http=%d msg=%s", e.HTTPStatus, e.Message)
}

type envelope struct {
	Success   bool            `json:"success"`
	Code      string          `json:"code"`
	Message   string          `json:"message"`
	RequestID string          `json:"requestId"`
	Job       json.RawMessage `json:"job"`
	Poll      string          `json:"poll"`
	Data      json.RawMessage `json:"data"`
}

// resolveURL 拼接请求地址。
// poll 字段可能返回三种形式：完整 URL、站点绝对路径（/api/v3/...）、接口相对路径（/account-check/...）。
func (c *Client) resolveURL(path string) string {
	p := strings.TrimSpace(path)
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	// 站点绝对路径：只取 BaseURL 的 scheme+host，避免拼出 /api/v3/automate/api/v3/... 的双前缀
	if strings.HasPrefix(p, "/api/") {
		if u, err := url.Parse(c.cfg.BaseURL); err == nil && u.Host != "" {
			return u.Scheme + "://" + u.Host + p
		}
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return c.cfg.BaseURL + p
}

func (c *Client) do(ctx context.Context, method, path string, body any, idem string) (*envelope, []byte, error) {
	if c.cfg.APIKey == "" {
		return nil, nil, &APIError{HTTPStatus: 401, Code: "automate_key_required", Message: "pro20x_api_key not configured"}
	}
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.resolveURL(path), rdr)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Automate-Key", c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if idem != "" {
		req.Header.Set("X-Request-Id", idem)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return nil, nil, err
	}
	var env envelope
	_ = json.Unmarshal(raw, &env)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || (!env.Success && env.Code != "") {
		return &env, raw, &APIError{HTTPStatus: resp.StatusCode, Code: env.Code, Message: env.Message}
	}
	return &env, raw, nil
}

// ---- 配置 / 余额 ----

// AutomateConfig 是 /config 的关键字段。
type AutomateConfig struct {
	Raw          json.RawMessage            `json:"-"`
	Plans        json.RawMessage            `json:"plans"`
	PlanRegions  map[string][]string        `json:"planRegions"`
	PaymentSourc []PaymentSource            `json:"paymentSources"`
	PointsBilling struct {
		TaskCosts map[string]int `json:"taskCosts"`
	} `json:"pointsBilling"`
}

// PaymentSource 是支付方式项（对象数组，按 id 判断）。
type PaymentSource struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// GetConfig 拉取 /config。
func (c *Client) GetConfig(ctx context.Context) (*AutomateConfig, error) {
	_, raw, err := c.do(ctx, http.MethodGet, "/config", nil, "")
	if err != nil {
		return nil, err
	}
	var out AutomateConfig
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("pro20x /config 返回非 JSON")
	}
	out.Raw = raw
	return &out, nil
}

// SupportsAutoCard 判断是否开放自动选卡。
func (cfg *AutomateConfig) SupportsAutoCard() bool {
	for _, s := range cfg.PaymentSourc {
		if s.ID == "auto_card" {
			return true
		}
	}
	return false
}

// TaskCost 返回该套餐的积分消耗（0 表示未配置）。
func (cfg *AutomateConfig) TaskCost(planID string) int {
	if cfg.PointsBilling.TaskCosts == nil {
		return 0
	}
	return cfg.PointsBilling.TaskCosts[planID]
}

// ---- 账号查询 ----

// AccountInfo 是 account-check 结果里的账号条目。
type AccountInfo struct {
	AccountID             string `json:"accountId"`
	PlanLabel             string `json:"planLabel"`
	Tier                  string `json:"tier"`
	HasActiveSubscription bool   `json:"hasActiveSubscription"`
	BillingCurrency       string `json:"billingCurrency"`
	WillRenew             *bool  `json:"willRenew"`
	ExpirationTime        string `json:"expirationTime"`
}

// CheckResult 是 account-check 的结果。
type CheckResult struct {
	Email             string        `json:"email"`
	SelectedAccountID string        `json:"selectedAccountId"`
	Accounts          []AccountInfo `json:"accounts"`
}

// Job 是异步任务。
type Job struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Status    string          `json:"status"`
	ErrorCode string          `json:"errorCode"`
	Result    json.RawMessage `json:"result"`
	ClientOrderID string      `json:"clientOrderId"`
	TaskID        string      `json:"taskId"`
}

// AccountCheck 提交并轮询账号查询（只读，不消耗积分）。
func (c *Client) AccountCheck(ctx context.Context, authSession any) (*CheckResult, error) {
	_, raw, err := c.do(ctx, http.MethodPost, "/account-check",
		map[string]any{"authSessionJson": authSession}, "")
	if err != nil {
		return nil, err
	}
	var env struct {
		Job  Job    `json:"job"`
		Poll string `json:"poll"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("pro20x /account-check 返回非 JSON")
	}
	job, err := c.waitJob(ctx, env.Job, env.Poll, "/account-check/jobs/")
	if err != nil {
		return nil, err
	}
	var res CheckResult
	if err := json.Unmarshal(job.Result, &res); err != nil {
		return nil, fmt.Errorf("account-check 结果解析失败")
	}
	return &res, nil
}

// ---- 提交任务 ----

// CreateTaskRequest 是 /tasks 的请求体（只列常用字段）。
type CreateTaskRequest struct {
	ClientOrderID   string `json:"clientOrderId"`
	PlanID          string `json:"planId"`
	CheckoutCountry string `json:"checkoutCountry"`
	AuthSession     any    `json:"authSessionJson,omitempty"`
	PaymentSource   string `json:"paymentSource,omitempty"`
	Card             *Card  `json:"card,omitempty"`
	CardPlatformID   string `json:"cardPlatformId,omitempty"`
	CreditsQuantity  int    `json:"creditsQuantity,omitempty"`
	CancelAutoRenew  *bool  `json:"cancelAutoRenew,omitempty"`
	DeletePaymentCard *bool `json:"deletePaymentCard,omitempty"`
	DetectPlusPromo   *bool `json:"detectPlusPromo,omitempty"`
}

// Card 是完整卡信息（paymentSource=card 时必填）。
type Card struct {
	Number   string `json:"number"`
	CVC      string `json:"cvc"`
	ExpMonth string `json:"expMonth"`
	ExpYear  string `json:"expYear"`
}

// Task 是充值任务详情。
type Task struct {
	ID            string          `json:"id"`
	ClientOrderID string          `json:"clientOrderId"`
	Status        string          `json:"status"`
	Stage         string          `json:"stage"`
	ErrorCode     string          `json:"errorCode"`
	Error         json.RawMessage `json:"error"`
	PaymentSource string          `json:"paymentSource"`
	Result        json.RawMessage `json:"result"`
	Cleanup       json.RawMessage `json:"gptPaymentMethodCleanup"`
}

// unwrapTask 从多种响应外壳里取出任务对象。
// 兼容：顶层 Task / {task:{...}} / {job:{...}} / {data:{task:{...}}} / {taskId} / {id}
func unwrapTask(raw []byte) (*Task, error) {
	var direct Task
	if err := json.Unmarshal(raw, &direct); err == nil && direct.ID != "" {
		return &direct, nil
	}
	var env struct {
		Task   Task   `json:"task"`
		Job    Task   `json:"job"`
		TaskID string `json:"taskId"`
		ID     string `json:"id"`
		Data   struct {
			Task   Task   `json:"task"`
			Job    Task   `json:"job"`
			TaskID string `json:"taskId"`
			ID     string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("pro20x 返回非 JSON")
	}
	for _, cand := range []*Task{&env.Task, &env.Job, &env.Data.Task, &env.Data.Job} {
		if cand.ID != "" {
			return cand, nil
		}
	}
	for _, id := range []string{env.TaskID, env.ID, env.Data.TaskID, env.Data.ID} {
		if strings.TrimSpace(id) != "" {
			return &Task{ID: strings.TrimSpace(id)}, nil
		}
	}
	return nil, fmt.Errorf("pro20x 响应里没有任务 ID")
}

// CreateTask 提交充值任务。
func (c *Client) CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error) {
	_, raw, err := c.do(ctx, http.MethodPost, "/tasks", req, req.ClientOrderID)
	if err != nil {
		return nil, err
	}
	return unwrapTask(raw)
}

// GetTaskByClientOrderID 按幂等订单号查任务（用于任务 ID 丢失后的补偿查询）。
func (c *Client) GetTaskByClientOrderID(ctx context.Context, clientOrderID string) (*Task, error) {
	if strings.TrimSpace(clientOrderID) == "" {
		return nil, fmt.Errorf("clientOrderId 为空")
	}
	path := "/tasks?clientOrderId=" + url.QueryEscape(strings.TrimSpace(clientOrderID))
	_, raw, err := c.do(ctx, http.MethodGet, path, nil, "")
	if err != nil {
		return nil, err
	}
	var env struct {
		Tasks []Task `json:"tasks"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("pro20x /tasks 列表返回非 JSON")
	}
	if len(env.Tasks) == 0 {
		return nil, nil
	}
	return &env.Tasks[0], nil
}

// GetTask 查询任务详情。
func (c *Client) GetTask(ctx context.Context, id string) (*Task, error) {
	_, raw, err := c.do(ctx, http.MethodGet, "/tasks/"+url.PathEscape(id), nil, "")
	if err != nil {
		return nil, err
	}
	return unwrapTask(raw)
}

// CancelTask 取消任务（仅支付提交前）。
func (c *Client) CancelTask(ctx context.Context, id, note string) error {
	_, _, err := c.do(ctx, http.MethodPost, "/tasks/"+url.PathEscape(id)+"/force-terminate",
		map[string]string{"note": note}, "")
	return err
}

// ---- 任务状态 ----

const (
	StatusQueued          = "queued"
	StatusRunning         = "running"
	StatusCaptchaRequired = "captcha_required"
	StatusSucceeded       = "succeeded"
	StatusFailed          = "failed"
	StatusManualReview    = "manual_review"
	StatusCancelled       = "cancelled"
)

// IsTerminal 判断任务是否终态。
// 注意：异步 job（account-check / upgrade-quote / cancel）用 "success"，
// 充值任务用 "succeeded"，两者都要认。
func IsTerminal(status string) bool {
	switch status {
	case StatusSucceeded, "success", StatusFailed, StatusManualReview, StatusCancelled:
		return true
	}
	return false
}

// isFailure 判断终态是否为失败。
func isFailure(status string) bool {
	return status == StatusFailed || status == "error"
}

// waitJob 轮询异步 job 直到终态。
func (c *Client) waitJob(ctx context.Context, job Job, poll, pathPrefix string) (*Job, error) {
	if poll == "" {
		poll = pathPrefix + job.ID
	}
	deadline := time.Now().Add(3 * time.Minute)
	for {
		if IsTerminal(job.Status) {
			if isFailure(job.Status) {
				return nil, &APIError{Code: job.ErrorCode, Message: "job failed"}
			}
			return &job, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("pro20x job 轮询超时")
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(3 * time.Second):
		}
		_, raw, err := c.do(ctx, http.MethodGet, poll, nil, "")
		if err != nil {
			return nil, err
		}
		var env struct {
			Job Job `json:"job"`
		}
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, fmt.Errorf("job 返回非 JSON")
		}
		job = env.Job
	}
}
