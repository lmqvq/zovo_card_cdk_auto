package cardplatform

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXPremiumCatalogueIsProductScoped(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/openapi/v1/gpt-direct/plans" || r.URL.Query().Get("product") != "x" {
			t.Errorf("wrong product request %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0,"data":{"version":1,"plans":{"x_basic_monthly":{"enabled":true,"currency":"JPY","serviceFeeUsdMinor":15}},"registry":[{"key":"basic_monthly","label":"X Basic","flow":"direct"}]}}`))
	}))
	defer s.Close()
	c := New(Config{SiteBase: s.URL, APIKey: "fixture-key"})
	plans, err := c.GetPlans(context.Background(), "x")
	if err != nil {
		t.Fatal(err)
	}
	rows := plans.SellablePlans()
	if len(rows) != 1 || rows[0].Key != "basic_monthly" || rows[0].ServiceFeeUsdMinor != 15 {
		t.Fatalf("bad X catalogue %+v", rows)
	}
	if len(plans.PaymentRegions) != 1 || plans.PaymentRegions[0].Country != "JP" {
		t.Fatal("wrong X region")
	}
}
func TestXCredentialsAreNotInvoiceSessions(t *testing.T) {
	if !IsXPremiumCredential(`{"auth_token":"fixture","ct0":"fixture","billing_email":"fixture@example.test"}`) {
		t.Fatal("X credential stored as GPT session")
	}
	if IsXPremiumCredential(`{"sessionToken":"fixture"}`) {
		t.Fatal("GPT credential misclassified")
	}
	for _, plan := range []string{"basic_monthly", "premium_yearly", "premium_plus_yearly", "x_premium_monthly"} {
		if !IsXPremiumPlan(plan) {
			t.Fatal(plan)
		}
	}
	if IsXPremiumPlan("plus") || IsXPremiumPlan("grok_plus_monthly") {
		t.Fatal("cross product plan")
	}
}
