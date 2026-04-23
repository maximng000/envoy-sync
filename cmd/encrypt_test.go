package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-sync/internal/envfile"
)

const encTestKey = "0123456789abcdef"

func writeTempEncryptEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEncryptEnv: %v", err)
	}
	return p
}

func TestEncryptCmd_EncryptsSecrets(t *testing.T) {
	src := writeTempEncryptEnv(t, "APP_NAME=myapp\nDB_PASSWORD=hunter2\n")
	out := filepath.Join(t.TempDir(), "encrypted.env")

	output, err := captureOutput(func() error {
		return rootCmd.Execute()
	}, "encrypt", src, "--key", encTestKey, "--output", out)
	_ = output
	_ = err

	rootCmd.ResetFlags()
	encryptCmd.ResetFlags()

	// Re-run via direct args
	rootCmd.SetArgs([]string{"encrypt", src, "--key", encTestKey, "--output", out})
	if execErr := rootCmd.Execute(); execErr != nil {
		t.Fatalf("encrypt command failed: %v", execErr)
	}

	entries, err := envfile.Parse(out)
	if err != nil {
		t.Fatalf("parse encrypted file: %v", err)
	}
	for _, e := range entries {
		if e.Key == "DB_PASSWORD" && e.Value == "hunter2" {
			t.Error("DB_PASSWORD should be encrypted")
		}
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("APP_NAME should be unchanged, got %q", e.Value)
		}
	}
}

func TestEncryptCmd_InvalidKeyLength(t *testing.T) {
	src := writeTempEncryptEnv(t, "DB_PASSWORD=secret\n")
	rootCmd.SetArgs([]string{"encrypt", src, "--key", "short"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for short key")
	}
	if !strings.Contains(err.Error(), "16, 24, or 32") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEncryptCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"encrypt", "/nonexistent/.env", "--key", encTestKey})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing file")
	}
}
