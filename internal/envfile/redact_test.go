package envfile

import (
	"testing"
)

func TestRedact_AutoDetectSecrets(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
		{Key: "API_KEY", Value: "abc123"},
	}

	result := Redact(entries, RedactModeMask, nil)

	if result.Entries[0].Value != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", result.Entries[0].Value)
	}
	if result.Entries[1].Value != "***" {
		t.Errorf("expected DB_PASSWORD masked, got %q", result.Entries[1].Value)
	}
	if result.Entries[2].Value != "***" {
		t.Errorf("expected API_KEY masked, got %q", result.Entries[2].Value)
	}
	if len(result.Redacted) != 2 {
		t.Errorf("expected 2 redacted keys, got %d", len(result.Redacted))
	}
}

func TestRedact_BlankMode(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET_TOKEN", Value: "topsecret"},
	}
	result := Redact(entries, RedactModeBlank, nil)
	if result.Entries[0].Value != "" {
		t.Errorf("expected blank value, got %q", result.Entries[0].Value)
	}
}

func TestRedact_PlaceholderMode(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "hunter2"},
	}
	result := Redact(entries, RedactModePlaceholder, nil)
	if result.Entries[0].Value != "{{DB_PASSWORD}}" {
		t.Errorf("expected placeholder, got %q", result.Entries[0].Value)
	}
}

func TestRedact_ExplicitKeys(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
	}
	// Neither is a secret by name, but we force redact APP_ENV
	result := Redact(entries, RedactModeMask, []string{"APP_ENV"})
	if result.Entries[0].Value != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", result.Entries[0].Value)
	}
	if result.Entries[1].Value != "***" {
		t.Errorf("expected APP_ENV masked, got %q", result.Entries[1].Value)
	}
	if len(result.Redacted) != 1 || result.Redacted[0] != "APP_ENV" {
		t.Errorf("expected [APP_ENV] in redacted list, got %v", result.Redacted)
	}
}

func TestRedact_OriginalUnmodified(t *testing.T) {
	entries := []Entry{
		{Key: "API_SECRET", Value: "original"},
	}
	Redact(entries, RedactModeMask, nil)
	if entries[0].Value != "original" {
		t.Error("Redact must not modify original entries slice")
	}
}
