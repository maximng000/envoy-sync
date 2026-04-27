package envfile

import (
	"testing"
)

func diffEntries(specs ...struct {
	key    string
	status DiffStatus
}) []DiffSummaryEntry {
	out := make([]DiffSummaryEntry, 0, len(specs))
	for _, s := range specs {
		out = append(out, DiffSummaryEntry{Key: s.key, Status: s.status})
	}
	return out
}

func TestSummarizeDiff_Empty(t *testing.T) {
	r := SummarizeDiff(nil)
	if r.Total != 0 || r.HasChanges {
		t.Errorf("expected empty summary, got %+v", r)
	}
}

func TestSummarizeDiff_CountsCorrectly(t *testing.T) {
	entries := []DiffSummaryEntry{
		{Key: "A", Status: DiffAdded},
		{Key: "B", Status: DiffRemoved},
		{Key: "C", Status: DiffChanged},
		{Key: "D", Status: DiffUnchanged},
		{Key: "E", Status: DiffAdded},
	}
	r := SummarizeDiff(entries)
	if r.Added != 2 {
		t.Errorf("Added: want 2, got %d", r.Added)
	}
	if r.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", r.Removed)
	}
	if r.Changed != 1 {
		t.Errorf("Changed: want 1, got %d", r.Changed)
	}
	if r.Unchanged != 1 {
		t.Errorf("Unchanged: want 1, got %d", r.Unchanged)
	}
	if r.Total != 5 {
		t.Errorf("Total: want 5, got %d", r.Total)
	}
}

func TestSummarizeDiff_HasChanges(t *testing.T) {
	entries := []DiffSummaryEntry{
		{Key: "X", Status: DiffChanged},
	}
	r := SummarizeDiff(entries)
	if !r.HasChanges {
		t.Error("expected HasChanges to be true")
	}
}

func TestSummarizeDiff_NoChangesWhenUnchanged(t *testing.T) {
	entries := []DiffSummaryEntry{
		{Key: "A", Status: DiffUnchanged},
		{Key: "B", Status: DiffUnchanged},
	}
	r := SummarizeDiff(entries)
	if r.HasChanges {
		t.Error("expected HasChanges to be false")
	}
	if r.Total != 2 {
		t.Errorf("Total: want 2, got %d", r.Total)
	}
}

func TestSummarizeDiff_OnlyRemoved(t *testing.T) {
	entries := []DiffSummaryEntry{
		{Key: "GONE", Status: DiffRemoved},
	}
	r := SummarizeDiff(entries)
	if r.Removed != 1 || !r.HasChanges {
		t.Errorf("unexpected result: %+v", r)
	}
}
