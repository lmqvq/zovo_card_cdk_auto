package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// 自建码字母表：去掉易混字符 0/O/1/I。
const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// GenerateCDKCode 生成 ZC-XXXXXXXX-XXXX 格式的本地码（12 位随机，≈ 59 bit）。
func GenerateCDKCode() (string, error) {
	const n = 12
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	chars := make([]byte, n)
	for i := 0; i < n; i++ {
		chars[i] = codeAlphabet[int(buf[i])%len(codeAlphabet)]
	}
	return "ZC-" + string(chars[:8]) + "-" + string(chars[8:]), nil
}

// LocalCDK 是一条自建卡密。
type LocalCDK struct {
	Code        string
	Prefix      string
	Plan        string
	Country     string
	Status      string
	BatchNo     string
	AgentID     string
	PriceMinor  int64
	BoundEmail  string
	ExpiresAt   string
	CreatedAt   string
}

// InsertLocalCDK 写入一条自建码（source=local，upstream_id=0）。
func InsertLocalCDK(c LocalCDK) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	status := strings.TrimSpace(c.Status)
	if status == "" {
		status = "unused"
	}
	var exp any
	if strings.TrimSpace(c.ExpiresAt) != "" {
		exp = strings.TrimSpace(c.ExpiresAt)
	}
	_, err := DB.Exec(`
		INSERT INTO cardplatform_cdk_codes
			(upstream_id, code, code_prefix, plan, fee_amount_minor, status, payment_country,
			 source, expires_at, batch_no, agent_id, price_minor, bound_email)
		VALUES (0, ?, ?, ?, 0, ?, ?, 'local', ?, ?, ?, ?, ?)`,
		c.Code, c.Prefix, c.Plan, status, c.Country,
		exp, c.BatchNo, c.AgentID, c.PriceMinor, c.BoundEmail)
	return err
}

// GenerateLocalCDKs 批量生成并入库，返回生成成功的码。
// 冲突（极小概率）时自动重试一次。
func GenerateLocalCDKs(plan, country, batch, agent string, count int, priceMinor int64, validDays int) ([]LocalCDK, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	if count <= 0 || count > 1000 {
		count = 1
	}
	var expires string
	if validDays > 0 {
		expires = time.Now().UTC().AddDate(0, 0, validDays).Format("2006-01-02 15:04:05")
	}
	out := make([]LocalCDK, 0, count)
	for i := 0; i < count; i++ {
		var code string
		var err error
		ok := false
		for attempt := 0; attempt < 3; attempt++ {
			code, err = GenerateCDKCode()
			if err != nil {
				return out, err
			}
			c := LocalCDK{
				Code:       code,
				Prefix:     code[:11],
				Plan:       plan,
				Country:    country,
				BatchNo:    batch,
				AgentID:    agent,
				PriceMinor: priceMinor,
				ExpiresAt:  expires,
			}
			if err := InsertLocalCDK(c); err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "unique") {
					continue // 撞码，重生成
				}
				return out, err
			}
			out = append(out, c)
			ok = true
			break
		}
		if !ok {
			return out, fmt.Errorf("生成卡密失败（连续撞码）")
		}
	}
	return out, nil
}

// 自建码兑换时的原子状态流转：
//
//	unused/''  --ClaimLocalCDK-->  consuming  --ConsumeLocalCDK-->  consumed
//	                                   └------ReleaseLocalCDK----->  unused
//
// Claim 用条件 UPDATE 实现抢占，保证同一张码只会被一个请求拿到。

// ClaimLocalCDK 原子占用卡密；false 表示已被别人占用或已使用。
func ClaimLocalCDK(code string) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("db not ready")
	}
	c := normalizeCDKCode(code)
	res, err := DB.Exec(`
		UPDATE cardplatform_cdk_codes
		SET status = 'consuming'
		WHERE code = ? AND COALESCE(source,'zovo') = 'local'
		  AND TRIM(COALESCE(status,'')) IN ('', 'unused')`, c)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}

// ConsumeLocalCDK 兑换成功：consuming → consumed，并记录绑定邮箱。
func ConsumeLocalCDK(code, email string) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	c := normalizeCDKCode(code)
	_, err := DB.Exec(`
		UPDATE cardplatform_cdk_codes
		SET status = 'consumed', bound_email = ?
		WHERE code = ? AND COALESCE(source,'zovo') = 'local'`, strings.TrimSpace(email), c)
	return err
}

// ReleaseLocalCDK 兑换失败：consuming → unused，允许重试。
func ReleaseLocalCDK(code string) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	c := normalizeCDKCode(code)
	_, err := DB.Exec(`
		UPDATE cardplatform_cdk_codes
		SET status = 'unused'
		WHERE code = ? AND COALESCE(source,'zovo') = 'local' AND status = 'consuming'`, c)
	return err
}

// UpdateLocalCDKStatus 管理端手动改状态（禁用/启用）。
func UpdateLocalCDKStatus(code, status string) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	c := normalizeCDKCode(code)
	_, err := DB.Exec(`
		UPDATE cardplatform_cdk_codes SET status = ?
		WHERE code = ? AND COALESCE(source,'zovo') = 'local'`, strings.TrimSpace(status), c)
	return err
}

// LookupLocalCDK 按完整码查自建码；不存在返回 ok=false。
func LookupLocalCDK(code string) (*LocalCDK, bool, error) {
	if DB == nil {
		return nil, false, fmt.Errorf("db not ready")
	}
	c := normalizeCDKCode(code)
	if c == "" {
		return nil, false, nil
	}
	var (
		row                                    LocalCDK
		expires, created                       sql.NullString
		price                                  sql.NullInt64
	)
	err := DB.QueryRow(`
		SELECT code, COALESCE(code_prefix,''), COALESCE(plan,''), COALESCE(payment_country,''),
		       COALESCE(status,''), COALESCE(batch_no,''), COALESCE(agent_id,''),
		       COALESCE(price_minor,0), COALESCE(bound_email,''), expires_at, created_at
		FROM cardplatform_cdk_codes
		WHERE code = ? AND COALESCE(source,'zovo') = 'local'`, c).
		Scan(&row.Code, &row.Prefix, &row.Plan, &row.Country, &row.Status,
			&row.BatchNo, &row.AgentID, &price, &row.BoundEmail, &expires, &created)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	row.PriceMinor = price.Int64
	row.ExpiresAt = expires.String
	row.CreatedAt = created.String
	return &row, true, nil
}

// MarkLocalCDKUsed 兑换成功后把自建码标记为已用。
func MarkLocalCDKUsed(code, email string) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	c := normalizeCDKCode(code)
	_, err := DB.Exec(`
		UPDATE cardplatform_cdk_codes
		SET status = 'consumed', bound_email = ?
		WHERE code = ? AND COALESCE(source,'zovo') = 'local'`, strings.TrimSpace(email), c)
	return err
}

// ListLocalCDKs 列出自建码（管理页）。
func ListLocalCDKs(limit int, status, batch string) ([]LocalCDK, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	q := `SELECT code, COALESCE(code_prefix,''), COALESCE(plan,''), COALESCE(payment_country,''),
		       COALESCE(status,''), COALESCE(batch_no,''), COALESCE(agent_id,''),
		       COALESCE(price_minor,0), COALESCE(bound_email,''),
		       COALESCE(expires_at,''), COALESCE(created_at,'')
		FROM cardplatform_cdk_codes
		WHERE COALESCE(source,'zovo') = 'local'`
	args := []any{}
	if s := strings.TrimSpace(status); s != "" {
		q += ` AND COALESCE(status,'') = ?`
		args = append(args, s)
	}
	if b := strings.TrimSpace(batch); b != "" {
		q += ` AND COALESCE(batch_no,'') = ?`
		args = append(args, b)
	}
	q += ` ORDER BY created_at DESC LIMIT ?`
	args = append(args, limit)

	rows, err := DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LocalCDK
	for rows.Next() {
		var c LocalCDK
		if err := rows.Scan(&c.Code, &c.Prefix, &c.Plan, &c.Country, &c.Status,
			&c.BatchNo, &c.AgentID, &c.PriceMinor, &c.BoundEmail, &c.ExpiresAt, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// CountLocalCDKs 统计自建码数量（按状态过滤，空表示全部）。
func CountLocalCDKs(status string) (int, error) {
	if DB == nil {
		return 0, fmt.Errorf("db not ready")
	}
	q := `SELECT COUNT(*) FROM cardplatform_cdk_codes WHERE COALESCE(source,'zovo') = 'local'`
	args := []any{}
	if s := strings.TrimSpace(status); s != "" {
		q += ` AND COALESCE(status,'') = ?`
		args = append(args, s)
	}
	var n int
	err := DB.QueryRow(q, args...).Scan(&n)
	return n, err
}
