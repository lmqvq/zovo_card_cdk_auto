package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// GetFXRate 读某天某货币对 CNY 的汇率缓存。
func GetFXRate(date, base, quote string) (rate float64, source string, ok bool, err error) {
	if DB == nil {
		return 0, "", false, fmt.Errorf("db not ready")
	}
	date = strings.TrimSpace(date)
	base = strings.ToUpper(strings.TrimSpace(base))
	quote = strings.ToUpper(strings.TrimSpace(quote))
	if date == "" || base == "" || quote == "" {
		return 0, "", false, nil
	}
	err = DB.QueryRow(`SELECT rate, COALESCE(source,'') FROM fx_rates WHERE rate_date = ? AND base = ? AND quote = ?`,
		date, base, quote).Scan(&rate, &source)
	if err == sql.ErrNoRows {
		return 0, "", false, nil
	}
	if err != nil {
		return 0, "", false, err
	}
	return rate, source, true, nil
}

// LatestFXRate 取该货币最近一条缓存（当天数据缺失时兜底）。
func LatestFXRate(base, quote string) (rate float64, date string, source string, ok bool) {
	if DB == nil {
		return 0, "", "", false
	}
	base = strings.ToUpper(strings.TrimSpace(base))
	quote = strings.ToUpper(strings.TrimSpace(quote))
	err := DB.QueryRow(`SELECT rate, rate_date, COALESCE(source,'') FROM fx_rates
		WHERE base = ? AND quote = ? ORDER BY rate_date DESC LIMIT 1`, base, quote).
		Scan(&rate, &date, &source)
	if err != nil {
		return 0, "", "", false
	}
	return rate, date, source, true
}

// SaveFXRate 写入/更新一条汇率缓存。
func SaveFXRate(date, base, quote string, rate float64, source string) error {
	if DB == nil {
		return fmt.Errorf("db not ready")
	}
	date = strings.TrimSpace(date)
	base = strings.ToUpper(strings.TrimSpace(base))
	quote = strings.ToUpper(strings.TrimSpace(quote))
	if date == "" || base == "" || quote == "" || rate <= 0 {
		return nil
	}
	_, err := DB.Exec(`
		INSERT INTO fx_rates (rate_date, base, quote, rate, source, fetched_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(rate_date, base, quote) DO UPDATE SET
			rate = excluded.rate,
			source = excluded.source,
			fetched_at = CURRENT_TIMESTAMP`, date, base, quote, rate, source)
	return err
}

// ListCompletedOrderDates 已完成的 GPT 订单按「北京时间日期」去重，最近在前。
// 用于只补齐真正有订单的日期，避免无意义的外部请求。
func ListCompletedOrderDates(limit int) ([]string, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	if limit <= 0 || limit > 400 {
		limit = 400
	}
	rows, err := DB.Query(`
		SELECT DISTINCT date(datetime(json_extract(payload,'$.completed_at'),'unixepoch','+8 hours')) AS d
		FROM webhook_events
		WHERE event_type = 'gpt_direct.completed'
		  AND json_extract(payload,'$.completed_at') IS NOT NULL
		ORDER BY d DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var d sql.NullString
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		if d.Valid && strings.TrimSpace(d.String) != "" {
			out = append(out, d.String)
		}
	}
	return out, rows.Err()
}

// ListCompletedPayloadsOnDate 取「北京时间某一天」完成的订单 payload。
func ListCompletedPayloadsOnDate(date string) ([]string, error) {
	if DB == nil {
		return nil, fmt.Errorf("db not ready")
	}
	date = strings.TrimSpace(date)
	if date == "" {
		return nil, nil
	}
	rows, err := DB.Query(`
		SELECT payload FROM webhook_events
		WHERE event_type = 'gpt_direct.completed'
		  AND date(datetime(json_extract(payload,'$.completed_at'),'unixepoch','+8 hours')) = ?
		ORDER BY id DESC`, date)
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
