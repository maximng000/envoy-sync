package envfile

import "strings"

// FilterOptions controls how entries are filtered.
type FilterOptions struct {
	Prefix    string
	Suffix    string
	Keys      []string
	SecretsOnly bool
	NonSecretsOnly bool
}

// Filter returns a subset of entries based on the provided options.
// If multiple options are set, all conditions must match (AND logic).
func Filter(entries []Entry, opts FilterOptions) []Entry {
	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	var result []Entry
	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if opts.Suffix != "" && !strings.HasSuffix(e.Key, opts.Suffix) {
			continue
		}
		if len(keySet) > 0 && !keySet[e.Key] {
			continue
		}
		if opts.SecretsOnly && !IsSecret(e.Key) {
			continue
		}
		if opts.NonSecretsOnly && IsSecret(e.Key) {
			continue
		}
		result = append(result, e)
	}
	return result
}

// FilterByPrefix returns entries whose keys start with the given prefix.
func FilterByPrefix(entries []Entry, prefix string) []Entry {
	return Filter(entries, FilterOptions{Prefix: prefix})
}

// FilterSecrets returns only entries identified as secrets.
func FilterSecrets(entries []Entry) []Entry {
	return Filter(entries, FilterOptions{SecretsOnly: true})
}

// FilterNonSecrets returns only entries not identified as secrets.
func FilterNonSecrets(entries []Entry) []Entry {
	return Filter(entries, FilterOptions{NonSecretsOnly: true})
}
