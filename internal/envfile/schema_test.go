package envfile

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTempSchema(t *testing.T, s Schema) string {
	t.Helper()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(t.TempDir(), "schema-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadSchema_ValidFile(t *testing.T) {
	s := Schema{Fields: []SchemaField{{Key: "PORT", Required: true}}}
	path := writeTempSchema(t, s)
	loaded, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Fields) != 1 || loaded.Fields[0].Key != "PORT" {
		t.Errorf("unexpected fields: %+v", loaded.Fields)
	}
}

func TestLoadSchema_MissingFile(t *testing.T) {
	_, err := LoadSchema("/nonexistent/schema.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheckSchema_AllPresent(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "HOST", Required: true}}}
	entries := []Entry{{Key: "HOST", Value: "localhost"}}
	violations := CheckSchema(entries, s)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCheckSchema_MissingRequired(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "DB_URL", Required: true}}}
	violations := CheckSchema([]Entry{}, s)
	if len(violations) != 1 || violations[0].Key != "DB_URL" {
		t.Errorf("expected missing key violation, got %v", violations)
	}
}

func TestCheckSchema_PatternMismatch(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "PORT", Required: true, Pattern: `^\d+$`}}}
	entries := []Entry{{Key: "PORT", Value: "not-a-port"}}
	violations := CheckSchema(entries, s)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheckSchema_PatternMatch(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "PORT", Required: true, Pattern: `^\d+$`}}}
	entries := []Entry{{Key: "PORT", Value: "8080"}}
	violations := CheckSchema(entries, s)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCheckSchema_OptionalMissingNoViolation(t *testing.T) {
	s := &Schema{Fields: []SchemaField{{Key: "LOG_LEVEL", Required: false}}}
	violations := CheckSchema([]Entry{}, s)
	if len(violations) != 0 {
		t.Errorf("expected no violations for optional missing key, got %v", violations)
	}
}
