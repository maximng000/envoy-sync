package envfile

import (
	"testing"
)

func diffEntries(statuses ...string) []DiffEntry {
	out := make([]DiffEntry, 0, len(statuses))
	for i, s := range statuses {
		out = append(out, DiffEntry{
			Key:    fmt.Sprintf("KEY_%d", i),
			Status: s,
		})
	}
	return out
}

func TestSummarizeDiff_Empty(t *testing.T) {
	s := SummarizeDiff(nil)
	if s.Total != 0 || s.HasChanges() {
		t.Errorf("expected empty summary, got %+v", s)
	}
}

func TestSummarizeDiff_CountsCorrectly(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: "added"},
		{Key: "B", Status: "added"},
		{Key: "C", Status: "removed"},
		{Key: "D", Status: "changed"},
	}
	s := SummarizeDiff(entries)
	if s.Added != 2 {
		t.Errorf("Added: want 2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("Changed: want 1, got %d", s.Changed)
	}
	if s.Total != 4 {
		t.Errorf("Total: want 4, got %d", s.Total)
	}
}

func TestSummarizeDiff_HasChanges(t *testing.T) {
	entries := []DiffEntry{{Key: "X", Status: "added"}}
	if !SummarizeDiff(entries).HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestSummarizeDiff_NoChangesWhenUnchanged(t *testing.T) {
	entries := []DiffEntry{{Key: "X", Status: "unchanged"}}
	if SummarizeDiff(entries).HasChanges() {
		t.Error("expected HasChanges to be false for unchanged entries")
	}
}

func TestDiffSummary_String(t *testing.T) {
	s := DiffSummary{Added: 1, Removed: 2, Changed: 3, Total: 6}
	want := "total=6 added=1 removed=2 changed=3"
	if got := s.String(); got != want {
		t.Errorf("String(): want %q, got %q", want, got)
	}
}
