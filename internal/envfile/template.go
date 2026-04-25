package envfile

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// TemplateResult holds the output of a template rendering operation.
type TemplateResult struct {
	Rendered string
	Missing  []string
}

// RenderTemplate takes a template string with {{KEY}} placeholders and
// substitutes values from entries. If failOnMissing is true, missing keys
// return an error; otherwise they are left as empty strings and tracked.
func RenderTemplate(tmpl string, entries map[string]string, failOnMissing bool) (TemplateResult, error) {
	result := TemplateResult{}
	missingSet := map[string]struct{}{}

	output := tmpl
	for key, val := range entries {
		placeholder := "{{" + key + "}}"
		output = strings.ReplaceAll(output, placeholder, val)
	}

	// Detect any remaining placeholders
	for {
		start := strings.Index(output, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(output[start:], "}}")
		if end == -1 {
			break
		}
		key := output[start+2 : start+end]
		missingSet[key] = struct{}{}
		if failOnMissing {
			return result, fmt.Errorf("template key not found in env: %s", key)
		}
		output = output[:start] + output[start+end+2:]
	}

	for k := range missingSet {
		result.Missing = append(result.Missing, k)
	}
	sort.Strings(result.Missing)
	result.Rendered = output
	return result, nil
}

// RenderTemplateFile reads a template file from disk and renders it using
// the provided entries map.
func RenderTemplateFile(path string, entries map[string]string, failOnMissing bool) (TemplateResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TemplateResult{}, fmt.Errorf("reading template file: %w", err)
	}
	return RenderTemplate(string(data), entries, failOnMissing)
}
