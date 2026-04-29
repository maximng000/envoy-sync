package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func archiveEntries(pairs ...string) map[string]Entry {
	m := make(map[string]Entry)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = Entry{Key: pairs[i], Value: pairs[i+1]}
	}
	return m
}

func TestAddToArchive_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "archive.json")

	entries := archiveEntries("APP_ENV", "production")
	if err := AddToArchive(path, entries, "v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("archive file not created: %v", err)
	}
}

func TestAddToArchive_AccumulatesVersions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "archive.json")

	if err := AddToArchive(path, archiveEntries("KEY", "v1"), "first"); err != nil {
		t.Fatal(err)
	}
	if err := AddToArchive(path, archiveEntries("KEY", "v2"), "second"); err != nil {
		t.Fatal(err)
	}

	arch, err := LoadArchive(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(arch.Versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(arch.Versions))
	}
	if arch.Versions[0].Label != "first" || arch.Versions[1].Label != "second" {
		t.Error("labels not preserved correctly")
	}
}

func TestLoadArchive_MissingFile(t *testing.T) {
	_, err := LoadArchive("/nonexistent/archive.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLatestArchiveEntry_ReturnsLast(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "archive.json")

	_ = AddToArchive(path, archiveEntries("A", "1"), "old")
	_ = AddToArchive(path, archiveEntries("A", "2"), "new")

	arch, _ := LoadArchive(path)
	latest, err := LatestArchiveEntry(arch)
	if err != nil {
		t.Fatal(err)
	}
	if latest.Label != "new" {
		t.Errorf("expected label 'new', got %q", latest.Label)
	}
}

func TestDiffArchiveVersions_DetectsChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "archive.json")

	_ = AddToArchive(path, archiveEntries("APP_ENV", "staging"), "")
	_ = AddToArchive(path, archiveEntries("APP_ENV", "production", "NEW_KEY", "hello"), "")

	arch, _ := LoadArchive(path)
	diffs, err := DiffArchiveVersions(arch, 0, 1)
	if err != nil {
		t.Fatal(err)
	}

	changed, added := 0, 0
	for _, d := range diffs {
		switch d.Status {
		case StatusChanged:
			changed++
		case StatusAdded:
			added++
		}
	}
	if changed != 1 {
		t.Errorf("expected 1 changed, got %d", changed)
	}
	if added != 1 {
		t.Errorf("expected 1 added, got %d", added)
	}
}

func TestDiffArchiveVersions_OutOfRange(t *testing.T) {
	arch := &Archive{Versions: []ArchiveEntry{{}, {}}}
	if _, err := DiffArchiveVersions(arch, 0, 5); err == nil {
		t.Error("expected error for out-of-range index")
	}
}
