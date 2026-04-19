package envfile

import "fmt"

// SyncResult holds the outcome of a sync operation.
type SyncResult struct {
	Applied  []string
	Skipped  []string
	Conflicts []string
}

// SyncStrategy controls how conflicts are resolved.
type SyncStrategy int

const (
	StrategySkip     SyncStrategy = iota // keep base value
	StrategyOverride                      // use source value
)

// Sync copies keys from src into dst according to the chosen strategy.
// Keys already present in dst are treated as conflicts.
func Sync(dst, src map[string]string, strategy SyncStrategy) (map[string]string, SyncResult) {
	result := map[string]string{}
	for k, v := range dst {
		result[k] = v
	}

	var sr SyncResult
	for k, v := range src {
		if existing, ok := dst[k]; ok {
			if existing == v {
				sr.Skipped = append(sr.Skipped, k)
				continue
			}
			sr.Conflicts = append(sr.Conflicts, fmt.Sprintf("%s: %q -> %q", k, existing, v))
			if strategy == StrategyOverride {
				result[k] = v
				sr.Applied = append(sr.Applied, k)
			} else {
				sr.Skipped = append(sr.Skipped, k)
			}
		} else {
			result[k] = v
			sr.Applied = append(sr.Applied, k)
		}
	}
	return result, sr
}
