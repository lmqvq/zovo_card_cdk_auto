package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// 解绑卡任务状态
const (
	CardDeletePending = "pending" // 排队/等待取消续费
	CardDeleteDone    = "done"    // 已解绑
	CardDeleteSkipped = "skipped" // 按规则不处理
	CardDeleteReview  = "review"  // 需要人工确认
	CardDeleteFailed  = "failed"  // 执行失败
)

// CardDeleteJob 是一条解绑卡任务。
type CardDeleteJob struct {
	ID               int64
	CDKID            int64
	OrderID          int64
	Plan             string
	AccountEmail     string
	CardLastFour     string
	PMID             string
	PMBrand          string
	CardsBefore      int
	Action           string
	Reason           string
	HTTPStatus       int
	Attempts         int
	SetDefaultUsed   bool
	SessionState     string
	SessionCheckedAt string
	NextRunAt        string
	CreatedAt        string
	UpdatedAt        string
}

func sqliteNow() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

// InsertCardDeleteJob 幂等插入：同 order_id + 尾号只留一条。
// inserted=false 表示已存在（重复 webhook / 扫描补跑）。
func InsertCardDeleteJob(j CardDeleteJob) (id int64, inserted bool, err error) {
	if DB == nil {
		return 0, false, fmt.Errorf("db not ready")
	}
	if strings.TrimSpace(j.Action) == "" {
		j.Action = CardDeletePending
	}
	res, err := DB.Exec(`
		INSERT OR IGNORE INTO card_delete_jobs
			(cdk_id, order_id, plan, account_email, card_last_four, action, reason, attempts, next_run_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?)`,
		j.CDKID, j.OrderID, j.Plan, j.AccountEmail, j.CardLastFour, j.Action, j.Reason, j.NextRunAt,
	)
	if err != nil {
		return 0, false, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return 0, false, nil
	}
	newID, err := res.LastInsertId()
	if err != nil {
		return 0, false, err
	}
	return newID, true, nil
}

// UpdateCardDeleteJob 按 id 更新结果字段。
func UpdateCardDeleteJob(j CardDeleteJob) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	setDefault := 0
	if j.SetDefaultUsed {
		setDefault = 1
	}
	_, err := DB.Exec(`
		UPDATE card_delete_jobs SET
			pm_id = ?, pm_brand = ?, cards_before = ?, action = ?, reason = ?,
			http_status = ?, attempts = ?, setdefault_used = ?,
			session_state = ?, session_checked_at = ?, next_run_at = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		j.PMID, j.PMBrand, j.CardsBefore, j.Action, j.Reason,
		j.HTTPStatus, j.Attempts, setDefault,
		j.SessionState, nullIfEmpty(j.SessionCheckedAt), nullIfEmpty(j.NextRunAt),
		j.ID,
	)
	return err
}

// GetCardDeleteJob 按 id 取任务。
func GetCardDeleteJob(id int64) (*CardDeleteJob, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	row := DB.QueryRow(`
		SELECT id, cdk_id, order_id, COALESCE(plan,''), COALESCE(account_email,''), COALESCE(card_last_four,''),
		       COALESCE(pm_id,''), COALESCE(pm_brand,''), COALESCE(cards_before,0),
		       COALESCE(action,''), COALESCE(reason,''), COALESCE(http_status,0), COALESCE(attempts,0),
		       COALESCE(setdefault_used,0), COALESCE(session_state,'unknown'),
		       COALESCE(session_checked_at,''), COALESCE(next_run_at,''), COALESCE(created_at,''), COALESCE(updated_at,'')
		FROM card_delete_jobs WHERE id = ?`, id)
	return scanCardDeleteJob(row)
}

// FindCardDeleteJob 按订单 + 尾号取任务（唯一键）。
func FindCardDeleteJob(orderID int64, last4 string) (*CardDeleteJob, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	row := DB.QueryRow(`
		SELECT id, cdk_id, order_id, COALESCE(plan,''), COALESCE(account_email,''), COALESCE(card_last_four,''),
		       COALESCE(pm_id,''), COALESCE(pm_brand,''), COALESCE(cards_before,0),
		       COALESCE(action,''), COALESCE(reason,''), COALESCE(http_status,0), COALESCE(attempts,0),
		       COALESCE(setdefault_used,0), COALESCE(session_state,'unknown'),
		       COALESCE(session_checked_at,''), COALESCE(next_run_at,''), COALESCE(created_at,''), COALESCE(updated_at,'')
		FROM card_delete_jobs WHERE order_id = ? AND card_last_four = ?`, orderID, last4)
	return scanCardDeleteJob(row)
}

// ListCardDeleteJobs 列出任务（管理页用）。
func ListCardDeleteJobs(limit int, action string) ([]CardDeleteJob, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	const cols = `SELECT id, cdk_id, order_id, COALESCE(plan,''), COALESCE(account_email,''), COALESCE(card_last_four,''),
		       COALESCE(pm_id,''), COALESCE(pm_brand,''), COALESCE(cards_before,0),
		       COALESCE(action,''), COALESCE(reason,''), COALESCE(http_status,0), COALESCE(attempts,0),
		       COALESCE(setdefault_used,0), COALESCE(session_state,'unknown'),
		       COALESCE(session_checked_at,''), COALESCE(next_run_at,''), COALESCE(created_at,''), COALESCE(updated_at,'')
		FROM card_delete_jobs`
	var (
		rows *sql.Rows
		err  error
	)
	if strings.TrimSpace(action) != "" {
		rows, err = DB.Query(cols+` WHERE action = ? ORDER BY id DESC LIMIT ?`, action, limit)
	} else {
		rows, err = DB.Query(cols+` ORDER BY id DESC LIMIT ?`, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CardDeleteJob
	for rows.Next() {
		j, err := scanCardDeleteJob(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *j)
	}
	return out, rows.Err()
}

// ListDueCardDeleteJobs 取到期的 pending 任务。
func ListDueCardDeleteJobs(limit int) ([]CardDeleteJob, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := DB.Query(`
		SELECT id, cdk_id, order_id, COALESCE(plan,''), COALESCE(account_email,''), COALESCE(card_last_four,''),
		       COALESCE(pm_id,''), COALESCE(pm_brand,''), COALESCE(cards_before,0),
		       COALESCE(action,''), COALESCE(reason,''), COALESCE(http_status,0), COALESCE(attempts,0),
		       COALESCE(setdefault_used,0), COALESCE(session_state,'unknown'),
		       COALESCE(session_checked_at,''), COALESCE(next_run_at,''), COALESCE(created_at,''), COALESCE(updated_at,'')
		FROM card_delete_jobs
		WHERE action = 'pending' AND (next_run_at IS NULL OR next_run_at <= CURRENT_TIMESTAMP)
		ORDER BY next_run_at IS NULL, next_run_at
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CardDeleteJob
	for rows.Next() {
		j, err := scanCardDeleteJob(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *j)
	}
	return out, rows.Err()
}

// ListCardDeleteJobsByOrderIDs 批量取订单的解绑状态（对账页 join 用）。
func ListCardDeleteJobsByOrderIDs(orderIDs []int64) (map[int64]CardDeleteJob, error) {
	if DB == nil || len(orderIDs) == 0 {
		return map[int64]CardDeleteJob{}, nil
	}
	placeholders := make([]string, 0, len(orderIDs))
	args := make([]any, 0, len(orderIDs))
	for _, id := range orderIDs {
		if id <= 0 {
			continue
		}
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}
	if len(args) == 0 {
		return map[int64]CardDeleteJob{}, nil
	}
	rows, err := DB.Query(`
		SELECT id, cdk_id, order_id, COALESCE(plan,''), COALESCE(account_email,''), COALESCE(card_last_four,''),
		       COALESCE(pm_id,''), COALESCE(pm_brand,''), COALESCE(cards_before,0),
		       COALESCE(action,''), COALESCE(reason,''), COALESCE(http_status,0), COALESCE(attempts,0),
		       COALESCE(setdefault_used,0), COALESCE(session_state,'unknown'),
		       COALESCE(session_checked_at,''), COALESCE(next_run_at,''), COALESCE(created_at,''), COALESCE(updated_at,'')
		FROM card_delete_jobs WHERE order_id IN (`+strings.Join(placeholders, ",")+`)`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[int64]CardDeleteJob, len(args))
	for rows.Next() {
		j, err := scanCardDeleteJob(rows)
		if err != nil {
			return nil, err
		}
		out[j.OrderID] = *j
	}
	return out, rows.Err()
}

// CountCardDeleteJobsToday 同账号当日尝试次数（按脱敏邮箱）。
func CountCardDeleteJobsToday(accountEmail string) (int, error) {
	if DB == nil || strings.TrimSpace(accountEmail) == "" {
		return 0, nil
	}
	var n int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM card_delete_jobs
		WHERE account_email = ? AND date(created_at) = date('now')`, accountEmail).Scan(&n)
	return n, err
}

// ListCompletedWebhookPayloads 取最近 N 天内 gpt_direct.completed 的原始 payload（扫描补跑用）。
// since 非空时只返回 created_at >= since 的记录（格式 "2006-01-02 15:04:05"，UTC）。
func ListCompletedWebhookPayloads(days, limit int, since string) ([]string, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	if days <= 0 || days > 30 {
		days = 7
	}
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	query := `SELECT payload FROM webhook_events
		WHERE event_type = 'gpt_direct.completed'
		  AND created_at > datetime('now', ?)`
	args := []any{fmt.Sprintf("-%d day", days)}
	if s := strings.TrimSpace(since); s != "" {
		query += ` AND created_at >= ?`
		args = append(args, s)
	}
	query += ` ORDER BY id DESC LIMIT ?`
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCardDeleteJob(row rowScanner) (*CardDeleteJob, error) {
	var j CardDeleteJob
	var setDefault int
	err := row.Scan(
		&j.ID, &j.CDKID, &j.OrderID, &j.Plan, &j.AccountEmail, &j.CardLastFour,
		&j.PMID, &j.PMBrand, &j.CardsBefore,
		&j.Action, &j.Reason, &j.HTTPStatus, &j.Attempts,
		&setDefault, &j.SessionState,
		&j.SessionCheckedAt, &j.NextRunAt, &j.CreatedAt, &j.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	j.SetDefaultUsed = setDefault == 1
	return &j, nil
}

func nullIfEmpty(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

// SaveCardSnapshot 记录兑换前账号已有的 payment_method id（逗号分隔）。
func SaveCardSnapshot(cdkCode, pmIDs string) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	code := strings.ToUpper(strings.TrimSpace(cdkCode))
	if code == "" {
		return nil
	}
	_, err := DB.Exec(`
		INSERT INTO card_snapshots (cdk_code, pm_ids, captured_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(cdk_code) DO UPDATE SET
			pm_ids = excluded.pm_ids,
			captured_at = CURRENT_TIMESTAMP`, code, strings.TrimSpace(pmIDs))
	return err
}

// GetCardSnapshot 读取兑换前快照；ok=false 表示没有快照（老单或抓取失败）。
func GetCardSnapshot(cdkCode string) (pmIDs string, ok bool, err error) {
	if DB == nil {
		return "", false, fmt.Errorf("db not ready")
	}
	code := strings.ToUpper(strings.TrimSpace(cdkCode))
	if code == "" {
		return "", false, nil
	}
	var ids string
	err = DB.QueryRow(`SELECT COALESCE(pm_ids,'') FROM card_snapshots WHERE cdk_code = ?`, code).Scan(&ids)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return ids, true, nil
}

// PruneCardSnapshots 清理过旧快照（默认保留 30 天）。
func PruneCardSnapshots(days int) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	if days <= 0 {
		days = 30
	}
	_, err := DB.Exec(`DELETE FROM card_snapshots WHERE captured_at < datetime('now', ?)`,
		fmt.Sprintf("-%d day", days))
	return err
}
