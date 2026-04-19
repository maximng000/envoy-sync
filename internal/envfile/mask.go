package envfile

import "strings"

// DefaultSecretPatterns contains substrings that indicate a key holds a secret.
var DefaultSecretPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
}

const maskedValue = "***"

// IsSecret reports whether a key is considered sensitive based on known patterns.
func IsSecret(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// MaskEntry returns the entry value, replacing it with "***" if the key is secret.
func MaskEntry(e Entry, patterns []string) Entry {
	if e.Key != "" && IsSecret(e.Key, patterns) {
		return Entry{Key: e.Key, Value: maskedValue, Comment: e.Comment}
	}
	return e
}

// MaskedMap returns a copy of the env map with secret values replaced.
func MaskedMap(m map[string]string, patterns []string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if IsSecret(k, patterns) {
			out[k] = maskedValue
		} else {
			out[k] = v
		}
	}
	return out
}
