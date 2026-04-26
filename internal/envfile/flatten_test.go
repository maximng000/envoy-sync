package envfile

import (
	"testing"
)

func TestFlatten_GroupsByPrefix(t *testing.T) {
	entries := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}

	results := Flatten(entries, "_", "")
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestFlatten_FilterByPrefix(t *testing.T) {
	entries := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}

	results := Flatten(entries, "_", "DB")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Prefix != "DB" {
			t.Errorf("expected prefix DB, got %s", r.Prefix)
		}
	}
}

func TestFlatten_NoPrefixKey(t *testing.T) {
	entries := map[string]string{
		"SIMPLE": "value",
	}

	results := Flatten(entries, "_", "")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Prefix != "" {
		t.Errorf("expected empty prefix, got %q", results[0].Prefix)
	}
}

func TestFlatten_DefaultDelimiter(t *testing.T) {
	entries := map[string]string{
		"DB_HOST": "localhost",
	}

	results := Flatten(entries, "", "")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Prefix != "DB" {
		t.Errorf("expected prefix DB, got %q", results[0].Prefix)
	}
}

func TestFlattenSummary_GroupsCorrectly(t *testing.T) {
	entries := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
		"APP_NAME": "myapp",
		"SIMPLE": "val",
	}

	summary := FlattenSummary(entries, "_")

	if len(summary["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(summary["DB"]))
	}
	if len(summary["APP"]) != 2 {
		t.Errorf("expected 2 APP keys, got %d", len(summary["APP"]))
	}
	if len(summary["(no prefix)"]) != 1 {
		t.Errorf("expected 1 no-prefix key, got %d", len(summary["(no prefix)"]))
	}
}

func TestFlattenSummary_EmptyEntries(t *testing.T) {
	summary := FlattenSummary(map[string]string{}, "_")
	if len(summary) != 0 {
		t.Errorf("expected empty summary, got %v", summary)
	}
}
