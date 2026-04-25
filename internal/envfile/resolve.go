package envfile

import (
	"fmt"
	"os"
	"strings"
)

// ResolveSource describes where a value came from during resolution.
type ResolveSource string

const (
	SourceFile   ResolveSource = "file"
	SourceEnv    ResolveSource = "env"
	SourceDefault ResolveSource = "default"
	SourceMissing ResolveSource = "missing"
)

// ResolvedEntry holds a key, its resolved value, and the source of that value.
type ResolvedEntry struct {
	Key    string
	Value  string
	Source ResolveSource
}

// ResolveOptions controls how resolution behaves.
type ResolveOptions struct {
	// Defaults provides fallback values when a key is absent from file and env.
	Defaults map[string]string
	// FailOnMissing causes Resolve to return an error if any key has no value.
	FailOnMissing bool
	// PreferEnv makes live environment variables take precedence over file values.
	PreferEnv bool
}

// Resolve merges file entries with live environment variables and optional
// defaults. Resolution priority (when PreferEnv is false): file > env > default.
// When PreferEnv is true: env > file > default.
func Resolve(entries []Entry, opts ResolveOptions) ([]ResolvedEntry, error) {
	results := make([]ResolvedEntry, 0, len(entries))

	for _, e := range entries {
		envVal, inEnv := os.LookupEnv(e.Key)

		var resolved ResolvedEntry
		resolved.Key = e.Key

		switch {
		case opts.PreferEnv && inEnv:
			resolved.Value = envVal
			resolved.Source = SourceEnv
		case e.Value != "":
			resolved.Value = e.Value
			resolved.Source = SourceFile
		case inEnv:
			resolved.Value = envVal
			resolved.Source = SourceEnv
		case opts.Defaults != nil:
			if def, ok := opts.Defaults[e.Key]; ok {
				resolved.Value = def
				resolved.Source = SourceDefault
			} else {
				resolved.Source = SourceMissing
			}
		default:
			resolved.Source = SourceMissing
		}

		if opts.FailOnMissing && resolved.Source == SourceMissing {
			return nil, fmt.Errorf("resolve: no value found for key %q", e.Key)
		}

		results = append(results, resolved)
	}

	return results, nil
}

// ResolvedToEntries converts resolved entries back to a plain Entry slice,
// omitting any entries whose source is SourceMissing.
func ResolvedToEntries(resolved []ResolvedEntry) []Entry {
	out := make([]Entry, 0, len(resolved))
	for _, r := range resolved {
		if r.Source == SourceMissing {
			continue
		}
		out = append(out, Entry{Key: r.Key, Value: r.Value})
	}
	return out
}

// ResolvedSummary returns a human-readable summary line for each resolved entry.
func ResolvedSummary(resolved []ResolvedEntry) string {
	var sb strings.Builder
	for _, r := range resolved {
		val := r.Value
		if IsSecret(r.Key) {
			val = MaskEntry(Entry{Key: r.Key, Value: val}).Value
		}
		fmt.Fprintf(&sb, "%-30s %-12s %s\n", r.Key, "["+string(r.Source)+"]", val)
	}
	return sb.String()
}
