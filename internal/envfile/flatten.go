package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenResult holds the result of flattening a nested prefix map.
type FlattenResult struct {
	Key    string
	Value  string
	Prefix string
}

// Flatten groups entries by a delimiter-separated prefix and returns a
// structured list. For example, DB_HOST and DB_PORT share the prefix "DB".
// If prefix is non-empty, only entries matching that prefix are returned.
func Flatten(entries map[string]string, delimiter string, prefix string) []FlattenResult {
	if delimiter == "" {
		delimiter = "_"
	}

	var results []FlattenResult

	for _, k := range sortedFlattenKeys(entries) {
		v := entries[k]
		parts := strings.SplitN(k, delimiter, 2)
		pfx := ""
		if len(parts) == 2 {
			pfx = parts[0]
		}

		if prefix != "" && !strings.EqualFold(pfx, prefix) {
			continue
		}

		results = append(results, FlattenResult{
			Key:    k,
			Value:  v,
			Prefix: pfx,
		})
	}

	return results
}

// FlattenSummary returns a map of prefix -> list of keys under that prefix.
func FlattenSummary(entries map[string]string, delimiter string) map[string][]string {
	if delimiter == "" {
		delimiter = "_"
	}

	summary := make(map[string][]string)

	for _, k := range sortedFlattenKeys(entries) {
		parts := strings.SplitN(k, delimiter, 2)
		pfx := fmt.Sprintf("(no prefix)")
		if len(parts) == 2 {
			pfx = parts[0]
		}
		summary[pfx] = append(summary[pfx], k)
	}

	return summary
}

func sortedFlattenKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
