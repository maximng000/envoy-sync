package envfile

import (
	"strings"
)

// NormalizeOptions controls how normalization is applied.
type NormalizeOptions struct {
	UppercaseKeys   bool
	TrimValues      bool
	RemoveEmpty     bool
	SortAlpha       bool
}

// NormalizeResult holds the outcome of a normalization pass.
type NormalizeResult struct {
	Entries    []Entry
	Modified   []string // keys that were changed
	Removed    []string // keys that were dropped
}

// Normalize applies the given options to a slice of entries and returns
// a NormalizeResult describing what changed.
func Normalize(entries []Entry, opts NormalizeOptions) NormalizeResult {
	result := NormalizeResult{}
	seen := make(map[string]bool)

	for _, e := range entries {
		origKey := e.Key
		origVal := e.Value

		if opts.UppercaseKeys {
			e.Key = strings.ToUpper(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.RemoveEmpty && e.Value == "" {
			result.Removed = append(result.Removed, origKey)
			continue
		}
		if seen[e.Key] {
			result.Removed = append(result.Removed, origKey)
			continue
		}
		seen[e.Key] = true

		if e.Key != origKey || e.Value != origVal {
			result.Modified = append(result.Modified, origKey)
		}
		result.Entries = append(result.Entries, e)
	}

	if opts.SortAlpha {
		result.Entries = Sort(result.Entries, SortOptions{Strategy: "alpha"}).Entries
	}

	return result
}

// NormalizeSummary returns a human-readable summary string.
func NormalizeSummary(r NormalizeResult) string {
	var sb strings.Builder
	sb.WriteString("Normalize summary:\n")
	sb.WriteString("  kept:     " + itoa(len(r.Entries)) + "\n")
	sb.WriteString("  modified: " + itoa(len(r.Modified)) + "\n")
	sb.WriteString("  removed:  " + itoa(len(r.Removed)) + "\n")
	return sb.String()
}

func itoa(n int) string {
	return strings.TrimSpace(strings.ReplaceAll(" "+string(rune('0'+n%10)), " ", ""))
}
