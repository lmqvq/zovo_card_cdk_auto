package pro20x

import "testing"

func TestUnwrapTaskVariants(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
		bad  bool
	}{
		{"顶层 Task", `{"id":"AUTO-1","status":"queued"}`, "AUTO-1", false},
		{"task 包裹", `{"success":true,"task":{"id":"AUTO-2","status":"queued"}}`, "AUTO-2", false},
		{"job 包裹", `{"success":true,"job":{"id":"ptj_3","status":"queued"}}`, "ptj_3", false},
		{"data.task 包裹", `{"success":true,"data":{"task":{"id":"AUTO-4"}}}`, "AUTO-4", false},
		{"taskId 扁平", `{"success":true,"taskId":"AUTO-5"}`, "AUTO-5", false},
		{"data.id 扁平", `{"success":true,"data":{"id":"AUTO-6"}}`, "AUTO-6", false},
		{"无 ID", `{"success":true}`, "", true},
		{"非 JSON", `<html>404</html>`, "", true},
	}
	for _, c := range cases {
		got, err := unwrapTask([]byte(c.raw))
		if c.bad {
			if err == nil {
				t.Errorf("%s: want error, got %+v", c.name, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error %v", c.name, err)
			continue
		}
		if got.ID != c.want {
			t.Errorf("%s: id=%q want %q", c.name, got.ID, c.want)
		}
	}
}

func TestUnwrapTaskKeepsFields(t *testing.T) {
	raw := `{"success":true,"task":{"id":"AUTO-9","status":"failed","stage":"card","errorCode":"card_recharge_failed","gptPaymentMethodCleanup":{"status":"failed"}}}`
	tk, err := unwrapTask([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if tk.Status != "failed" || tk.ErrorCode != "card_recharge_failed" {
		t.Fatalf("fields lost: %+v", tk)
	}
	if len(tk.Cleanup) == 0 {
		t.Fatal("cleanup lost")
	}
}
