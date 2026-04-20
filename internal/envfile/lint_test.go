package envfile

import (
	"testing"
)

func TestLint_CleanFile(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	}
	result := Lint(entries)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "staging"},
		{Key: "APP_ENV", Value: "production"},
	}
	result := Lint(entries)
	if !result.HasErrors() {
		t.Fatal("expected an error for duplicate key")
	}
	if result.Issues[0].Severity != "error" {
		t.Errorf("expected severity 'error', got %s", result.Issues[0].Severity)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	entries := []Entry{
		{Key: "app_env", Value: "production"},
	}
	result := Lint(entries)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warn" {
		t.Errorf("expected severity 'warn', got %s", result.Issues[0].Severity)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: ""},
	}
	result := Lint(entries)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Message != "value is empty" {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_WhitespaceValue(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: " localhost "},
	}
	result := Lint(entries)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Message != "value has leading or trailing whitespace" {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_HasErrors_False(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: ""},
	}
	result := Lint(entries)
	if result.HasErrors() {
		t.Error("expected no errors, only warnings")
	}
}
