package envfile

import (
	"fmt"
	"sort"
	"time"
)

// PinnedEntry represents a key whose value has been locked at a specific time.
type PinnedEntry struct {
	Key       string
	Value     string
	PinnedAt  time.Time
	Comment   string
}

// PinResult holds the outcome of a Pin operation.
type PinResult struct {
	Pinned  []PinnedEntry
	Skipped []string
}

// Pin locks the values of the specified keys (or all keys if keys is empty),
// returning a PinResult with the frozen entries. Secret values are masked in
// the Comment field but stored as-is in Value for downstream use.
func Pin(entries []Entry, keys []string, now time.Time) (PinResult, error) {
	if now.IsZero() {
		now = time.Now().UTC()
	}

	index := make(map[string]string, len(entries))
	for _, e := range entries {
		index[e.Key] = e.Value
	}

	wantAll := len(keys) == 0
	if wantAll {
		for _, e := range entries {
			keys = append(keys, e.Key)
		}
	}

	sort.Strings(keys)

	var result PinResult
	for _, k := range keys {
		v, ok := index[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		display := v
		if IsSecret(k) {
			display = "***"
		}
		result.Pinned = append(result.Pinned, PinnedEntry{
			Key:      k,
			Value:    v,
			PinnedAt: now,
			Comment:  fmt.Sprintf("pinned at %s (value: %s)", now.Format(time.RFC3339), display),
		})
	}
	return result, nil
}

// PinnedToEntries converts a slice of PinnedEntry back to []Entry.
func PinnedToEntries(pinned []PinnedEntry) []Entry {
	out := make([]Entry, len(pinned))
	for i, p := range pinned {
		out[i] = Entry{Key: p.Key, Value: p.Value}
	}
	return out
}
