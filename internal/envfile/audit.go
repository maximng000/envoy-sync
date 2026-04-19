package envfile

import (
	"fmt"
	"time"
)

// AuditAction represents the type of change recorded in an audit entry.
type AuditAction string

const (
	ActionAdded   AuditAction = "added"
	ActionRemoved AuditAction = "removed"
	ActionChanged AuditAction = "changed"
)

// AuditEntry records a single change event for a key.
type AuditEntry struct {
	Timestamp time.Time
	Key       string
	Action    AuditAction
	OldValue  string
	NewValue  string
	Secret    bool
}

// String returns a human-readable representation of the audit entry.
func (e AuditEntry) String() string {
	old := e.OldValue
	new_ := e.NewValue
	if e.Secret {
		old = mask(old)
		new_ = mask(new_)
	}
	switch e.Action {
	case ActionAdded:
		return fmt.Sprintf("[%s] ADD %s = %q", e.Timestamp.Format(time.RFC3339), e.Key, new_)
	case ActionRemoved:
		return fmt.Sprintf("[%s] DEL %s (was %q)", e.Timestamp.Format(time.RFC3339), e.Key, old)
	case ActionChanged:
		return fmt.Sprintf("[%s] MOD %s: %q -> %q", e.Timestamp.Format(time.RFC3339), e.Key, old, new_)
	}
	return ""
}

func mask(v string) string {
	if len(v) == 0 {
		return ""
	}
	return "***"
}

// Audit compares base and updated maps and returns a log of changes.
func Audit(base, updated map[string]string) []AuditEntry {
	now := time.Now().UTC()
	var entries []AuditEntry

	for k, newVal := range updated {
		oldVal, exists := base[k]
		if !exists {
			entries = append(entries, AuditEntry{
				Timestamp: now, Key: k, Action: ActionAdded,
				NewValue: newVal, Secret: IsSecret(k),
			})
		} else if oldVal != newVal {
			entries = append(entries, AuditEntry{
				Timestamp: now, Key: k, Action: ActionChanged,
				OldValue: oldVal, NewValue: newVal, Secret: IsSecret(k),
			})
		}
	}

	for k, oldVal := range base {
		if _, exists := updated[k]; !exists {
			entries = append(entries, AuditEntry{
				Timestamp: now, Key: k, Action: ActionRemoved,
				OldValue: oldVal, Secret: IsSecret(k),
			})
		}
	}

	return entries
}
