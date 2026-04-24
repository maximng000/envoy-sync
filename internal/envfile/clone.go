package envfile

import (
	"fmt"
	"sort"
)

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	Key     string
	OldEnv  string
	NewEnv  string
	Skipped bool
	Reason  string
}

// CloneOptions controls behaviour during cloning.
type CloneOptions struct {
	// FilterKeys, if non-empty, clones only the listed keys.
	FilterKeys []string
	// Overwrite allows existing keys in dst to be overwritten.
	Overwrite bool
	// MaskSecrets redacts secret values in the returned results (does not
	// affect the actual entries written to dst).
	MaskSecrets bool
}

// Clone copies entries from src into dst, returning one CloneResult per key
// that was considered. Keys already present in dst are skipped unless
// Overwrite is set.
func Clone(src, dst map[string]string, opts CloneOptions) (map[string]string, []CloneResult) {
	out := copyMap(dst)

	keys := keysToClone(src, opts.FilterKeys)
	results := make([]CloneResult, 0, len(keys))

	for _, k := range keys {
		v := src[k]
		display := v
		if opts.MaskSecrets && IsSecret(k) {
			display = MaskEntry(k, v)
		}

		if _, exists := out[k]; exists && !opts.Overwrite {
			results = append(results, CloneResult{
				Key:     k,
				Skipped: true,
				Reason:  fmt.Sprintf("key already exists in destination (value: %s)", display),
			})
			continue
		}

		out[k] = v
		results = append(results, CloneResult{
			Key:    k,
			OldEnv: display,
			NewEnv: display,
		})
	}

	return out, results
}

func keysToClone(src map[string]string, filter []string) []string {
	if len(filter) > 0 {
		sorted := make([]string, len(filter))
		copy(sorted, filter)
		sort.Strings(sorted)
		return sorted
	}
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
