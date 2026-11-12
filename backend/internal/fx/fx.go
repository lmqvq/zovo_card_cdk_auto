// Package fx 汇率缓存：把订单金额按「北京时间当日汇率」折算成人民币。
//
// 数据源：fawazahmed0/currency-api（免费、无需 Key、URL 带日期即可取历史汇率）。
// 只缓存到本地 SQLite，页面刷新不会反复打外部接口。
package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tuzi/cdk-recharge-system/internal/db"
)

// cst 北京时间。用 FixedZone 避免依赖系统 tzdata。
var cst = time.FixedZone("CST", 8*3600)

const (
	quoteCNY     = "CNY"
	httpTimeout  = 8 * time.Second
	maxDates     = 120
	refreshEvery = 6 * time.Hour
)

// DateOf 把 unix 秒转成北京时间日期（YYYY-MM-DD）；0 返回空串。
func DateOf(unix int64) string {
	if unix <= 0 {
		return ""
	}
	return time.Unix(unix, 0).In(cst).Format("2006-01-02")
}

// Today 返回北京时间今天。
func Today() string {
	return time.Now().In(cst).Format("2006-01-02")
}

// Yesterday 返回北京时间昨天。
func Yesterday() string {
	return time.Now().In(cst).AddDate(0, 0, -1).Format("2006-01-02")
}

// Rate 取某天某货币对 CNY 的汇率；缓存优先，缺失时拉取，再失败回落最近可用汇率。
// ok=false 表示完全拿不到。source 里的 "|fallback" 表示用了兜底汇率。
func Rate(date, base string) (rate float64, source string, ok bool) {
	date = strings.TrimSpace(date)
	base = strings.ToUpper(strings.TrimSpace(base))
	if date == "" || base == "" {
		return 0, "", false
	}
	if r, src, hit, _ := db.GetFXRate(date, base, quoteCNY); hit {
		return r, src, true
	}
	if r, err := Ensure(date, base); err == nil {
		src := "api"
		if _, s, _, _ := db.GetFXRate(date, base, quoteCNY); s != "" {
			src = s
		}
		return r, src, true
	}
	if lr, _, lsrc, hit := db.LatestFXRate(base, quoteCNY); hit {
		return lr, lsrc + "|fallback", true
	}
	return 0, "", false
}

// Ensure 确保某天某货币的汇率已缓存。
func Ensure(date, base string) (float64, error) {
	date = strings.TrimSpace(date)
	base = strings.ToUpper(strings.TrimSpace(base))
	if date == "" || base == "" {
		return 0, fmt.Errorf("invalid args")
	}
	if r, _, hit, _ := db.GetFXRate(date, base, quoteCNY); hit {
		return r, nil
	}
	r, src, err := fetch(date, base)
	if err != nil {
		return 0, err
	}
	_ = db.SaveFXRate(date, base, quoteCNY, r, src)
	return r, nil
}

func fetch(date, base string) (float64, string, error) {
	lower := strings.ToLower(base)
	sources := []struct {
		url string
		src string
	}{
		{fmt.Sprintf("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@%s/v1/currencies/%s.json", date, lower), "jsdelivr"},
		{fmt.Sprintf("https://%s.currency-api.pages.dev/v1/currencies/%s.json", date, lower), "pages.dev"},
	}
	client := &http.Client{Timeout: httpTimeout}
	var lastErr error
	for _, s := range sources {
		req, err := http.NewRequest(http.MethodGet, s.url, nil)
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("User-Agent", "cdk-recharge-system/1.0")
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		_ = resp.Body.Close()
		if readErr != nil || resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("%s http=%d", s.src, resp.StatusCode)
			continue
		}
		var data map[string]any
		if json.Unmarshal(body, &data) != nil {
			lastErr = fmt.Errorf("%s 非 JSON", s.src)
			continue
		}
		cur, _ := data[lower].(map[string]any)
		if v, ok := cur["cny"].(float64); ok && v > 0 {
			return v, s.src, nil
		}
		lastErr = fmt.Errorf("%s 无 cny 字段", s.src)
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no source")
	}
	return 0, "", lastErr
}

// Backfill 补齐所有「有订单的日期 + 今天/昨天」的汇率；已有缓存跳过。
// 返回本次新拉取的条数。
func Backfill() (int, error) {
	dates, err := db.ListCompletedOrderDates(maxDates)
	if err != nil {
		return 0, err
	}
	set := map[string]bool{Today(): true, Yesterday(): true}
	for _, d := range dates {
		if strings.TrimSpace(d) != "" {
			set[strings.TrimSpace(d)] = true
		}
	}
	// 订单上出现过的货币 + USD 服务费
	bases := []string{"PHP", "USD", "EGP", "CLP"}
	n := 0
	for d := range set {
		for _, b := range bases {
			if _, _, ok, _ := db.GetFXRate(d, b, quoteCNY); ok {
				continue
			}
			if _, err := Ensure(d, b); err == nil {
				n++
			}
			time.Sleep(120 * time.Millisecond)
		}
	}
	return n, nil
}

// Start 启动后台补汇率：启动时补一次，之后每 6 小时补今天/昨天。
func Start(ctx context.Context) {
	go func() {
		if n, err := Backfill(); err != nil {
			log.Printf("[fx] backfill: %v", err)
		} else if n > 0 {
			log.Printf("[fx] backfilled %d rates", n)
		}
		ticker := time.NewTicker(refreshEvery)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if n, err := Backfill(); err == nil && n > 0 {
					log.Printf("[fx] refreshed %d rates", n)
				}
			}
		}
	}()
}
