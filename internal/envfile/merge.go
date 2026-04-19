package envfile

// MergeStrategy defines how conflicts are resolved during merge.
type MergeStrategy int

const (
	// PreferBase keeps the base value on conflict.
	PreferBase MergeStrategy = iota
	// PreferOverride uses the override value on conflict.
	PreferOverride
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Merged   map[string]string
	Conflicts []MergeConflict
}

// MergeConflict describes a key that existed in both maps with different values.
type MergeConflict struct {
	Key          string
	BaseValue    string
	OverrideValue string
}

// Merge combines base and override maps according to the given strategy.
// Keys present only in base or only in override are always included.
// Conflicting keys are resolved by strategy and recorded in MergeResult.Conflicts.
func Merge(base, override map[string]string, strategy MergeStrategy) MergeResult {
	merged := make(map[string]string, len(base))
	var conflicts []MergeConflict

	for k, v := range base {
		merged[k] = v
	}

	for k, ov := range override {
		bv, exists := merged[k]
		if !exists {
			merged[k] = ov
			continue
		}
		if bv != ov {
			conflicts = append(conflicts, MergeConflict{
				Key:           k,
				BaseValue:     bv,
				OverrideValue: ov,
			})
			if strategy == PreferOverride {
				merged[k] = ov
			}
		}
	}

	return MergeResult{
		Merged:    merged,
		Conflicts: conflicts,
	}
}
