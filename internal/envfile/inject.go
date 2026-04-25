package envfile

import (
	"fmt"
	"os"
	"sort"
)

// InjectResult holds the outcome of injecting env entries into the process environment.
type InjectResult struct {
	Injected []string
	Skipped  []string
}

// InjectOptions controls behaviour of Inject.
type InjectOptions struct {
	// Overwrite existing OS env vars when true.
	Overwrite bool
	// Keys restricts injection to only the listed keys. Empty means all keys.
	Keys []string
}

// Inject sets entries as environment variables in the current process.
// It returns an InjectResult describing what was injected or skipped.
func Inject(entries map[string]string, opts InjectOptions) (InjectResult, error) {
	filter := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		filter[k] = struct{}{}
	}

	keys := sortedInjectKeys(entries)
	var result InjectResult

	for _, k := range keys {
		if len(filter) > 0 {
			if _, ok := filter[k]; !ok {
				continue
			}
		}

		_, exists := os.LookupEnv(k)
		if exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if err := os.Setenv(k, entries[k]); err != nil {
			return result, fmt.Errorf("inject: failed to set %q: %w", k, err)
		}
		result.Injected = append(result.Injected, k)
	}

	return result, nil
}

func sortedInjectKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
