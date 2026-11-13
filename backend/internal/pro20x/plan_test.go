package pro20x

import (
	"encoding/json"
	"testing"
)

func TestMapPlan(t *testing.T) {
	cases := []struct {
		in       string
		planID   string
		credits  int
		ok       bool
	}{
		{"plus", "plus-monthly", 0, true},
		{"PLUS", "plus-monthly", 0, true},
		{"pro_5x", "pro5x-monthly", 0, true},
		{"pro_20x", "pro20x-direct-monthly", 0, true},
		{"credit250", "credits", 250, true},
		{"credit500", "credits", 500, true},
		{"credit1000", "credits", 1000, true},
		{"unknown-plan", "", 0, false},
		{"", "", 0, false},
	}
	for _, c := range cases {
		m, ok := MapPlan(c.in)
		if ok != c.ok || m.PlanID != c.planID || m.CreditsQuantity != c.credits {
			t.Errorf("MapPlan(%q) = %+v,%v; want %s,%d,%v", c.in, m, ok, c.planID, c.credits, c.ok)
		}
	}
}

func TestNormalizeCountry(t *testing.T) {
	if got := NormalizeCountry("ph", "US"); got != "PH" {
		t.Fatalf("ph -> %s", got)
	}
	if got := NormalizeCountry("XX", "PH"); got != "PH" {
		t.Fatalf("unsupported should fallback, got %s", got)
	}
	if got := NormalizeCountry("", "EG"); got != "EG" {
		t.Fatalf("empty should fallback, got %s", got)
	}
}

func TestBuildCreateRequestAutoCard(t *testing.T) {
	req, err := BuildCreateRequest(BuildOptions{ClientOrderID: "FB-1001", CDKPlan: "plus", Country: "PH", FallbackCountry: "PH", AuthSession: map[string]any{"accessToken": "at"}, PaymentSource: "auto_card", CardPlatformID: "zp_1"})
	if err != nil {
		t.Fatalf("BuildCreateRequest: %v", err)
	}
	if req.PlanID != "plus-monthly" || req.CheckoutCountry != "PH" {
		t.Fatalf("req = %+v", req)
	}
	if req.DeletePaymentCard == nil || !*req.DeletePaymentCard {
		t.Fatal("deletePaymentCard must be true")
	}
	if req.CancelAutoRenew == nil || !*req.CancelAutoRenew {
		t.Fatal("cancelAutoRenew must be true")
	}
	if req.DetectPlusPromo == nil || !*req.DetectPlusPromo {
		t.Fatal("plus 应开启 detectPlusPromo")
	}
	if req.PaymentSource != "auto_card" || req.CardPlatformID != "zp_1" {
		t.Fatalf("payment = %s/%s", req.PaymentSource, req.CardPlatformID)
	}
}

func TestBuildCreateRequestCredits(t *testing.T) {
	req, err := BuildCreateRequest(BuildOptions{ClientOrderID: "FB-1002", CDKPlan: "credit500", FallbackCountry: "PH", PaymentSource: "wallet"})
	if err != nil {
		t.Fatalf("BuildCreateRequest: %v", err)
	}
	if req.PlanID != "credits" || req.CreditsQuantity != 500 {
		t.Fatalf("req = %+v", req)
	}
	if req.DetectPlusPromo != nil {
		t.Fatal("credits 不应开启 detectPlusPromo")
	}
	if req.CheckoutCountry != "PH" {
		t.Fatalf("empty country should fallback to PH, got %s", req.CheckoutCountry)
	}
}

func TestBuildCreateRequestUnsupported(t *testing.T) {
	_, err := BuildCreateRequest(BuildOptions{ClientOrderID: "FB-1003", CDKPlan: "nope", Country: "PH", FallbackCountry: "PH", PaymentSource: "auto_card"})
	if err == nil {
		t.Fatal("want error for unmapped plan")
	}
}

func TestCreateTaskRequestJSONOmitsEmpty(t *testing.T) {
	yes := true
	req := CreateTaskRequest{
		ClientOrderID:     "FB-1",
		PlanID:            "plus-monthly",
		CheckoutCountry:   "PH",
		PaymentSource:     "auto_card",
		DeletePaymentCard: &yes,
	}
	b, _ := json.Marshal(req)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if _, ok := m["card"]; ok {
		t.Fatal("card should be omitted")
	}
	if _, ok := m["creditsQuantity"]; ok {
		t.Fatal("creditsQuantity should be omitted")
	}
	if m["deletePaymentCard"] != true {
		t.Fatal("deletePaymentCard missing")
	}
}
