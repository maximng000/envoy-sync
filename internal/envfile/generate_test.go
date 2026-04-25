package envfile

import (
	"strings"
	"testing"
)

func TestGenerate_BasicKeys(t *testing.T) {
	result := Generate([]string{"APP_NAME", "PORT"}, GenerateOptions{})
	if result.Count != 2 {
		t.Fatalf("expected count 2, got %d", result.Count)
	}
	if !contains(result.Lines, "APP_NAME=CHANGEME") {
		t.Error("expected APP_NAME=CHANGEME in output")
	}
	if !contains(result.Lines, "PORT=CHANGEME") {
		t.Error("expected PORT=CHANGEME in output")
	}
}

func TestGenerate_SecretKeyMasked(t *testing.T) {
	result := Generate([]string{"DB_PASSWORD", "API_KEY"}, GenerateOptions{})
	for _, line := range result.Lines {
		if strings.Contains(line, "CHANGEME") {
			t.Errorf("secret key should not use CHANGEME placeholder, got: %s", line)
		}
	}
	if !contains(result.Lines, "DB_PASSWORD=***") {
		t.Error("expected DB_PASSWORD=*** in output")
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	result := Generate([]string{"HOST"}, GenerateOptions{Placeholder: "TODO"})
	if !contains(result.Lines, "HOST=TODO") {
		t.Errorf("expected HOST=TODO, got: %v", result.Lines)
	}
}

func TestGenerate_WithComments(t *testing.T) {
	result := Generate([]string{"API_SECRET", "APP_ENV"}, GenerateOptions{IncludeComments: true})
	foundSecretComment := false
	for _, line := range result.Lines {
		if strings.Contains(line, "API_SECRET") && strings.Contains(line, "secret value") {
			foundSecretComment = true
		}
	}
	if !foundSecretComment {
		t.Error("expected secret comment for API_SECRET")
	}
}

func TestGenerate_SkipsEmptyKeys(t *testing.T) {
	result := Generate([]string{"APP_NAME", "", "  "}, GenerateOptions{})
	for _, line := range result.Lines {
		if strings.TrimSpace(line) == "=CHANGEME" || line == "=CHANGEME" {
			t.Error("empty key should be skipped")
		}
	}
}

func TestGenerateFromEntries_UsesKeys(t *testing.T) {
	entries := map[string]string{
		"APP_ENV":    "production",
		"DB_PASSWORD": "supersecret",
	}
	result := GenerateFromEntries(entries, GenerateOptions{})
	if result.Count != 2 {
		t.Fatalf("expected count 2, got %d", result.Count)
	}
	if !contains(result.Lines, "APP_ENV=CHANGEME") {
		t.Error("expected APP_ENV=CHANGEME")
	}
	if !contains(result.Lines, "DB_PASSWORD=***") {
		t.Error("expected DB_PASSWORD=***")
	}
}

func contains(lines []string, target string) bool {
	for _, l := range lines {
		if l == target {
			return true
		}
	}
	return false
}
