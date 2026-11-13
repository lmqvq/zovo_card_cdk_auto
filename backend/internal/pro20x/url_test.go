package pro20x

import "testing"

func TestResolveURL(t *testing.T) {
	c := New(Config{BaseURL: "https://pro20x.com/api/v3/automate", APIKey: "k"})
	cases := []struct{ in, want string }{
		// 站点绝对路径：只补 scheme+host
		{"/api/v3/automate/account-check/jobs/ptj_1", "https://pro20x.com/api/v3/automate/account-check/jobs/ptj_1"},
		// 接口相对路径：拼完整 BaseURL
		{"/account-check", "https://pro20x.com/api/v3/automate/account-check"},
		{"account-check", "https://pro20x.com/api/v3/automate/account-check"},
		{"/tasks/AUTO-1", "https://pro20x.com/api/v3/automate/tasks/AUTO-1"},
		// 完整 URL：原样返回
		{"https://pro20x.com/api/v3/automate/tasks/x", "https://pro20x.com/api/v3/automate/tasks/x"},
	}
	for _, tc := range cases {
		if got := c.resolveURL(tc.in); got != tc.want {
			t.Errorf("resolveURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
