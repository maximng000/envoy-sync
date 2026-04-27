package envfile

import (
	"testing"
)

func scopeEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "AWS_ACCESS_KEY", Value: "AKIA123"},
		{Key: "AWS_SECRET", Value: "abc"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "NOPREFIX", Value: "val"},
	}
}

func TestScope_FiltersByPrefix(t *testing.T) {
	result := Scope(scopeEntries(), "DB")
	if result.Scope != "DB" {
		t.Fatalf("expected scope DB, got %s", result.Scope)
	}
	if len(result.Entries) != 3 {
		t.Fatalf("expected 3 DB entries, got %d", len(result.Entries))
	}
}

func TestScope_CaseInsensitiveInput(t *testing.T) {
	result := Scope(scopeEntries(), "aws")
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 AWS entries, got %d", len(result.Entries))
	}
}

func TestScope_NoMatch(t *testing.T) {
	result := Scope(scopeEntries(), "REDIS")
	if len(result.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result.Entries))
	}
}

func TestListScopes_ReturnsDistinctPrefixes(t *testing.T) {
	scopes := ListScopes(scopeEntries())
	expected := []string{"APP", "AWS", "DB"}
	if len(scopes) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, scopes)
	}
	for i, s := range expected {
		if scopes[i] != s {
			t.Errorf("expected %s at index %d, got %s", s, i, scopes[i])
		}
	}
}

func TestListScopes_NoPrefixKeyExcluded(t *testing.T) {
	scopes := ListScopes(scopeEntries())
	for _, s := range scopes {
		if s == "NOPREFIX" {
			t.Error("NOPREFIX should not appear as a scope")
		}
	}
}

func TestScopeSummaryOf_CountsCorrect(t *testing.T) {
	summary := ScopeSummaryOf(scopeEntries())
	if summary.Counts["DB"] != 3 {
		t.Errorf("expected DB count 3, got %d", summary.Counts["DB"])
	}
	if summary.Counts["AWS"] != 2 {
		t.Errorf("expected AWS count 2, got %d", summary.Counts["AWS"])
	}
	if summary.Counts["APP"] != 1 {
		t.Errorf("expected APP count 1, got %d", summary.Counts["APP"])
	}
}

func TestFormatScopeSummary_Empty(t *testing.T) {
	summary := ScopeSummaryOf([]Entry{})
	out := FormatScopeSummary(summary)
	if out != "no scopes detected" {
		t.Errorf("unexpected output: %s", out)
	}
}
