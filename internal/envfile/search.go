package envfile

import (
	"regexp"
	"strings"
)

// SearchMode controls how key/value matching is performed.
type SearchMode string

const (
	SearchModeExact  SearchMode = "exact"
	SearchModePrefix SearchMode = "prefix"
	SearchModeRegex  SearchMode = "regex"
)

// SearchOptions configures a Search operation.
type SearchOptions struct {
	// Mode determines how the Query is applied (default: exact).
	Mode SearchMode

	// Query is the search term applied to keys and/or values.
	Query string

	// SearchKeys enables matching against entry keys.
	SearchKeys bool

	// SearchValues enables matching against entry values.
	SearchValues bool

	// CaseSensitive controls whether matching is case-sensitive.
	CaseSensitive bool

	// MaskSecrets redacts secret values in the returned results.
	MaskSecrets bool
}

// SearchResult holds a matched entry along with metadata about the match.
type SearchResult struct {
	Entry    Entry
	MatchedOn string // "key", "value", or "both"
}

// Search scans a slice of Entry values and returns those that match the
// provided SearchOptions. It supports exact, prefix, and regex matching
// against keys, values, or both.
func Search(entries []Entry, opts SearchOptions) ([]SearchResult, error) {
	if opts.Mode == "" {
		opts.Mode = SearchModeExact
	}
	if !opts.SearchKeys && !opts.SearchValues {
		opts.SearchKeys = true
	}

	query := opts.Query
	if !opts.CaseSensitive && opts.Mode != SearchModeRegex {
		query = strings.ToLower(query)
	}

	var re *regexp.Regexp
	if opts.Mode == SearchModeRegex {
		flags := ""
		if !opts.CaseSensitive {
			flags = "(?i)"
		}
		var err error
		re, err = regexp.Compile(flags + opts.Query)
		if err != nil {
			return nil, err
		}
	}

	var results []SearchResult

	for _, e := range entries {
		keyMatch := false
		valMatch := false

		if opts.SearchKeys {
			keyMatch = matches(e.Key, query, opts.Mode, opts.CaseSensitive, re)
		}
		if opts.SearchValues {
			valMatch = matches(e.Value, query, opts.Mode, opts.CaseSensitive, re)
		}

		if !keyMatch && !valMatch {
			continue
		}

		matched := "key"
		if keyMatch && valMatch {
			matched = "both"
		} else if valMatch {
			matched = "value"
		}

		result := e
		if opts.MaskSecrets && IsSecret(e.Key) {
			result = MaskEntry(e)
		}

		results = append(results, SearchResult{
			Entry:     result,
			MatchedOn: matched,
		})
	}

	return results, nil
}

// matches checks whether target satisfies the query under the given mode.
func matches(target, query string, mode SearchMode, caseSensitive bool, re *regexp.Regexp) bool {
	compare := target
	if !caseSensitive && mode != SearchModeRegex {
		compare = strings.ToLower(target)
	}

	switch mode {
	case SearchModeExact:
		return compare == query
	case SearchModePrefix:
		return strings.HasPrefix(compare, query)
	case SearchModeRegex:
		if re == nil {
			return false
		}
		return re.MatchString(target)
	default:
		return compare == query
	}
}
