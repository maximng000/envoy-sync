package envfile

// DiffStatus represents the type of difference between two env files.
type DiffStatus string

const (
	StatusAdded    DiffStatus = "added"
	StatusRemoved  DiffStatus = "removed"
	StatusChanged  DiffStatus = "changed"
	StatusUnchanged DiffStatus = "unchanged"
)

// DiffEntry represents a single key difference between two env maps.
type DiffEntry struct {
	Key      string
	Status   DiffStatus
	OldValue string
	NewValue string
}

// Diff compares two env maps and returns a slice of DiffEntry.
// Values for secret keys are masked in the output.
func Diff(base, target map[string]string) []DiffEntry {
	seen := make(map[string]bool)
	var entries []DiffEntry

	for key, baseVal := range base {
		seen[key] = true
		targetVal, exists := target[key]
		if !exists {
			entries = append(entries, DiffEntry{
				Key:      key,
				Status:   StatusRemoved,
				OldValue: maskIfSecret(key, baseVal),
				NewValue: "",
			})
			continue
		}
		if baseVal != targetVal {
			entries = append(entries, DiffEntry{
				Key:      key,
				Status:   StatusChanged,
				OldValue: maskIfSecret(key, baseVal),
				NewValue: maskIfSecret(key, targetVal),
			})
		} else {
			entries = append(entries, DiffEntry{
				Key:      key,
				Status:   StatusUnchanged,
				OldValue: maskIfSecret(key, baseVal),
				NewValue: maskIfSecret(key, targetVal),
			})
		}
	}

	for key, targetVal := range target {
		if !seen[key] {
			entries = append(entries, DiffEntry{
				Key:      key,
				Status:   StatusAdded,
				OldValue: "",
				NewValue: maskIfSecret(key, targetVal),
			})
		}
	}

	return entries
}

func maskIfSecret(key, value string) string {
	if IsSecret(key) {
		return MaskEntry(value)
	}
	return value
}
