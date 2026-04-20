package envfile

import "fmt"

// RenameResult describes the outcome of a single key rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Renamed bool
	Reason  string
}

// RenameEntry renames oldKey to newKey in the provided entries map.
// It returns a RenameResult and an updated copy of the map.
// If oldKey does not exist, Renamed is false and the original map is returned unchanged.
// If newKey already exists and overwrite is false, Renamed is false.
func RenameEntry(entries map[string]string, oldKey, newKey string, overwrite bool) (map[string]string, RenameResult) {
	result := RenameResult{OldKey: oldKey, NewKey: newKey}

	if oldKey == newKey {
		result.Reason = "old and new key are identical"
		return copyMap(entries), result
	}

	value, exists := entries[oldKey]
	if !exists {
		result.Reason = fmt.Sprintf("key %q not found", oldKey)
		return copyMap(entries), result
	}

	if _, conflict := entries[newKey]; conflict && !overwrite {
		result.Reason = fmt.Sprintf("key %q already exists; use overwrite to replace", newKey)
		return copyMap(entries), result
	}

	updated := copyMap(entries)
	delete(updated, oldKey)
	updated[newKey] = value

	result.Renamed = true
	result.Reason = "ok"
	return updated, result
}

// BulkRename applies multiple renames sequentially.
// Each pair in renames is [oldKey, newKey].
// Returns the final map and a slice of results for each rename.
func BulkRename(entries map[string]string, renames [][2]string, overwrite bool) (map[string]string, []RenameResult) {
	current := copyMap(entries)
	results := make([]RenameResult, 0, len(renames))
	for _, pair := range renames {
		var res RenameResult
		current, res = RenameEntry(current, pair[0], pair[1], overwrite)
		results = append(results, res)
	}
	return current, results
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
