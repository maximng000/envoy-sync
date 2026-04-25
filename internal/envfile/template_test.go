package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenderTemplate_BasicSubstitution(t *testing.T) {
	entries := map[string]string{"APP_NAME": "envoy", "PORT": "8080"}
	res, err := RenderTemplate("App: {{APP_NAME}} on port {{PORT}}", entries, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "App: envoy on port 8080" {
		t.Errorf("unexpected rendered output: %q", res.Rendered)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
}

func TestRenderTemplate_MissingKeySilent(t *testing.T) {
	entries := map[string]string{"APP_NAME": "envoy"}
	res, err := RenderTemplate("App: {{APP_NAME}} db: {{DB_HOST}}", entries, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_HOST" {
		t.Errorf("expected missing DB_HOST, got %v", res.Missing)
	}
	if res.Rendered != "App: envoy db: " {
		t.Errorf("unexpected rendered output: %q", res.Rendered)
	}
}

func TestRenderTemplate_MissingKeyError(t *testing.T) {
	entries := map[string]string{"APP_NAME": "envoy"}
	_, err := RenderTemplate("App: {{APP_NAME}} db: {{DB_HOST}}", entries, true)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRenderTemplate_NoPlaceholders(t *testing.T) {
	entries := map[string]string{"APP_NAME": "envoy"}
	res, err := RenderTemplate("plain text", entries, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "plain text" {
		t.Errorf("unexpected rendered output: %q", res.Rendered)
	}
}

func TestRenderTemplateFile_ReadsAndRenders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "template.txt")
	content := "HOST={{DB_HOST}} PORT={{DB_PORT}}"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}
	entries := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	res, err := RenderTemplateFile(path, entries, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "HOST=localhost PORT=5432" {
		t.Errorf("unexpected rendered output: %q", res.Rendered)
	}
}

func TestRenderTemplateFile_MissingFile(t *testing.T) {
	_, err := RenderTemplateFile("/nonexistent/template.txt", map[string]string{}, false)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
