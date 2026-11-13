package pro20x

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, h http.HandlerFunc) (*Client, func()) {
	t.Helper()
	srv := httptest.NewServer(h)
	return New(Config{BaseURL: srv.URL, APIKey: "test-key"}), srv.Close
}

func TestGetConfigParses(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Automate-Key") != "test-key" {
			t.Errorf("missing key header")
		}
		if r.URL.Path != "/config" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{
			"success": true,
			"planRegions": {"plus-monthly": ["PH","EG","JP","US","CL","SG"]},
			"paymentSources": [{"id":"card","label":"填卡支付"},{"id":"auto_card","label":"自动选卡"}],
			"pointsBilling": {"taskCosts": {"plus-monthly": 1}}
		}`))
	})
	defer done()

	cfg, err := c.GetConfig(context.Background())
	if err != nil {
		t.Fatalf("GetConfig: %v", err)
	}
	if !cfg.SupportsAutoCard() {
		t.Fatal("auto_card should be supported")
	}
	if cfg.TaskCost("plus-monthly") != 1 {
		t.Fatalf("task cost = %d", cfg.TaskCost("plus-monthly"))
	}
	if len(cfg.PlanRegions["plus-monthly"]) != 6 {
		t.Fatalf("regions = %v", cfg.PlanRegions["plus-monthly"])
	}
}

func TestAccountCheckPollsUntilDone(t *testing.T) {
	calls := 0
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/account-check":
			_, _ = w.Write([]byte(`{"success":true,"job":{"id":"ptj_1","type":"account_check","status":"queued"},"poll":"/account-check/jobs/ptj_1"}`))
		case strings.HasPrefix(r.URL.Path, "/account-check/jobs/"):
			calls++
			if calls < 2 {
				_, _ = w.Write([]byte(`{"success":true,"job":{"id":"ptj_1","status":"running"}}`))
				return
			}
			_, _ = w.Write([]byte(`{"success":true,"job":{"id":"ptj_1","status":"success","result":{"email":"a@b.com","selectedAccountId":"acc1","accounts":[{"accountId":"acc1","tier":"free","hasActiveSubscription":false}]}}}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})
	defer done()

	res, err := c.AccountCheck(context.Background(), map[string]any{"accessToken": "at"})
	if err != nil {
		t.Fatalf("AccountCheck: %v", err)
	}
	if res.SelectedAccountID != "acc1" || len(res.Accounts) != 1 {
		t.Fatalf("result = %+v", res)
	}
}

func TestCreateTaskSendsFields(t *testing.T) {
	var got map[string]any
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_ = json.NewDecoder(r.Body).Decode(&got)
		_, _ = w.Write([]byte(`{"success":true,"id":"AUTO-1","clientOrderId":"FB-1","status":"queued"}`))
	})
	defer done()

	yes := true
	task, err := c.CreateTask(context.Background(), CreateTaskRequest{
		ClientOrderID:     "FB-1",
		PlanID:            "plus-monthly",
		CheckoutCountry:   "PH",
		PaymentSource:     "auto_card",
		CardPlatformID:    "zp_1",
		CancelAutoRenew:   &yes,
		DeletePaymentCard: &yes,
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if task.ID != "AUTO-1" || task.Status != "queued" {
		t.Fatalf("task = %+v", task)
	}
	if got["deletePaymentCard"] != true {
		t.Fatalf("deletePaymentCard not sent: %v", got)
	}
	if got["paymentSource"] != "auto_card" {
		t.Fatalf("paymentSource = %v", got["paymentSource"])
	}
}

func TestGetTaskAndTerminal(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":true,"id":"AUTO-2","status":"succeeded","gptPaymentMethodCleanup":{"status":"deleted_verified","verified":true}}`))
	})
	defer done()

	task, err := c.GetTask(context.Background(), "AUTO-2")
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if !IsTerminal(task.Status) {
		t.Fatal("succeeded should be terminal")
	}
	if !strings.Contains(string(task.Cleanup), "deleted_verified") {
		t.Fatalf("cleanup = %s", task.Cleanup)
	}
}

func TestAPIErrorOnFailure(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
		_, _ = w.Write([]byte(`{"success":false,"code":"automate_points_insufficient","message":"not enough points"}`))
	})
	defer done()

	_, err := c.CreateTask(context.Background(), CreateTaskRequest{ClientOrderID: "x", PlanID: "plus-monthly"})
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("err = %v", err)
	}
	if apiErr.Code != "automate_points_insufficient" || apiErr.HTTPStatus != 402 {
		t.Fatalf("apiErr = %+v", apiErr)
	}
}

func TestMissingKey(t *testing.T) {
	c := New(Config{BaseURL: "http://127.0.0.1:1"})
	_, err := c.GetConfig(context.Background())
	if _, ok := err.(*APIError); !ok {
		t.Fatalf("want APIError, got %v", err)
	}
}
