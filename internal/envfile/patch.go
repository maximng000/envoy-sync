package envfile

import "fmt"

// PatchOp represents a single patch operation on an env entry.
type PatchOp struct {
	Key    string
	Op     string // "set", "delete", "rename"
	Value  string // used by "set"
	NewKey string // used by "rename"
}

// PatchResult describes the outcome of a single patch operation.
type PatchResult struct {
	Op      string
	Key     string
	Applied bool
	Reason  string
}

// Patch applies a list of PatchOps to a copy of the given entries.
// It returns the updated entries and a slice of PatchResults.
func Patch(entries []Entry, ops []PatchOp) ([]Entry, []PatchResult, error) {
	result := make([]Entry, len(entries))
	copy(result, entries)

	var results []PatchResult

	for _, op := range ops {
		switch op.Op {
		case "set":
			found := false
			for i, e := range result {
				if e.Key == op.Key {
					result[i].Value = op.Value
					found = true
					break
				}
			}
			if !found {
				result = append(result, Entry{Key: op.Key, Value: op.Value})
			}
			results = append(results, PatchResult{Op: "set", Key: op.Key, Applied: true})

		case "delete":
			next := result[:0]
			deleted := false
			for _, e := range result {
				if e.Key == op.Key {
					deleted = true
					continue
				}
				next = append(next, e)
			}
			result = next
			if deleted {
				results = append(results, PatchResult{Op: "delete", Key: op.Key, Applied: true})
			} else {
				results = append(results, PatchResult{Op: "delete", Key: op.Key, Applied: false, Reason: "key not found"})
			}

		case "rename":
			if op.NewKey == "" {
				return nil, nil, fmt.Errorf("rename op for %q requires a non-empty new_key", op.Key)
			}
			found := false
			for i, e := range result {
				if e.Key == op.Key {
					result[i].Key = op.NewKey
					found = true
					break
				}
			}
			if found {
				results = append(results, PatchResult{Op: "rename", Key: op.Key, Applied: true})
			} else {
				results = append(results, PatchResult{Op: "rename", Key: op.Key, Applied: false, Reason: "key not found"})
			}

		default:
			return nil, nil, fmt.Errorf("unknown patch op %q for key %q", op.Op, op.Key)
		}
	}

	return result, results, nil
}
