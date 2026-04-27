package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// ScopeResult holds entries filtered to a specific scope.
type ScopeResult struct {
	Scope   string
	Entries []Entry
}

// ScopeSummary describes the outcome of a scope operation.
type ScopeSummary struct {
	Scopes []string
	Counts map[string]int
}

// Scope filters entries whose keys match the given scope prefix (e.g. "DB", "AWS").
// The prefix is matched case-insensitively against KEY_PREFIX patterns.
// Keys like DB_HOST, DB_PORT match scope "DB".
func Scope(entries []Entry, scope string) ScopeResult {
	scope = strings.ToUpper(strings.TrimSpace(scope))
	prefix := scope + "_"
	var matched []Entry
	for _, e := range entries {
		if strings.HasPrefix(strings.ToUpper(e.Key), prefix) || strings.ToUpper(e.Key) == scope {
			matched = append(matched, e)
		}
	}
	return ScopeResult{Scope: scope, Entries: matched}
}

// ListScopes returns all distinct top-level prefixes found in the entry keys.
// A prefix is the part of the key before the first underscore.
func ListScopes(entries []Entry) []string {
	seen := make(map[string]struct{})
	for _, e := range entries {
		parts := strings.SplitN(e.Key, "_", 2)
		if len(parts) == 2 && parts[0] != "" {
			seen[parts[0]] = struct{}{}
		}
	}
	scopes := make([]string, 0, len(seen))
	for s := range seen {
		scopes = append(scopes, s)
	}
	sort.Strings(scopes)
	return scopes
}

// ScopeSummaryOf builds a summary of how many keys belong to each scope.
func ScopeSummaryOf(entries []Entry) ScopeSummary {
	scopes := ListScopes(entries)
	counts := make(map[string]int, len(scopes))
	for _, s := range scopes {
		r := Scope(entries, s)
		counts[s] = len(r.Entries)
	}
	return ScopeSummary{Scopes: scopes, Counts: counts}
}

// FormatScopeSummary returns a human-readable summary string.
func FormatScopeSummary(s ScopeSummary) string {
	if len(s.Scopes) == 0 {
		return "no scopes detected"
	}
	var sb strings.Builder
	for _, sc := range s.Scopes {
		sb.WriteString(fmt.Sprintf("  %-20s %d key(s)\n", sc, s.Counts[sc]))
	}
	return strings.TrimRight(sb.String(), "\n")
}
