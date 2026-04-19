package envfile

import (
	"strings"
	"testing"
)

func TestAudit_AddedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	updated := map[string]string{"FOO": "bar", "NEW_KEY": "value"}
	entries := Audit(base, updated)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != ActionAdded || entries[0].Key != "NEW_KEY" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestAudit_RemovedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "gone"}
	updated := map[string]string{"FOO": "bar"}
	entries := Audit(base, updated)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != ActionRemoved || entries[0].Key != "OLD" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestAudit_ChangedKey(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	updated := map[string]string{"FOO": "new"}
	entries := Audit(base, updated)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != ActionChanged || e.OldValue != "old" || e.NewValue != "new" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestAudit_NoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	entries := Audit(base, base)
	if len(entries) != 0 {
		t.Errorf("expected no entries, got %d", len(entries))
	}
}

func TestAudit_SecretMaskedInString(t *testing.T) {
	base := map[string]string{"DB_PASSWORD": "old_secret"}
	updated := map[string]string{"DB_PASSWORD": "new_secret"}
	entries := Audit(base, updated)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	s := entries[0].String()
	if strings.Contains(s, "old_secret") || strings.Contains(s, "new_secret") {
		t.Errorf("secret value exposed in audit string: %s", s)
	}
	if !strings.Contains(s, "***") {
		t.Errorf("expected masked value in string: %s", s)
	}
}
