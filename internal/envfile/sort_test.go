package envfile

import (
	"testing"
)

func sortEntries() []Entry {
	return []Entry{
		{Key: "ZEBRA", Value: "z"},
		{Key: "API_KEY", Value: "secret"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "pass"},
		{Key: "HOST", Value: "localhost"},
	}
}

func TestSort_Alpha(t *testing.T) {
	res, err := Sort(sortEntries(), SortAlpha)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY first, got %s", res.Entries[0].Key)
	}
	if res.Entries[len(res.Entries)-1].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA last, got %s", res.Entries[len(res.Entries)-1].Key)
	}
}

func TestSort_AlphaDesc(t *testing.T) {
	res, err := Sort(sortEntries(), SortAlphaDesc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA first, got %s", res.Entries[0].Key)
	}
}

func TestSort_BySecret(t *testing.T) {
	res, err := Sort(sortEntries(), SortBySecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Secrets should come first
	if !IsSecret(res.Entries[0].Key) {
		t.Errorf("expected first entry to be a secret, got %s", res.Entries[0].Key)
	}
}

func TestSort_ByLength(t *testing.T) {
	res, err := Sort(sortEntries(), SortByLength)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Shortest key first
	if len(res.Entries[0].Key) > len(res.Entries[1].Key) {
		t.Errorf("expected shorter key first")
	}
}

func TestSort_UnknownStrategy(t *testing.T) {
	_, err := Sort(sortEntries(), SortStrategy("bogus"))
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestSort_OriginalUnmodified(t *testing.T) {
	orig := sortEntries()
	firstKey := orig[0].Key
	_, _ = Sort(orig, SortAlpha)
	if orig[0].Key != firstKey {
		t.Error("original slice was modified")
	}
}

func TestSort_TotalCount(t *testing.T) {
	entries := sortEntries()
	res, err := Sort(entries, SortAlpha)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != len(entries) {
		t.Errorf("expected Total=%d, got %d", len(entries), res.Total)
	}
}
