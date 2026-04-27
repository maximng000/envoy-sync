package envfile

import (
	"strings"
	"testing"
)

func tagEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "API_KEY", Value: "key123"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestTag_AllEntries(t *testing.T) {
	entries := tagEntries()
	result, err := Tag(entries, nil, "v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tagged) != 4 {
		t.Errorf("expected 4 tagged, got %d", len(result.Tagged))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
	for _, te := range result.Tagged {
		if te.Tag != "v1" {
			t.Errorf("expected tag 'v1', got '%s'", te.Tag)
		}
	}
}

func TestTag_SpecificKeys(t *testing.T) {
	entries := tagEntries()
	result, err := Tag(entries, []string{"APP_NAME", "PORT"}, "release")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tagged) != 2 {
		t.Errorf("expected 2 tagged, got %d", len(result.Tagged))
	}
	if len(result.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(result.Skipped))
	}
}

func TestTag_EmptyTagReturnsError(t *testing.T) {
	_, err := Tag(tagEntries(), nil, "")
	if err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestGroupByTag_MultipleGroups(t *testing.T) {
	tagged := []TagEntry{
		{Key: "A", Value: "1", Tag: "prod"},
		{Key: "B", Value: "2", Tag: "dev"},
		{Key: "C", Value: "3", Tag: "prod"},
	}
	groups := GroupByTag(tagged)
	if len(groups["prod"]) != 2 {
		t.Errorf("expected 2 prod entries, got %d", len(groups["prod"]))
	}
	if len(groups["dev"]) != 1 {
		t.Errorf("expected 1 dev entry, got %d", len(groups["dev"]))
	}
}

func TestTagSummary_ContainsTagHeaders(t *testing.T) {
	tagged := []TagEntry{
		{Key: "APP_NAME", Value: "myapp", Tag: "stable"},
		{Key: "PORT", Value: "8080", Tag: "stable"},
	}
	summary := TagSummary(tagged)
	if !strings.Contains(summary, "[stable]") {
		t.Error("expected summary to contain '[stable]'")
	}
	if !strings.Contains(summary, "APP_NAME=myapp") {
		t.Error("expected summary to contain 'APP_NAME=myapp'")
	}
}
