package gptcard

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func fakeJWT(accountID string) string {
	h := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	p, _ := json.Marshal(map[string]any{
		"https://api.openai.com/auth": map[string]string{"chatgpt_account_id": accountID},
	})
	pl := base64.RawURLEncoding.EncodeToString(p)
	return h + "." + pl + ".sig"
}

func TestParseSessionJSON(t *testing.T) {
	raw := `{"accessToken":"at-1","sessionToken":"st-1","account":{"id":"acc-1"}}`
	s, err := ParseSession(raw)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if s.AccessToken != "at-1" || s.SessionToken != "st-1" || s.AccountID != "acc-1" {
		t.Fatalf("parsed = %+v", s)
	}
}

func TestParseSessionBareAccessToken(t *testing.T) {
	tok := fakeJWT("acc-jwt")
	s, err := ParseSession(tok)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if s.AccessToken != tok {
		t.Fatalf("access token mismatch")
	}
	if s.AccountID != "acc-jwt" {
		t.Fatalf("account id = %q, want acc-jwt", s.AccountID)
	}
}

func TestParseSessionEmpty(t *testing.T) {
	if _, err := ParseSession("  "); err == nil {
		t.Fatal("want error for empty session")
	}
}

func TestParseSessionBareSessionToken(t *testing.T) {
	tok := "a.b.c.d.e" // JWE 5 段
	s, err := ParseSession(tok)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if s.SessionToken != tok || s.AccessToken != "" {
		t.Fatalf("parsed = %+v", s)
	}
}

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	c, err := NewClient()
	if err != nil {
		srv.Close()
		t.Fatalf("NewClient: %v", err)
	}
	c.SetBaseURL(srv.URL)
	return c, srv.Close
}

func TestListParsesCards(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/backend-api/payments/payment_methods" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("account_id"); got != "acc-1" {
			t.Errorf("account_id = %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer at-1" {
			t.Errorf("auth = %q", got)
		}
		if got := r.Header.Get("chatgpt-account-id"); got != "acc-1" {
			t.Errorf("chatgpt-account-id = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"payment_methods":[
			{"id":"pm_1","card":{"brand":"mastercard","last4":"3948","exp_month":10,"exp_year":2027},"created":1700000000},
			{"id":"pm_2","card":{"brand":"mastercard","last4":"4880"}}
		],"default_payment_method_id":"pm_1"}`))
	})
	defer done()

	list, err := c.List(&Session{AccessToken: "at-1", AccountID: "acc-1"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list.Cards) != 2 {
		t.Fatalf("cards = %d, want 2", len(list.Cards))
	}
	if list.Cards[0].Last4 != "3948" || list.Cards[0].Brand != "mastercard" || !list.Cards[0].IsDefault {
		t.Fatalf("card0 = %+v", list.Cards[0])
	}
	if list.Cards[0].ExpYear != 2027 || list.Cards[0].Created != 1700000000 {
		t.Fatalf("card0 fields = %+v", list.Cards[0])
	}
	if list.Cards[1].IsDefault {
		t.Fatalf("card1 should not be default")
	}
}

func TestListUnauthorized(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"bad token"}`))
	})
	defer done()

	_, err := c.List(&Session{AccessToken: "x", AccountID: "a"})
	apiErr, ok := err.(*APIError)
	if !ok || !apiErr.Unauthorized() || apiErr.Status != 401 {
		t.Fatalf("err = %v", err)
	}
}

func TestDeleteSendsDeleteRequest(t *testing.T) {
	var gotPath, gotAccount string
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAccount = r.URL.Query().Get("account_id")
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		_, _ = w.Write([]byte(`{"success":true}`))
	})
	defer done()

	st, ok, err := c.Delete(&Session{AccessToken: "at", AccountID: "acc"}, "pm_9")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if st != 200 || !ok {
		t.Fatalf("status=%d ok=%v", st, ok)
	}
	if gotPath != "/backend-api/payments/payment_method/pm_9" {
		t.Fatalf("path = %q", gotPath)
	}
	if gotAccount != "acc" {
		t.Fatalf("account = %q", gotAccount)
	}
}

func TestSetDefaultBodyTemplate(t *testing.T) {
	var body string
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/backend-api/payments/payment_method/default" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		_, _ = w.Write([]byte(`{"success":true}`))
	})
	defer done()

	opt := DefaultSetDefaultOptions()
	st, err := c.SetDefault(&Session{AccessToken: "at", AccountID: "acc"}, "pm_2", opt)
	if err != nil || st != 200 {
		t.Fatalf("SetDefault st=%d err=%v", st, err)
	}
	if !strings.Contains(body, `"payment_method_id":"pm_2"`) || !strings.Contains(body, `"account_id":"acc"`) {
		t.Fatalf("body = %s", body)
	}
}

func TestSessionTokenCookiesChunked(t *testing.T) {
	tok := strings.Repeat("x", sessionChunk*2+10)
	cookies := sessionTokenCookies(tok)
	if len(cookies) != 3 {
		t.Fatalf("cookies = %d, want 3", len(cookies))
	}
	if cookies[0].Name != "__Secure-next-auth.session-token.0" {
		t.Fatalf("first cookie = %q", cookies[0].Name)
	}
	var joined string
	for _, c := range cookies {
		joined += c.Value
	}
	if joined != tok {
		t.Fatal("chunked cookies do not reassemble")
	}
	// 短 token 不分片，直接原名
	short := sessionTokenCookies("abc")
	if len(short) != 1 || short[0].Name != "__Secure-next-auth.session-token" {
		t.Fatalf("short cookie = %+v", short)
	}
}

func TestPickTarget(t *testing.T) {
	list := &CardList{Cards: []Card{
		{ID: "pm_1", Last4: "3948"},
		{ID: "pm_2", Last4: "4880"},
	}}
	target, reason := pickTarget(list, "3948")
	if reason != "" || target == nil || target.ID != "pm_1" {
		t.Fatalf("target=%v reason=%s", target, reason)
	}
	if _, reason := pickTarget(list, "0000"); reason != "not_found" {
		t.Fatalf("reason = %s", reason)
	}
	dup := &CardList{Cards: []Card{
		{ID: "pm_1", Last4: "1111"},
		{ID: "pm_2", Last4: "1111"},
	}}
	if _, reason := pickTarget(dup, "1111"); reason != "ambiguous" {
		t.Fatalf("reason = %s, want ambiguous", reason)
	}
	if pickOther(list, "pm_1") == nil || pickOther(list, "pm_1").ID != "pm_2" {
		t.Fatal("pickOther failed")
	}
}

func TestCardLabel(t *testing.T) {
	c := Card{Brand: "mastercard", Last4: "5794"}
	if got := c.Label(); got != "mastercard •••• 5794" {
		t.Fatalf("label = %q", got)
	}
}
