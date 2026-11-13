package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/tuzi/cdk-recharge-system/internal/cardplatform"
)

// TestIssueAmbiguous 锁定"哪些发码失败值得用同一请求编号重放"的判定。
// 只有"结果不明"的错误才重放；参数/鉴权/限流属于确定未受理，重放没有意义。
func TestIssueAmbiguous(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"network", errors.New("dial tcp: i/o timeout"), true},
		{"request_timeout", &cardplatform.APIError{HTTPStatus: http.StatusRequestTimeout}, true},
		{"conflict_in_progress", &cardplatform.APIError{HTTPStatus: http.StatusConflict}, true},
		{"internal_error", &cardplatform.APIError{HTTPStatus: http.StatusInternalServerError}, true},
		{"bad_gateway", &cardplatform.APIError{HTTPStatus: http.StatusBadGateway}, true},
		{"service_unavailable", &cardplatform.APIError{HTTPStatus: http.StatusServiceUnavailable}, true},
		{"gateway_timeout", &cardplatform.APIError{HTTPStatus: http.StatusGatewayTimeout}, true},
		{"invalid_argument", &cardplatform.APIError{HTTPStatus: http.StatusBadRequest}, false},
		{"unauthorized", &cardplatform.APIError{HTTPStatus: http.StatusUnauthorized}, false},
		{"forbidden", &cardplatform.APIError{HTTPStatus: http.StatusForbidden}, false},
		{"rate_limited", &cardplatform.APIError{HTTPStatus: http.StatusTooManyRequests}, false},
		{"not_found", &cardplatform.APIError{HTTPStatus: http.StatusNotFound}, false},
	}
	for _, c := range cases {
		if got := issueAmbiguous(c.err); got != c.want {
			t.Fatalf("%s: issueAmbiguous=%v want %v", c.name, got, c.want)
		}
	}
}
