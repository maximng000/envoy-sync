package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

var interpolationPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls interpolation behaviour.
type InterpolateOptions struct {
	// FailOnMissing returns an error if a referenced variable is not found.
	FailOnMissing bool
}

// Interpolate resolves variable references within entry values using the
// provided env map as the source of truth. References can be in the form
// $VAR or ${VAR}. Entries are processed in order; earlier definitions are
// visible to later ones.
func Interpolate(entries []Entry, opts InterpolateOptions) ([]Entry, error) {
	resolved := make(map[string]string, len(entries))
	result := make([]Entry, 0, len(entries))

	for _, e := range entries {
		expanded, err := expand(e.Value, resolved, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolating %q: %w", e.Key, err)
		}
		resolved[e.Key] = expanded
		result = append(result, Entry{Key: e.Key, Value: expanded})
	}

	return result, nil
}

func expand(value string, env map[string]string, opts InterpolateOptions) (string, error) {
	var expandErr error

	result := interpolationPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}

		key := strings.TrimPrefix(match, "$")
		key = strings.TrimPrefix(key, "{")
		key = strings.TrimSuffix(key, "}")

		if v, ok := env[key]; ok {
			return v
		}

		if opts.FailOnMissing {
			expandErr = fmt.Errorf("undefined variable %q", key)
			return match
		}

		return ""
	})

	if expandErr != nil {
		return "", expandErr
	}

	return result, nil
}
