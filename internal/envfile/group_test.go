package envfile

import (
	"testing"
)

func groupEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "APP_SECRET", Value: "s3cr3t"},
		{Key: "STANDALONE", Value: ""},
		{Key: "API_KEY", Value: "abc123"},
	}
}

func TestGroupBy_Prefix(t *testing.T) {
	groups, err := GroupBy(groupEntries(), "prefix", "_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	groupMap := map[string][]Entry{}
	for _, g := range groups {
		groupMap[g.Key] = g.Entries
	}

	if len(groupMap["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(groupMap["DB"]))
	}
	if len(groupMap["APP"]) != 2 {
		t.Errorf("expected 2 APP entries, got %d", len(groupMap["APP"]))
	}
	if len(groupMap["API"]) != 1 {
		t.Errorf("expected 1 API entry, got %d", len(groupMap["API"]))
	}
	if len(groupMap["(no prefix)"]) != 1 {
		t.Errorf("expected 1 standalone entry, got %d", len(groupMap["(no prefix)"]))
	}
}

func TestGroupBy_Secret(t *testing.T) {
	groups, err := GroupBy(groupEntries(), "secret", "_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	groupMap := map[string][]Entry{}
	for _, g := range groups {
		groupMap[g.Key] = g.Entries
	}

	if len(groupMap["secrets"]) == 0 {
		t.Error("expected at least one secret entry")
	}
	if len(groupMap["non-secrets"]) == 0 {
		t.Error("expected at least one non-secret entry")
	}
	for _, e := range groupMap["secrets"] {
		if !IsSecret(e.Key) {
			t.Errorf("entry %q should be a secret", e.Key)
		}
	}
}

func TestGroupBy_Empty(t *testing.T) {
	groups, err := GroupBy(groupEntries(), "empty", "_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	groupMap := map[string][]Entry{}
	for _, g := range groups {
		groupMap[g.Key] = g.Entries
	}

	if len(groupMap["empty"]) != 1 {
		t.Errorf("expected 1 empty entry, got %d", len(groupMap["empty"]))
	}
	if len(groupMap["non-empty"]) != 5 {
		t.Errorf("expected 5 non-empty entries, got %d", len(groupMap["non-empty"]))
	}
}

func TestGroupBy_UnknownStrategy(t *testing.T) {
	_, err := GroupBy(groupEntries(), "unknown", "_")
	if err == nil {
		t.Error("expected error for unknown strategy, got nil")
	}
}

func TestGroupBy_DefaultDelimiter(t *testing.T) {
	groups, err := GroupBy(groupEntries(), "prefix", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) == 0 {
		t.Error("expected groups with default delimiter")
	}
}
