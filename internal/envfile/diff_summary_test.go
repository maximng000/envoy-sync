package envfile

import (
	"testing"
)

func diffEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "A", Status: "added", NewValue: "1"},
		{Key: "B", Status: "removed", OldValue: "2"},
		{Key: "C", Status: "changed", OldValue: "x", NewValue: "y"},
		{Key: "D", Status: "unchanged", OldValue: "z", NewValue: "z"},
		{Key: "E", Status: "unchanged", OldValue: "q", NewValue: "q"},
	}
}

func TestSummarizeDiff_Empty(t *testing.T) {
	s := SummarizeDiff(nil)
	if s.Total != 0 || s.HasChanges {
		t.Errorf("expected empty summary, got %+v", s)
	}
}

func TestSummarizeDiff_CountsCorrectly(t *testing.T) {
	s := SummarizeDiff(diffEntries())
	if s.Added != 1 {
		t.Errorf("Added: want 1, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("Changed: want 1, got %d", s.Changed)
	}
	if s.Unchanged != 2 {
		t.Errorf("Unchanged: want 2, got %d", s.Unchanged)
	}
	if s.Total != 5 {
		t.Errorf("Total: want 5, got %d", s.Total)
	}
}

func TestSummarizeDiff_HasChanges(t *testing.T) {
	s := SummarizeDiff(diffEntries())
	if !s.HasChanges {
		t.Error("expected HasChanges to be true")
	}
}

func TestSummarizeDiff_NoChangesWhenUnchanged(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: "unchanged", OldValue: "1", NewValue: "1"},
		{Key: "B", Status: "unchanged", OldValue: "2", NewValue: "2"},
	}
	s := SummarizeDiff(entries)
	if s.HasChanges {
		t.Error("expected HasChanges to be false")
	}
	if s.Unchanged != 2 {
		t.Errorf("Unchanged: want 2, got %d", s.Unchanged)
	}
}

func TestSummarizeDiff_OnlyAdded(t *testing.T) {
	entries := []DiffEntry{
		{Key: "X", Status: "added", NewValue: "v1"},
		{Key: "Y", Status: "added", NewValue: "v2"},
	}
	s := SummarizeDiff(entries)
	if s.Added != 2 {
		t.Errorf("Added: want 2, got %d", s.Added)
	}
	if !s.HasChanges {
		t.Error("expected HasChanges true")
	}
	if s.Total != 2 {
		t.Errorf("Total: want 2, got %d", s.Total)
	}
}
