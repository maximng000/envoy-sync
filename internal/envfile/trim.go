package envfile

import (
	"strings"
)

// TrimResult holds the outcome of a trim operation on a single entry.
type TrimResult struct {
	Key      string
	OldValue string
	NewValue string
	Changed  bool
}

// Trim removes leading and trailing whitespace from all entry values.
// It returns the updated entries map and a slice of TrimResults describing
// which entries were changed.
func Trim(entries map[string]string) (map[string]string, []TrimResult) {
	result := make(map[string]string, len(entries))
	var changes []TrimResult

	for k, v := range entries {
		trimmed := strings.TrimSpace(v)
		result[k] = trimmed
		if trimmed != v {
			changes = append(changes, TrimResult{
				Key:      k,
				OldValue: v,
				NewValue: trimmed,
				Changed:  true,
			})
		}
	}

	return result, changes
}

// TrimKeys removes leading and trailing whitespace from a specific set of keys.
// Keys not present in entries are ignored.
func TrimKeys(entries map[string]string, keys []string) (map[string]string, []TrimResult) {
	result := copyEntries(entries)
	var changes []TrimResult

	for _, k := range keys {
		v, ok := result[k]
		if !ok {
			continue
		}
		trimmed := strings.TrimSpace(v)
		result[k] = trimmed
		if trimmed != v {
			changes = append(changes, TrimResult{
				Key:      k,
				OldValue: v,
				NewValue: trimmed,
				Changed:  true,
			})
		}
	}

	return result, changes
}

func copyEntries(entries map[string]string) map[string]string {
	out := make(map[string]string, len(entries))
	for k, v := range entries {
		out[k] = v
	}
	return out
}
