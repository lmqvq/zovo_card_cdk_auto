package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// Pro20xOrder 是一条走 Pro20x 通道的兑换单。
type Pro20xOrder struct {
	ID              int64
	RedeemToken     string
	CDKCode         string
	Plan            string
	Country         string
	CreditsQuantity int
	ClientOrderID   string
	TaskID          string
	Status          string
	Stage           string
	AccountEmail    string
	ErrorCode       string
	CleanupStatus   string
	CostPoints      int
	CreatedAt       string
	UpdatedAt       string
}

const pro20xCols = `id, redeem_token, cdk_code, COALESCE(plan,''), COALESCE(country,''),
	COALESCE(credits_quantity,0), client_order_id, COALESCE(task_id,''),
	COALESCE(status,''), COALESCE(stage,''), COALESCE(account_email,''),
	COALESCE(error_code,''), COALESCE(cleanup_status,''), COALESCE(cost_points,0),
	COALESCE(created_at,''), COALESCE(updated_at,'')`

// InsertPro20xOrder 创建兑换单；UNIQUE 冲突返回错误。
func InsertPro20xOrder(o Pro20xOrder) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("db not ready")
	}
	if strings.TrimSpace(o.Status) == "" {
		o.Status = "pending"
	}
	res, err := DB.Exec(`
		INSERT INTO pro20x_orders
			(redeem_token, cdk_code, plan, country, credits_quantity, client_order_id, status, stage, account_email)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.RedeemToken, o.CDKCode, o.Plan, o.Country, o.CreditsQuantity,
		o.ClientOrderID, o.Status, o.Stage, o.AccountEmail)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdatePro20xOrder 按 id 更新状态字段。
func UpdatePro20xOrder(o Pro20xOrder) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	_, err := DB.Exec(`
		UPDATE pro20x_orders SET
			task_id = ?, status = ?, stage = ?, account_email = ?,
			error_code = ?, cleanup_status = ?, cost_points = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		o.TaskID, o.Status, o.Stage, o.AccountEmail,
		o.ErrorCode, o.CleanupStatus, o.CostPoints, o.ID)
	return err
}

func scanPro20xOrder(row interface{ Scan(...any) error }) (*Pro20xOrder, error) {
	var o Pro20xOrder
	err := row.Scan(&o.ID, &o.RedeemToken, &o.CDKCode, &o.Plan, &o.Country,
		&o.CreditsQuantity, &o.ClientOrderID, &o.TaskID,
		&o.Status, &o.Stage, &o.AccountEmail,
		&o.ErrorCode, &o.CleanupStatus, &o.CostPoints,
		&o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetPro20xOrderByToken 按兑换令牌查。
func GetPro20xOrderByToken(token string) (*Pro20xOrder, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	row := DB.QueryRow(`SELECT `+pro20xCols+` FROM pro20x_orders WHERE redeem_token = ?`, strings.TrimSpace(token))
	return scanPro20xOrder(row)
}

// GetPro20xOrderByCDK 按卡密取最近一条。
func GetPro20xOrderByCDK(code string) (*Pro20xOrder, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	row := DB.QueryRow(`SELECT `+pro20xCols+` FROM pro20x_orders WHERE cdk_code = ? ORDER BY id DESC LIMIT 1`, strings.TrimSpace(code))
	return scanPro20xOrder(row)
}

// ClaimPro20xOrder 原子把会话从 preflight 推进到 redeeming；
// false 表示会话已不在 preflight（重复提交/并发）。
func ClaimPro20xOrder(token string) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("db not ready")
	}
	res, err := DB.Exec(`
		UPDATE pro20x_orders SET status = 'redeeming', updated_at = CURRENT_TIMESTAMP
		WHERE redeem_token = ? AND status = 'preflight'`, strings.TrimSpace(token))
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}

// ReleasePro20xOrder 兑换失败时把会话退回 preflight，允许重试。
func ReleasePro20xOrder(token string) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	_, err := DB.Exec(`
		UPDATE pro20x_orders SET status = 'preflight', updated_at = CURRENT_TIMESTAMP
		WHERE redeem_token = ? AND status = 'redeeming'`, strings.TrimSpace(token))
	return err
}

// ListPro20xOrders 列出最近的兑换单（管理页用）。
func ListPro20xOrders(limit int) ([]Pro20xOrder, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := DB.Query(`SELECT `+pro20xCols+` FROM pro20x_orders ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Pro20xOrder
	for rows.Next() {
		o, err := scanPro20xOrder(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *o)
	}
	return out, rows.Err()
}
