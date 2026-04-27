package envfile

import "fmt"

// DiffSummary holds aggregated statistics about a diff result.
type DiffSummary struct {
	Added   int
	Removed int
	Changed int
	Total   int
}

// String returns a human-readable summary line.
func (s DiffSummary) String() string {
	return fmt.Sprintf("total=%d added=%d removed=%d changed=%d",
		s.Total, s.Added, s.Removed, s.Changed)
}

// HasChanges returns true when any difference was detected.
func (s DiffSummary) HasChanges() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}

// SummarizeDiff computes a DiffSummary from a slice of DiffEntry values
// as returned by Diff.
func SummarizeDiff(entries []DiffEntry) DiffSummary {
	var s DiffSummary
	for _, e := range entries {
		switch e.Status {
		case "added":
			s.Added++
		case "removed":
			s.Removed++
		case "changed":
			s.Changed++
		}
	}
	s.Total = s.Added + s.Removed + s.Changed
	return s
}
