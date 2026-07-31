package compat

import "testing"

func TestResponseHasExactlyOnePayload(t *testing.T) {
	if _, err := Success("one", map[string]any{"ok": true}); err != nil {
		t.Fatal(err)
	}
	if _, err := Failure("one", "bad request"); err != nil {
		t.Fatal(err)
	}
	if err := (Response{ID: "one"}).Validate(); err == nil {
		t.Fatal("expected neither payload error")
	}
	if err := (Response{ID: "one", Result: map[string]any{}, Error: "bad"}).Validate(); err == nil {
		t.Fatal("expected both payloads error")
	}
}
