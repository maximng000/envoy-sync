package envfile

import (
	"fmt"
	"strings"
)

// RotateResult describes the outcome of rotating a single secret entry.
type RotateResult struct {
	Key      string
	OldValue string
	NewValue string
	Rotated  bool
	Skipped  bool
	Reason   string
}

// RotateFunc is a user-supplied function that generates a new value for a given key.
type RotateFunc func(key, oldValue string) (string, error)

// Rotate replaces the values of secret keys in entries using the provided RotateFunc.
// Non-secret keys are skipped unless their key appears in forceKeys.
// Returns updated entries and a slice of RotateResult for reporting.
func Rotate(entries []Entry, fn RotateFunc, forceKeys []string) ([]Entry, []RotateResult, error) {
	forced := make(map[string]bool, len(forceKeys))
	for _, k := range forceKeys {
		forced[strings.ToUpper(k)] = true
	}

	updated := make([]Entry, 0, len(entries))
	results := make([]RotateResult, 0, len(entries))

	for _, e := range entries {
		result := RotateResult{Key: e.Key, OldValue: e.Value}

		if !IsSecret(e.Key) && !forced[strings.ToUpper(e.Key)] {
			result.Skipped = true
			result.Reason = "not a secret"
			updated = append(updated, e)
			results = append(results, result)
			continue
		}

		newVal, err := fn(e.Key, e.Value)
		if err != nil {
			return nil, nil, fmt.Errorf("rotate: key %q: %w", e.Key, err)
		}

		result.NewValue = newVal
		result.Rotated = true
		e.Value = newVal
		updated = append(updated, e)
		results = append(results, result)
	}

	return updated, results, nil
}
