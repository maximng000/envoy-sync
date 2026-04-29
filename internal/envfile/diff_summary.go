package envfile

// DiffSummary holds aggregated counts from a diff result.
type DiffSummary struct {
	Added     int
	Removed   int
	Changed   int
	Unchanged int
	Total     int
	HasChanges bool
}

// SummarizeDiff aggregates a slice of DiffEntry values into a DiffSummary.
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
		case "unchanged":
			s.Unchanged++
		}
	}
	s.Total = s.Added + s.Removed + s.Changed + s.Unchanged
	s.HasChanges = s.Added > 0 || s.Removed > 0 || s.Changed > 0
	return s
}
