package gptcard

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tuzi/cdk-recharge-system/internal/db"
	"github.com/tuzi/cdk-recharge-system/internal/gptcheck"
)

// 后台配置键（都存在 site_settings，后台可改，改完无需发版）
const (
	settingEnabled          = "card_delete_enabled"
	settingDelaySeconds     = "card_delete_delay_seconds"
	settingRenewAttempts    = "card_delete_renew_check_attempts"
	settingMaxAttempts      = "card_delete_max_attempts"
	settingRatePerMinute    = "card_delete_rate_per_minute"
	settingConcurrency      = "card_delete_concurrency"
	settingSetDefaultMethod = "card_delete_setdefault_method"
	settingSetDefaultPath   = "card_delete_setdefault_path"
	settingSetDefaultBody   = "card_delete_setdefault_body"
)

const (
	defaultDelaySeconds  = 60
	defaultRenewAttempts = 3
	defaultMaxAttempts   = 3
	defaultRatePerMinute = 30
	defaultConcurrency   = 2
	scanInterval         = 5 * time.Minute
	accountGap           = 5 * time.Second
)

var (
	rateMu   sync.Mutex
	lastReq  time.Time
	acctMu   sync.Mutex
	acctLast = map[string]time.Time{}
	running  sync.Map // job_id -> struct{}，保证同一任务不会被并发执行
)

// Enabled 读取总开关。
func Enabled() bool {
	v, _ := db.GetSetting(settingEnabled)
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "on", "yes":
		return true
	}
	return false
}

func settingInt(key string, def int) int {
	v, _ := db.GetSetting(key)
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil || n <= 0 {
		return def
	}
	return n
}

func settingStr(key, def string) string {
	v, _ := db.GetSetting(key)
	if s := strings.TrimSpace(v); s != "" {
		return s
	}
	return def
}

func setDefaultOptions() SetDefaultOptions {
	def := DefaultSetDefaultOptions()
	return SetDefaultOptions{
		Method: settingStr(settingSetDefaultMethod, def.Method),
		Path:   settingStr(settingSetDefaultPath, def.Path),
		Body:   settingStr(settingSetDefaultBody, def.Body),
	}
}

// EnqueueFromWebhook 由 gpt_direct.completed webhook 调用：
// 只做过滤 + 幂等入库 + 延迟异步执行，绝不阻塞 webhook 返回。
func EnqueueFromWebhook(payload map[string]any) {
	if payload == nil || !Enabled() {
		return
	}
	if strings.ToLower(str(payload["product"])) != "gpt" {
		return
	}
	plan := strings.ToLower(str(payload["plan"]))
	if strings.Contains(plan, "renew") {
		return // 绑卡档（Pro 20x / 5x 续费）的卡是给客户长期扣费用的，永不删
	}
	if st := strings.ToLower(str(payload["status"])); st != "" && st != "completed" {
		return
	}
	orderID := int64Of(payload["order_id"])
	cdkID := int64Of(payload["cdk_id"])
	if orderID <= 0 {
		return
	}
	last4 := str(payload["card_last_four"])
	delay := settingInt(settingDelaySeconds, defaultDelaySeconds)
	next := time.Now().UTC().Add(time.Duration(delay) * time.Second).Format("2006-01-02 15:04:05")

	id, inserted, err := db.InsertCardDeleteJob(db.CardDeleteJob{
		CDKID:        cdkID,
		OrderID:      orderID,
		Plan:         str(payload["plan"]),
		AccountEmail: str(payload["account_email"]),
		CardLastFour: last4,
		Action:       db.CardDeletePending,
		NextRunAt:    next,
	})
	if err != nil {
		log.Printf("[card-delete] enqueue order=%d failed: %v", orderID, err)
		return
	}
	if !inserted {
		return // 重复 webhook / 已处理过
	}
	go func() {
		time.Sleep(time.Duration(delay) * time.Second)
		RunJob(id)
	}()
}

// RunJob 执行一条解绑任务。幂等：只有 pending 才会真正执行，同一任务不会并发跑。
func RunJob(jobID int64) {
	if _, loaded := running.LoadOrStore(jobID, struct{}{}); loaded {
		return
	}
	defer running.Delete(jobID)

	job, err := db.GetCardDeleteJob(jobID)
	if err != nil || job == nil {
		return
	}
	if job.Action != db.CardDeletePending {
		return
	}
	if !Enabled() {
		return
	}

	fail := func(action, reason string, httpStatus int) {
		job.Action = action
		job.Reason = reason
		if httpStatus > 0 {
			job.HTTPStatus = httpStatus
		}
		job.NextRunAt = ""
		if err := db.UpdateCardDeleteJob(*job); err != nil {
			log.Printf("[card-delete] save job=%d failed: %v", job.ID, err)
		}
		log.Printf("[card-delete] order=%d last4=%s -> %s(%s)", job.OrderID, job.CardLastFour, action, reason)
	}

	if len(job.CardLastFour) != 4 {
		fail(db.CardDeleteSkipped, "no_last4", 0)
		return
	}

	sess, err := loadSession(job)
	if err != nil {
		fail(db.CardDeleteSkipped, "no_session", 0)
		return
	}

	cli, err := NewClient()
	if err != nil {
		fail(db.CardDeleteFailed, "client_init", 0)
		return
	}

	if sess.AccessToken == "" {
		if err := cli.Refresh(sess); err != nil {
			job.SessionState = "invalid"
			job.SessionCheckedAt = nowString()
			fail(db.CardDeleteFailed, "token_invalid", 0)
			return
		}
	}

	waitAccount(job.AccountEmail)

	// 1) 列卡（同时验证凭据）
	waitRate()
	cards, err := cli.List(sess)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok && apiErr.Unauthorized() {
			if rerr := cli.Refresh(sess); rerr == nil {
				waitRate()
				cards, err = cli.List(sess)
			}
		}
	}
	if err != nil {
		if apiErr, ok := err.(*APIError); ok && apiErr.Unauthorized() {
			job.SessionState = "invalid"
			job.SessionCheckedAt = nowString()
			fail(db.CardDeleteFailed, "token_invalid", apiErr.Status)
			return
		}
		backoff(job, "network")
		return
	}
	job.SessionState = "valid"
	job.SessionCheckedAt = nowString()

	// 2) 必须确认不会再自动续费，才删卡
	willRenew, queried := queryWillRenew(sess)
	if !queried {
		backoff(job, "subscription_unknown")
		return
	}
	if willRenew {
		backoff(job, "still_renewing")
		return
	}

	// 3) 判定
	job.CardsBefore = len(cards.Cards)
	if len(cards.Cards) <= 1 {
		fail(db.CardDeleteSkipped, "last_card", 0)
		return
	}
	target, reason := pickTarget(cards, job.CardLastFour)
	if reason == "ambiguous" {
		fail(db.CardDeleteReview, "ambiguous", 0)
		return
	}
	if target == nil {
		fail(db.CardDeleteSkipped, "not_found", 0)
		return
	}
	job.PMID = target.ID
	job.PMBrand = target.Brand

	// 4) 默认卡必须先换默认，否则 Stripe 会拒绝 detach
	if target.IsDefault {
		other := pickOther(cards, target.ID)
		if other == nil {
			fail(db.CardDeleteSkipped, "default_card", 0)
			return
		}
		waitRate()
		if st, err := cli.SetDefault(sess, other.ID, setDefaultOptions()); err != nil {
			log.Printf("[card-delete] order=%d setdefault pm=%s http=%d err=%v", job.OrderID, other.ID, st, err)
			fail(db.CardDeleteSkipped, "default_card", st)
			return
		}
		job.SetDefaultUsed = true
		waitRate()
	}

	// 5) 删卡
	st, ok, err := cli.Delete(sess, target.ID)
	if err != nil {
		backoff(job, "network")
		return
	}
	job.HTTPStatus = st
	if st == 200 && ok {
		fail(db.CardDeleteDone, "", 200)
		return
	}
	if st == 401 || st == 403 {
		job.SessionState = "invalid"
		job.SessionCheckedAt = nowString()
		fail(db.CardDeleteFailed, "token_invalid", st)
		return
	}
	if st == 400 || st == 409 {
		fail(db.CardDeleteSkipped, "default_card", st)
		return
	}
	backoff(job, fmt.Sprintf("http_%d", st))
}

// backoff 暂缓执行；达到上限后落到 skipped/failed。
func backoff(job *db.CardDeleteJob, reason string) {
	job.Attempts++
	switch reason {
	case "still_renewing", "subscription_unknown":
		limit := settingInt(settingRenewAttempts, defaultRenewAttempts)
		if job.Attempts >= limit {
			job.Action = db.CardDeleteSkipped
			job.Reason = reason
			job.NextRunAt = ""
			_ = db.UpdateCardDeleteJob(*job)
			log.Printf("[card-delete] order=%d give up after %d renew checks (%s)", job.OrderID, job.Attempts, reason)
			return
		}
		// 10m / 1h / 3h 退避
		steps := []time.Duration{10 * time.Minute, 60 * time.Minute, 3 * time.Hour}
		idx := job.Attempts - 1
		if idx < 0 {
			idx = 0
		}
		if idx >= len(steps) {
			idx = len(steps) - 1
		}
		job.NextRunAt = time.Now().UTC().Add(steps[idx]).Format("2006-01-02 15:04:05")
	default:
		limit := settingInt(settingMaxAttempts, defaultMaxAttempts)
		if job.Attempts >= limit {
			job.Action = db.CardDeleteFailed
			job.Reason = reason
			job.NextRunAt = ""
			_ = db.UpdateCardDeleteJob(*job)
			log.Printf("[card-delete] order=%d failed after %d attempts (%s)", job.OrderID, job.Attempts, reason)
			return
		}
		job.NextRunAt = time.Now().UTC().Add(2 * time.Minute).Format("2006-01-02 15:04:05")
	}
	job.Action = db.CardDeletePending
	job.Reason = reason
	_ = db.UpdateCardDeleteJob(*job)
}

// RetryJob 手动重试：重置计数后立即执行。
func RetryJob(jobID int64) error {
	job, err := db.GetCardDeleteJob(jobID)
	if err != nil {
		return err
	}
	if job == nil {
		return fmt.Errorf("任务不存在")
	}
	job.Action = db.CardDeletePending
	job.Reason = ""
	job.Attempts = 0
	job.NextRunAt = nowString()
	if err := db.UpdateCardDeleteJob(*job); err != nil {
		return err
	}
	go RunJob(job.ID)
	return nil
}

// CheckSession 手动校验某条任务背后的 session 是否还有效（管理页按钮）。
func CheckSession(cdkID int64) (string, error) {
	if cdkID <= 0 {
		return "unknown", fmt.Errorf("cdk_id required")
	}
	raw, err := sessionByCDK(cdkID)
	if err != nil || strings.TrimSpace(raw) == "" {
		return "unknown", fmt.Errorf("未找到 session")
	}
	sess, err := ParseSession(raw)
	if err != nil {
		return "invalid", nil
	}
	cli, err := NewClient()
	if err != nil {
		return "unknown", err
	}
	if sess.AccessToken == "" {
		if err := cli.Refresh(sess); err != nil {
			return "invalid", nil
		}
	}
	waitRate()
	if _, err := cli.List(sess); err != nil {
		if apiErr, ok := err.(*APIError); ok && apiErr.Unauthorized() {
			return "invalid", nil
		}
		return "unknown", err
	}
	return "valid", nil
}

// Start 启动后台补跑与到期扫描。
func Start(ctx context.Context) {
	go run(ctx)
}

func run(ctx context.Context) {
	scanOnce(ctx)
	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[card-delete] stopped")
			return
		case <-ticker.C:
			scanOnce(ctx)
			runDue(ctx)
		}
	}
}

// scanOnce 从 webhook_events 补入漏掉的任务（服务重启会丢 goroutine）。
func scanOnce(ctx context.Context) {
	if !Enabled() {
		return
	}
	payloads, err := db.ListCompletedWebhookPayloads(7, 200)
	if err != nil {
		log.Printf("[card-delete] scan payloads: %v", err)
		return
	}
	for _, p := range payloads {
		select {
		case <-ctx.Done():
			return
		default:
		}
		var m map[string]any
		if json.Unmarshal([]byte(p), &m) != nil {
			continue
		}
		// 已经入过队的跳过（唯一键；重复入队也会被 INSERT OR IGNORE 挡掉）
		job, err := db.FindCardDeleteJob(int64Of(m["order_id"]), str(m["card_last_four"]))
		if err == nil && job != nil {
			continue
		}
		EnqueueFromWebhook(m)
	}
}

// runDue 执行到期的 pending 任务（受速率与并发限制）。
func runDue(ctx context.Context) {
	limit := settingInt(settingConcurrency, defaultConcurrency)
	sem := make(chan struct{}, limit)
	for round := 0; round < 20; round++ {
		jobs, err := db.ListDueCardDeleteJobs(50)
		if err != nil || len(jobs) == 0 {
			return
		}
		for _, j := range jobs {
			select {
			case <-ctx.Done():
				return
			default:
			}
			sem <- struct{}{}
			go func(id int64) {
				defer func() {
					<-sem
					if r := recover(); r != nil {
						log.Printf("[card-delete] panic: %v", r)
					}
				}()
				RunJob(id)
			}(j.ID)
		}
		time.Sleep(2 * time.Second)
	}
}

func loadSession(job *db.CardDeleteJob) (*Session, error) {
	raw, err := sessionByCDK(job.CDKID)
	if err != nil || strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("no session")
	}
	return ParseSession(raw)
}

func sessionByCDK(cdkID int64) (string, error) {
	code, ok := db.LookupCardplatformCDKCode(cdkID, "")
	if !ok || strings.TrimSpace(code) == "" {
		return "", fmt.Errorf("cdk %d not found", cdkID)
	}
	return db.GetSessionByCDK(code)
}

func queryWillRenew(s *Session) (willRenew bool, queried bool) {
	input := s.Raw
	if s.AccessToken != "" {
		b, _ := json.Marshal(map[string]string{
			"accessToken":  s.AccessToken,
			"sessionToken": s.SessionToken,
		})
		input = string(b)
	}
	res, err := gptcheck.Check(input)
	if err != nil || res == nil {
		return false, false
	}
	willRenew = false
	if b, ok := res.Summary["will_renew"].(bool); ok {
		willRenew = b
	}
	if active, ok := res.Summary["has_active_subscription"].(bool); ok && !active {
		willRenew = false
	}
	return willRenew, true
}

// pickTarget 按尾号定位本次支付卡。多张命中 → ambiguous（人工）。
func pickTarget(list *CardList, last4 string) (*Card, string) {
	if list == nil {
		return nil, "not_found"
	}
	var matched []Card
	for _, c := range list.Cards {
		if c.Last4 == last4 {
			matched = append(matched, c)
		}
	}
	switch len(matched) {
	case 0:
		return nil, "not_found"
	case 1:
		c := matched[0]
		return &c, ""
	default:
		return nil, "ambiguous"
	}
}

func pickOther(list *CardList, targetID string) *Card {
	for _, c := range list.Cards {
		if c.ID != targetID {
			cc := c
			return &cc
		}
	}
	return nil
}

// waitRate 全局速率闸门（默认 30 请求/分钟）。
func waitRate() {
	rateMu.Lock()
	defer rateMu.Unlock()
	perMin := settingInt(settingRatePerMinute, defaultRatePerMinute)
	gap := time.Minute / time.Duration(perMin)
	if gap < time.Second {
		gap = time.Second
	}
	if d := gap - time.Since(lastReq); d > 0 {
		time.Sleep(d)
	}
	lastReq = time.Now()
}

// waitAccount 同一账号两次操作至少间隔 5 秒（预约式，不阻塞其他账号）。
func waitAccount(key string) {
	if strings.TrimSpace(key) == "" {
		return
	}
	acctMu.Lock()
	t, ok := acctLast[key]
	acctLast[key] = time.Now().Add(accountGap)
	acctMu.Unlock()
	if ok {
		if d := t.Sub(time.Now()); d > 0 {
			time.Sleep(d)
		}
	}
}

func nowString() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}
