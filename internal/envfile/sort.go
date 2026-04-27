package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// SortStrategy defines how entries should be sorted.
type SortStrategy string

const (
	SortAlpha     SortStrategy = "alpha"
	SortAlphaDesc SortStrategy = "alpha-desc"
	SortBySecret  SortStrategy = "secret"
	SortByLength  SortStrategy = "length"
)

// SortResult holds sorted entries and metadata.
type SortResult struct {
	Entries  []Entry
	Strategy SortStrategy
	Total    int
}

// Sort returns entries ordered by the given strategy.
func Sort(entries []Entry, strategy SortStrategy) (SortResult, error) {
	copied := make([]Entry, len(entries))
	copy(copied, entries)

	switch strategy {
	case SortAlpha, "":
		sort.Slice(copied, func(i, j int) bool {
			return strings.ToLower(copied[i].Key) < strings.ToLower(copied[j].Key)
		})
	case SortAlphaDesc:
		sort.Slice(copied, func(i, j int) bool {
			return strings.ToLower(copied[i].Key) > strings.ToLower(copied[j].Key)
		})
	case SortBySecret:
		sort.SliceStable(copied, func(i, j int) bool {
			iSecret := IsSecret(copied[i].Key)
			jSecret := IsSecret(copied[j].Key)
			if iSecret != jSecret {
				return iSecret
			}
			return strings.ToLower(copied[i].Key) < strings.ToLower(copied[j].Key)
		})
	case SortByLength:
		sort.SliceStable(copied, func(i, j int) bool {
			if len(copied[i].Key) != len(copied[j].Key) {
				return len(copied[i].Key) < len(copied[j].Key)
			}
			return strings.ToLower(copied[i].Key) < strings.ToLower(copied[j].Key)
		})
	default:
		return SortResult{}, fmt.Errorf("unknown sort strategy: %q", strategy)
	}

	return SortResult{
		Entries:  copied,
		Strategy: strategy,
		Total:    len(copied),
	}, nil
}
