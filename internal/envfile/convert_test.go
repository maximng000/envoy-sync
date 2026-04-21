package envfile

import (
	"strings"
	"testing"
)

func TestConvert_DockerCompose(t *testing.T) {
	entries := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	res, err := Convert(entries, FormatDockerCompose, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(res.Content, "environment:\n") {
		t.Errorf("expected docker-compose header, got: %s", res.Content)
	}
	if !strings.Contains(res.Content, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME entry in output")
	}
}

func TestConvert_Shell(t *testing.T) {
	entries := map[string]string{
		"DB_HOST": "localhost",
	}
	res, err := Convert(entries, FormatShell, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Content, "export DB_HOST=localhost") {
		t.Errorf("expected export statement, got: %s", res.Content)
	}
}

func TestConvert_Dotenv(t *testing.T) {
	entries := map[string]string{
		"API_URL": "https://example.com",
	}
	res, err := Convert(entries, FormatDotenv, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Content, "API_URL=") {
		t.Errorf("expected API_URL in dotenv output, got: %s", res.Content)
	}
}

func TestConvert_SecretMasked(t *testing.T) {
	entries := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "myapp",
	}
	res, err := Convert(entries, FormatDockerCompose, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(res.Content, "supersecret") {
		t.Errorf("expected secret to be masked, but found plain value")
	}
	if !strings.Contains(res.Content, "APP_NAME=myapp") {
		t.Errorf("expected non-secret value to remain unmasked")
	}
}

func TestConvert_UnsupportedFormat(t *testing.T) {
	entries := map[string]string{"KEY": "val"}
	_, err := Convert(entries, ConvertFormat("xml"), false)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestConvert_ShellQuotesSpaces(t *testing.T) {
	entries := map[string]string{
		"APP_DESC": "hello world",
	}
	res, err := Convert(entries, FormatShell, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Content, `"hello world"`) {
		t.Errorf("expected quoted value for space-containing string, got: %s", res.Content)
	}
}
