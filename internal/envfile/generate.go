package envfile

import (
	"fmt"
	"strings"
)

// GenerateOptions controls how a template .env file is generated.
type GenerateOptions struct {
	IncludeComments bool
	Placeholder     string // default: "CHANGEME"
}

// GenerateResult holds the output of a Generate call.
type GenerateResult struct {
	Lines []string
	Count int
}

// Generate produces a template .env file from a list of key names.
// Secret keys receive a masked placeholder; others use the configured placeholder.
func Generate(keys []string, opts GenerateOptions) GenerateResult {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "CHANGEME"
	}

	var lines []string
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if opts.IncludeComments {
			if IsSecret(key) {
				lines = append(lines, fmt.Sprintf("# %s — secret value, handle with care", key))
			} else {
				lines = append(lines, fmt.Sprintf("# %s", key))
			}
		}
		value := placeholder
		if IsSecret(key) {
			value = "***"
		}
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	return GenerateResult{
		Lines: lines,
		Count: len(keys),
	}
}

// GenerateFromEntries builds a template from an existing parsed env map,
// stripping all values and replacing them with placeholders.
func GenerateFromEntries(entries map[string]string, opts GenerateOptions) GenerateResult {
	keys := sortedKeys(entries)
	return Generate(keys, opts)
}
