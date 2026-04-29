package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ArchiveEntry represents a single versioned snapshot of an env file.
type ArchiveEntry struct {
	Timestamp time.Time        `json:"timestamp"`
	Label     string           `json:"label,omitempty"`
	Entries   map[string]Entry `json:"entries"`
}

// Archive holds a collection of versioned env snapshots.
type Archive struct {
	Versions []ArchiveEntry `json:"versions"`
}

// AddToArchive appends the current entries as a new version in the archive file.
// If the archive file does not exist, it is created.
func AddToArchive(archivePath string, entries map[string]Entry, label string) error {
	arch, err := LoadArchive(archivePath)
	if err != nil {
		arch = &Archive{}
	}

	snap := make(map[string]Entry, len(entries))
	for k, v := range entries {
		snap[k] = v
	}

	arch.Versions = append(arch.Versions, ArchiveEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   snap,
	})

	if err := os.MkdirAll(filepath.Dir(archivePath), 0o755); err != nil {
		return fmt.Errorf("archive: create dir: %w", err)
	}

	data, err := json.MarshalIndent(arch, "", "  ")
	if err != nil {
		return fmt.Errorf("archive: marshal: %w", err)
	}
	return os.WriteFile(archivePath, data, 0o644)
}

// LoadArchive reads and parses an archive file from disk.
func LoadArchive(archivePath string) (*Archive, error) {
	data, err := os.ReadFile(archivePath)
	if err != nil {
		return nil, fmt.Errorf("archive: read: %w", err)
	}
	var arch Archive
	if err := json.Unmarshal(data, &arch); err != nil {
		return nil, fmt.Errorf("archive: unmarshal: %w", err)
	}
	return &arch, nil
}

// LatestArchiveEntry returns the most recently added entry, or an error if empty.
func LatestArchiveEntry(arch *Archive) (*ArchiveEntry, error) {
	if len(arch.Versions) == 0 {
		return nil, fmt.Errorf("archive: no versions stored")
	}
	return &arch.Versions[len(arch.Versions)-1], nil
}

// DiffArchiveVersions computes a Diff between two archive entries by index.
func DiffArchiveVersions(arch *Archive, fromIdx, toIdx int) ([]DiffEntry, error) {
	if fromIdx < 0 || fromIdx >= len(arch.Versions) {
		return nil, fmt.Errorf("archive: from index %d out of range", fromIdx)
	}
	if toIdx < 0 || toIdx >= len(arch.Versions) {
		return nil, fmt.Errorf("archive: to index %d out of range", toIdx)
	}
	return Diff(arch.Versions[fromIdx].Entries, arch.Versions[toIdx].Entries), nil
}
