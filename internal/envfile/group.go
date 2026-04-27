package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// GroupResult holds entries organized by a grouping key.
type GroupResult struct {
	Key     string
	Entries []Entry
}

// GroupBy groups entries by a strategy: "prefix" (by KEY prefix before delimiter),
// "secret" (secrets vs non-secrets), or "empty" (empty vs non-empty values).
func GroupBy(entries []Entry, strategy string, delimiter string) ([]GroupResult, error) {
	if delimiter == "" {
		delimiter = "_"
	}

	switch strategy {
	case "prefix":
		return groupByPrefix(entries, delimiter), nil
	case "secret":
		return groupBySecret(entries), nil
	case "empty":
		return groupByEmpty(entries), nil
	default:
		return nil, fmt.Errorf("unknown grouping strategy %q: must be one of prefix, secret, empty", strategy)
	}
}

func groupByPrefix(entries []Entry, delimiter string) []GroupResult {
	buckets := map[string][]Entry{}

	for _, e := range entries {
		parts := strings.SplitN(e.Key, delimiter, 2)
		prefix := parts[0]
		if len(parts) == 1 {
			prefix = "(no prefix)"
		}
		buckets[prefix] = append(buckets[prefix], e)
	}

	return sortedGroups(buckets)
}

func groupBySecret(entries []Entry) []GroupResult {
	buckets := map[string][]Entry{
		"secrets":     {},
		"non-secrets": {},
	}

	for _, e := range entries {
		if IsSecret(e.Key) {
			buckets["secrets"] = append(buckets["secrets"], e)
		} else {
			buckets["non-secrets"] = append(buckets["non-secrets"], e)
		}
	}

	return sortedGroups(buckets)
}

func groupByEmpty(entries []Entry) []GroupResult {
	buckets := map[string][]Entry{
		"empty":     {},
		"non-empty": {},
	}

	for _, e := range entries {
		if strings.TrimSpace(e.Value) == "" {
			buckets["empty"] = append(buckets["empty"], e)
		} else {
			buckets["non-empty"] = append(buckets["non-empty"], e)
		}
	}

	return sortedGroups(buckets)
}

func sortedGroups(buckets map[string][]Entry) []GroupResult {
	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]GroupResult, 0, len(keys))
	for _, k := range keys {
		results = append(results, GroupResult{Key: k, Entries: buckets[k]})
	}
	return results
}
