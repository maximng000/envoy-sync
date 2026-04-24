package envfile

import (
	"fmt"
	"strings"
)

// RedactResult holds the output of a redaction operation.
type RedactResult struct {
	Entries  []Entry
	Redacted []string // keys that were redacted
}

// RedactMode controls how secret values are replaced.
type RedactMode string

const (
	RedactModeBlank    RedactMode = "blank"    // replace with empty string
	RedactModeMask     RedactMode = "mask"     // replace with ***
	RedactModePlaceholder RedactMode = "placeholder" // replace with {{KEY}}
)

// Redact returns a copy of entries with secret values replaced according to mode.
// If keys is non-empty, only those keys are redacted regardless of IsSecret.
func Redact(entries []Entry, mode RedactMode, keys []string) RedactResult {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[strings.ToUpper(k)] = true
	}

	forceByKey := len(keySet) > 0

	out := make([]Entry, len(entries))
	var redacted []string

	for i, e := range entries {
		copy := e
		shouldRedact := (forceByKey && keySet[strings.ToUpper(e.Key)]) ||
			(!forceByKey && IsSecret(e.Key))

		if shouldRedact {
			copy.Value = redactValue(e.Key, mode)
			redacted = append(redacted, e.Key)
		}
		out[i] = copy
	}

	return RedactResult{Entries: out, Redacted: redacted}
}

func redactValue(key string, mode RedactMode) string {
	switch mode {
	case RedactModeBlank:
		return ""
	case RedactModePlaceholder:
		return fmt.Sprintf("{{%s}}", strings.ToUpper(key))
	default: // RedactModeMask
		return "***"
	}
}
