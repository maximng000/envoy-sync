package envfile

import (
	"testing"
	"time"
)

var pinTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func pinEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestPin_AllKeys(t *testing.T) {
	res, err := Pin(pinEntries(), nil, pinTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 3 {
		t.Fatalf("expected 3 pinned, got %d", len(res.Pinned))
	}
	if len(res.Skipped) != 0 {
		t.Fatalf("expected 0 skipped, got %d", len(res.Skipped))
	}
}

func TestPin_SpecificKeys(t *testing.T) {
	res, err := Pin(pinEntries(), []string{"PORT", "APP_HOST"}, pinTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 2 {
		t.Fatalf("expected 2 pinned, got %d", len(res.Pinned))
	}
}

func TestPin_MissingKeySkipped(t *testing.T) {
	res, err := Pin(pinEntries(), []string{"MISSING_KEY"}, pinTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 0 {
		t.Fatalf("expected 0 pinned, got %d", len(res.Pinned))
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING_KEY" {
		t.Fatalf("expected MISSING_KEY in skipped, got %v", res.Skipped)
	}
}

func TestPin_SecretValueMaskedInComment(t *testing.T) {
	res, err := Pin(pinEntries(), []string{"DB_PASSWORD"}, pinTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Pinned) != 1 {
		t.Fatalf("expected 1 pinned entry")
	}
	p := res.Pinned[0]
	if p.Value != "s3cr3t" {
		t.Errorf("expected raw value preserved, got %q", p.Value)
	}
	if !contains(p.Comment, "***") {
		t.Errorf("expected masked comment, got %q", p.Comment)
	}
}

func TestPin_PinnedToEntries(t *testing.T) {
	res, _ := Pin(pinEntries(), nil, pinTime)
	entries := PinnedToEntries(res.Pinned)
	if len(entries) != len(res.Pinned) {
		t.Fatalf("length mismatch: %d vs %d", len(entries), len(res.Pinned))
	}
	for i, e := range entries {
		if e.Key != res.Pinned[i].Key || e.Value != res.Pinned[i].Value {
			t.Errorf("entry %d mismatch: got %+v", i, e)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
