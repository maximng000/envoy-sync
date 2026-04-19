package envfile

import (
	"testing"
)

func TestDiff_AddedKey(t *testing.T) {
	base := map[string]string{"APP_NAME": "myapp"}
	target := map[string]string{"APP_NAME": "myapp", "NEW_KEY": "value"}

	entries := Diff(base, target)
	found := findEntry(entries, "NEW_KEY")
	if found == nil {
		t.Fatal("expected NEW_KEY in diff")
	}
	if found.Status != StatusAdded {
		t.Errorf("expected added, got %s", found.Status)
	}
	if found.NewValue != "value" {
		t.Errorf("unexpected new value: %s", found.NewValue)
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := map[string]string{"APP_NAME": "myapp", "OLD_KEY": "old"}
	target := map[string]string{"APP_NAME": "myapp"}

	entries := Diff(base, target)
	found := findEntry(entries, "OLD_KEY")
	if found == nil {
		t.Fatal("expected OLD_KEY in diff")
	}
	if found.Status != StatusRemoved {
		t.Errorf("expected removed, got %s", found.Status)
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	base := map[string]string{"APP_NAME": "myapp"}
	target := map[string]string{"APP_NAME": "newapp"}

	entries := Diff(base, target)
	found := findEntry(entries, "APP_NAME")
	if found == nil {
		t.Fatal("expected APP_NAME in diff")
	}
	if found.Status != StatusChanged {
		t.Errorf("expected changed, got %s", found.Status)
	}
	if found.OldValue != "myapp" || found.NewValue != "newapp" {
		t.Errorf("unexpected values: old=%s new=%s", found.OldValue, found.NewValue)
	}
}

func TestDiff_SecretValueMasked(t *testing.T) {
	base := map[string]string{"DB_PASSWORD": "secret123"}
	target := map[string]string{"DB_PASSWORD": "newsecret"}

	entries := Diff(base, target)
	found := findEntry(entries, "DB_PASSWORD")
	if found == nil {
		t.Fatal("expected DB_PASSWORD in diff")
	}
	if found.OldValue == "secret123" || found.NewValue == "newsecret" {
		t.Error("secret values should be masked")
	}
}

func TestDiff_UnchangedKey(t *testing.T) {
	base := map[string]string{"APP_ENV": "production"}
	target := map[string]string{"APP_ENV": "production"}

	entries := Diff(base, target)
	found := findEntry(entries, "APP_ENV")
	if found == nil {
		t.Fatal("expected APP_ENV in diff")
	}
	if found.Status != StatusUnchanged {
		t.Errorf("expected unchanged, got %s", found.Status)
	}
}

func findEntry(entries []DiffEntry, key string) *DiffEntry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}
