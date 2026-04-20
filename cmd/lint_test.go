package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLintEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func runLintCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	return captureOutput(t, func() error {
		rootCmd.SetArgs(append([]string{"lint"}, args...))
		return rootCmd.Execute()
	})
}

func TestLintCmd_CleanFile(t *testing.T) {
	path := writeTempLintEnv(t, "APP_ENV=production\nPORT=8080\n")
	out, err := runLintCmd(t, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No issues found") {
		t.Errorf("expected clean output, got: %s", out)
	}
}

func TestLintCmd_DuplicateKey(t *testing.T) {
	path := writeTempLintEnv(t, "APP_ENV=staging\nAPP_ENV=production\n")
	out, err := runLintCmd(t, path)
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
	if !strings.Contains(out, "duplicate key") {
		t.Errorf("expected duplicate key message, got: %s", out)
	}
}

func TestLintCmd_LowercaseKeyWarn(t *testing.T) {
	path := writeTempLintEnv(t, "app_env=production\n")
	out, err := runLintCmd(t, path)
	if err != nil {
		t.Fatalf("unexpected error (warnings should not fail): %v", err)
	}
	if !strings.Contains(out, "UPPER_SNAKE_CASE") {
		t.Errorf("expected uppercase warning, got: %s", out)
	}
}

func TestLintCmd_MissingFile(t *testing.T) {
	_, err := runLintCmd(t, "/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
