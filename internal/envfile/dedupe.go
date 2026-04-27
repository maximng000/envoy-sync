package envfile

import "fmt"

// DedupeStrategy controls how duplicate keys are resolved.
type DedupeStrategy string

const (
	DedupeKeepFirst DedupeStrategy = "first"
	DedupeKeepLast  DedupeStrategy = "last"
)

// DuplicateEntry records a key that appeared more than once.
type DuplicateEntry struct {
	Key    string
	Values []string
	Kept   string
}

// DedupeResult holds the deduplicated entries and a report of what was resolved.
type DedupeResult struct {
	Entries    []Entry
	Duplicates []DuplicateEntry
}

// Dedupe removes duplicate keys from a slice of entries.
// When a key appears multiple times, strategy controls which value is kept.
// The original order of first occurrence is preserved.
func Dedupe(entries []Entry, strategy DedupeStrategy) (DedupeResult, error) {
	if strategy != DedupeKeepFirst && strategy != DedupeKeepLast {
		return DedupeResult{}, fmt.Errorf("unknown dedupe strategy: %q", strategy)
	}

	// Track all values seen per key in order.
	type occurrence struct {
		index int
		entry Entry
	}
	seen := make(map[string][]occurrence)
	order := []string{}

	for i, e := range entries {
		if _, exists := seen[e.Key]; !exists {
			order = append(order, e.Key)
		}
		seen[e.Key] = append(seen[e.Key], occurrence{index: i, entry: e})
	}

	result := DedupeResult{}

	for _, key := range order {
		occurrences := seen[key]
		if len(occurrences) == 1 {
			result.Entries = append(result.Entries, occurrences[0].entry)
			continue
		}

		var kept Entry
		if strategy == DedupeKeepFirst {
			kept = occurrences[0].entry
		} else {
			kept = occurrences[len(occurrences)-1].entry
		}

		values := make([]string, len(occurrences))
		for i, o := range occurrences {
			values[i] = o.entry.Value
		}

		result.Entries = append(result.Entries, kept)
		result.Duplicates = append(result.Duplicates, DuplicateEntry{
			Key:    key,
			Values: values,
			Kept:   kept.Value,
		})
	}

	return result, nil
}
