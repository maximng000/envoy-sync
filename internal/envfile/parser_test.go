package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "APP_ENV" || ef.Entries[0].Value != "production" {
		t.Errorf("unexpected entry: %+v", ef.Entries[0])
	}
}

func TestParse_CommentsSkipped(t *testing.T) {
	path := writeTempEnv(t, "# this is a comment\nFOO=bar\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if _, ok := m["FOO"]; !ok {
		t.Error("expected FOO key in map")
	}
	if len(m) != 1 {
		t.Errorf("expected 1 map entry, got %d", len(m))
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Entries[0].Value != "my secret value" {
		t.Errorf("expected unquoted value, got %q", ef.Entries[0].Value)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
