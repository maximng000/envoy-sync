package envfile

// DiffSummaryResult holds aggregated counts from a slice of DiffEntry.
type DiffSummaryResult struct {
	Added     int
	Removed   int
	Changed   int
	Unchanged int
	Total     int
	HasChanges bool
}

// DiffStatus represents the kind of change in a diff entry.
type DiffStatus string

const (
	DiffAdded     DiffStatus = "added"
	DiffRemoved   DiffStatus = "removed"
	DiffChanged   DiffStatus = "changed"
	DiffUnchanged DiffStatus = "unchanged"
)

// DiffSummaryEntry is a single entry used for summarisation.
type DiffSummaryEntry struct {
	Key    string
	Status DiffStatus
}

// SummarizeDiff aggregates a slice of DiffSummaryEntry into a DiffSummaryResult.
func SummarizeDiff(entries []DiffSummaryEntry) DiffSummaryResult {
	var r DiffSummaryResult
	for _, e := range entries {
		r.Total++
		switch e.Status {
		case DiffAdded:
			r.Added++
		case DiffRemoved:
			r.Removed++
		case DiffChanged:
			r.Changed++
		case DiffUnchanged:
			r.Unchanged++
		}
	}
	r.HasChanges = r.Added > 0 || r.Removed > 0 || r.Changed > 0
	return r
}
