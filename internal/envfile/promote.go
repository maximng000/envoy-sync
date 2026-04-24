package envfile

import "fmt"

// PromoteResult describes what happened to a single key during promotion.
type PromoteResult struct {
	Key       string
	Action    string // "added", "updated", "skipped", "unchanged"
	OldValue  string
	NewValue  string
	IsSecret  bool
}

// PromoteOptions controls promotion behaviour.
type PromoteOptions struct {
	Overwrite bool // overwrite existing keys in dst
	DryRun    bool // do not mutate dst
	MaskSecrets bool
}

// Promote copies keys from src into dst, returning a summary of actions taken.
// Keys present only in src are always added. Keys present in both are
// overwritten only when opts.Overwrite is true.
func Promote(src, dst []Entry, opts PromoteOptions) ([]Entry, []PromoteResult, error) {
	if src == nil {
		return nil, nil, fmt.Errorf("promote: src is nil")
	}

	dstMap := make(map[string]string, len(dst))
	for _, e := range dst {
		dstMap[e.Key] = e.Value
	}

	results := make([]PromoteResult, 0, len(src))
	outMap := make(map[string]string, len(dst))
	for k, v := range dstMap {
		outMap[k] = v
	}

	for _, e := range src {
		secret := IsSecret(e.Key)
		old, exists := dstMap[e.Key]
		displayNew := e.Value
		displayOld := old
		if opts.MaskSecrets && secret {
			displayNew = MaskEntry(e).Value
			displayOld = MaskEntry(Entry{Key: e.Key, Value: old}).Value
		}

		switch {
		case !exists:
			results = append(results, PromoteResult{Key: e.Key, Action: "added", NewValue: displayNew, IsSecret: secret})
			if !opts.DryRun {
				outMap[e.Key] = e.Value
			}
		case old == e.Value:
			results = append(results, PromoteResult{Key: e.Key, Action: "unchanged", OldValue: displayOld, NewValue: displayNew, IsSecret: secret})
		case opts.Overwrite:
			results = append(results, PromoteResult{Key: e.Key, Action: "updated", OldValue: displayOld, NewValue: displayNew, IsSecret: secret})
			if !opts.DryRun {
				outMap[e.Key] = e.Value
			}
		default:
			results = append(results, PromoteResult{Key: e.Key, Action: "skipped", OldValue: displayOld, NewValue: displayNew, IsSecret: secret})
		}
	}

	out := make([]Entry, 0, len(outMap))
	for k, v := range outMap {
		out = append(out, Entry{Key: k, Value: v})
	}
	return out, results, nil
}
