package envfile

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation issue.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors from a validation run.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) OK() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) Add(key, msg string) {
	r.Errors = append(r.Errors, ValidationError{Key: key, Message: msg})
}

// ValidateAgainstSchema checks that env satisfies all keys required by schema.
// schema is typically a .env.example map where values may be empty.
func ValidateAgainstSchema(env map[string]string, schema map[string]string) ValidationResult {
	var result ValidationResult
	for key := range schema {
		val, ok := env[key]
		if !ok {
			result.Add(key, "missing required key")
			continue
		}
		if strings.TrimSpace(val) == "" {
			result.Add(key, "value is empty")
		}
	}
	return result
}

// ValidateKeys checks for keys with empty values in env.
func ValidateKeys(env map[string]string) ValidationResult {
	var result ValidationResult
	for k, v := range env {
		if strings.TrimSpace(v) == "" {
			result.Add(k, "value is empty")
		}
	}
	return result
}
