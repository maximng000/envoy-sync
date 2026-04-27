package envfile

import (
	"testing"
)

func normalizeEntries() []Entry {
	return []Entry{
		{Key: "app_name", Value: "  myapp  "},
		{Key: "db_host", Value: "localhost"},
		{Key: "empty_key", Value: ""},
		{Key: "DB_HOST", Value: "prod-host"},
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	entries := []Entry{
		{Key: "app_name", Value: "myapp"},
		{Key: "db_host", Value: "localhost"},
	}
	r := Normalize(entries, NormalizeOptions{UppercaseKeys: true})
	if r.Entries[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", r.Entries[0].Key)
	}
	if r.Entries[1].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", r.Entries[1].Key)
	}
	if len(r.Modified) != 2 {
		t.Errorf("expected 2 modified, got %d", len(r.Modified))
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	entries := []Entry{
		{Key: "APP", Value: "  hello  "},
	}
	r := Normalize(entries, NormalizeOptions{TrimValues: true})
	if r.Entries[0].Value != "hello" {
		t.Errorf("expected 'hello', got %q", r.Entries[0].Value)
	}
	if len(r.Modified) != 1 {
		t.Errorf("expected 1 modified")
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	entries := []Entry{
		{Key: "PRESENT", Value: "yes"},
		{Key: "EMPTY", Value: ""},
	}
	r := Normalize(entries, NormalizeOptions{RemoveEmpty: true})
	if len(r.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(r.Entries))
	}
	if len(r.Removed) != 1 || r.Removed[0] != "EMPTY" {
		t.Errorf("expected EMPTY in removed")
	}
}

func TestNormalize_DeduplicatesAfterUppercase(t *testing.T) {
	entries := []Entry{
		{Key: "db_host", Value: "first"},
		{Key: "DB_HOST", Value: "second"},
	}
	r := Normalize(entries, NormalizeOptions{UppercaseKeys: true})
	if len(r.Entries) != 1 {
		t.Errorf("expected 1 entry after dedup, got %d", len(r.Entries))
	}
	if r.Entries[0].Value != "first" {
		t.Errorf("expected first value to be kept")
	}
}

func TestNormalize_NoChanges(t *testing.T) {
	entries := []Entry{
		{Key: "KEY", Value: "value"},
	}
	r := Normalize(entries, NormalizeOptions{})
	if len(r.Modified) != 0 {
		t.Errorf("expected no modifications")
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected no removals")
	}
}

func TestNormalizeSummary_ContainsFields(t *testing.T) {
	r := NormalizeResult{
		Entries:  []Entry{{Key: "A", Value: "1"}},
		Modified: []string{"B"},
		Removed:  []string{"C", "D"},
	}
	s := NormalizeSummary(r)
	for _, want := range []string{"kept", "modified", "removed"} {
		if !containsStr(s, want) {
			t.Errorf("summary missing %q", want)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s[1:], sub) || s[:len(sub)] == sub)
}
