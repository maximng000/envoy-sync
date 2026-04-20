package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an env file.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Entries   map[string]string `json:"entries"`
}

// TakeSnapshot creates a Snapshot from a parsed env map.
func TakeSnapshot(source string, entries map[string]string) Snapshot {
	copy := make(map[string]string, len(entries))
	for k, v := range entries {
		copy[k] = v
	}
	return Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Entries:   copy,
	}
}

// SaveSnapshot writes a Snapshot to a JSON file at the given path.
func SaveSnapshot(path string, snap Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a Snapshot from a JSON file at the given path.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return snap, nil
}

// DiffSnapshot compares a current env map against a saved Snapshot and
// returns the differences using the existing Diff logic.
func DiffSnapshot(snap Snapshot, current map[string]string) []DiffEntry {
	return Diff(snap.Entries, current)
}
