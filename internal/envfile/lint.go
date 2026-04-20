package envfile

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting problem found in an env file.
type LintIssue struct {
	Line    int
	Key     string
	Message string
	Severity string // "warn" or "error"
}

func (l LintIssue) String() string {
	return fmt.Sprintf("[%s] line %d (%s): %s", l.Severity, l.Line, l.Key, l.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}

// Lint checks an env file's entries for common style and correctness issues.
func Lint(entries []Entry) *LintResult {
	result := &LintResult{}
	seen := make(map[string]int)

	for idx, entry := range entries {
		lineNum := idx + 1

		// Duplicate key check
		if prev, ok := seen[entry.Key]; ok {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  fmt.Sprintf("duplicate key, first seen at line %d", prev),
				Severity: "error",
			})
		}
		seen[entry.Key] = lineNum

		// Key naming convention: should be UPPER_SNAKE_CASE
		if entry.Key != strings.ToUpper(entry.Key) {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  "key should be UPPER_SNAKE_CASE",
				Severity: "warn",
			})
		}

		// Empty value warning
		if strings.TrimSpace(entry.Value) == "" {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  "value is empty",
				Severity: "warn",
			})
		}

		// Leading/trailing whitespace in value
		if entry.Value != strings.TrimSpace(entry.Value) {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  "value has leading or trailing whitespace",
				Severity: "warn",
			})
		}
	}

	return result
}
