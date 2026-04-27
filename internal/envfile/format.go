package envfile

import (
	"fmt"
	"strings"
)

// FormatStyle controls how entries are formatted.
type FormatStyle string

const (
	StyleAligned  FormatStyle = "aligned"
	StyleCompact  FormatStyle = "compact"
	StyleSpaced   FormatStyle = "spaced"
)

// FormatOptions configures the Format operation.
type FormatOptions struct {
	Style      FormatStyle
	SortKeys   bool
	MaskSecret bool
}

// FormatResult holds the formatted output and metadata.
type FormatResult struct {
	Lines    []string
	Modified int
}

// Format applies a consistent style to env entries and returns formatted lines.
func Format(entries []Entry, opts FormatOptions) FormatResult {
	if opts.Style == "" {
		opts.Style = StyleCompact
	}

	src := entries
	if opts.SortKeys {
		sorted := make([]Entry, len(entries))
		copy(sorted, entries)
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].Key > sorted[j].Key {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		src = sorted
	}

	maxLen := 0
	if opts.Style == StyleAligned {
		for _, e := range src {
			if len(e.Key) > maxLen {
				maxLen = len(e.Key)
			}
		}
	}

	var lines []string
	modified := 0

	for _, e := range src {
		val := e.Value
		if opts.MaskSecret && IsSecret(e.Key) {
			val = "***"
		}

		var line string
		switch opts.Style {
		case StyleAligned:
			padding := strings.Repeat(" ", maxLen-len(e.Key))
			line = fmt.Sprintf("%s%s = %s", e.Key, padding, val)
		case StyleSpaced:
			line = fmt.Sprintf("%s = %s", e.Key, val)
		default: // compact
			line = fmt.Sprintf("%s=%s", e.Key, val)
		}

		original := fmt.Sprintf("%s=%s", e.Key, e.Value)
		if line != original {
			modified++
		}
		lines = append(lines, line)
	}

	return FormatResult{Lines: lines, Modified: modified}
}
