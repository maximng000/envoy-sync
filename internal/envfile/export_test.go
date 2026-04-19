package envfile

import (
	"strings"
	"testing"
)

var sampleEntries = map[string]string{
	"APP_NAME":    "myapp",
	"DB_PASSWORD": "s3cr3t",
	"PORT":        "8080",
}

func TestExport_DotenvFormat(t *testing.T) {
	var sb strings.Builder
	if err := Export(&sb, sampleEntries, FormatDotenv, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME=myapp in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD=s3cr3t") {
		t.Errorf("expected unmasked password in output")
	}
}

func TestExport_DotenvMasked(t *testing.T) {
	var sb strings.Builder
	if err := Export(&sb, sampleEntries, FormatDotenv, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "DB_PASSWORD=***") {
		t.Errorf("expected masked password, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME unmasked")
	}
}

func TestExport_ExportFormat(t *testing.T) {
	var sb strings.Builder
	if err := Export(&sb, sampleEntries, FormatExport, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "export PORT=8080") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	var sb strings.Builder
	if err := Export(&sb, sampleEntries, FormatJSON, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, `"APP_NAME"`) {
		t.Errorf("expected APP_NAME key in JSON")
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var sb strings.Builder
	err := Export(&sb, sampleEntries, Format("xml"), false)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExport_QuoteValueWithSpace(t *testing.T) {
	entries := map[string]string{"GREETING": "hello world"}
	var sb strings.Builder
	_ = Export(&sb, entries, FormatDotenv, false)
	out := sb.String()
	if !strings.Contains(out, `GREETING="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}
