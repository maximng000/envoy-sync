package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTakeSnapshot_CopiesEntries(t *testing.T) {
	original := map[string]string{"FOO": "bar", "BAZ": "qux"}
	snap := TakeSnapshot("test.env", original)

	if snap.Source != "test.env" {
		t.Errorf("expected source test.env, got %s", snap.Source)
	}
	if snap.Entries["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
	// Mutating original should not affect snapshot
	original["FOO"] = "changed"
	if snap.Entries["FOO"] != "bar" {
		t.Errorf("snapshot should be independent of original map")
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	original := Snapshot{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Source:    "prod.env",
		Entries:   map[string]string{"KEY": "value", "SECRET_TOKEN": "abc123"},
	}

	if err := SaveSnapshot(path, original); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if loaded.Source != original.Source {
		t.Errorf("source mismatch: got %s", loaded.Source)
	}
	if loaded.Entries["KEY"] != "value" {
		t.Errorf("expected KEY=value")
	}
	if !loaded.Timestamp.Equal(original.Timestamp) {
		t.Errorf("timestamp mismatch")
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveSnapshot_InvalidPath(t *testing.T) {
	snap := TakeSnapshot("x.env", map[string]string{"A": "1"})
	err := SaveSnapshot("/nonexistent/dir/snap.json", snap)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	snap := TakeSnapshot("base.env", map[string]string{
		"FOO": "old",
		"BAR": "same",
	})

	current := map[string]string{
		"FOO": "new",
		"BAR": "same",
		"NEW": "added",
	}

	diffs := DiffSnapshot(snap, current)
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}

	_ = os.Getenv // suppress unused import
}
