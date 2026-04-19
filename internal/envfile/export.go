package envfile

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents an output format for env file export.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Export writes the given env map to w in the specified format.
func Export(w io.Writer, entries map[string]string, format Format, maskSecrets bool) error {
	keys := sortedKeys(entries)

	switch format {
	case FormatDotenv:
		for _, k := range keys {
			v := valueOf(k, entries, maskSecrets)
			fmt.Fprintf(w, "%s=%s\n", k, quoteIfNeeded(v))
		}
	case FormatExport:
		for _, k := range keys {
			v := valueOf(k, entries, maskSecrets)
			fmt.Fprintf(w, "export %s=%s\n", k, quoteIfNeeded(v))
		}
	case FormatJSON:
		return exportJSON(w, keys, entries, maskSecrets)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	return nil
}

func valueOf(key string, entries map[string]string, maskSecrets bool) string {
	v := entries[key]
	if maskSecrets && IsSecret(key) {
		return "***"
	}
	return v
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t#") {
		return `"` + v + `"`
	}
	return v
}

func exportJSON(w io.Writer, keys []string, entries map[string]string, maskSecrets bool) error {
	fmt.Fprintln(w, "{")
	for i, k := range keys {
		v := valueOf(k, entries, maskSecrets)
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(w, "  %q: %q%s\n", k, v, comma)
	}
	fmt.Fprintln(w, "}")
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
