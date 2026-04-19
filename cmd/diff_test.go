package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempDiffEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runDiffCmd(t *testing.T, args ...string) string {
	t.Helper()
	out, err := captureOutput(func() error {
		rootCmd.SetArgs(append([]string{"diff"}, args...))
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("diff command failed: %v", err)
	}
	return out
}

func TestDiffCmd_AddedKey(t *testing.T) {
	base := writeTempDiffEnv(t, "FOO=bar\n")
	other := writeTempDiffEnv(t, "FOO=bar\nBAZ=qux\n")
	out := runDiffCmd(t, base, other)
	if !strings.Contains(out, "+ BAZ") {
		t.Errorf("expected added key in output, got: %s", out)
	}
}

func TestDiffCmd_RemovedKey(t *testing.T) {
	base := writeTempDiffEnv(t, "FOO=bar\nBAZ=qux\n")
	other := writeTempDiffEnv(t, "FOO=bar\n")
	out := runDiffCmd(t, base, other)
	if !strings.Contains(out, "- BAZ") {
		t.Errorf("expected removed key in output, got: %s", out)
	}
}

func TestDiffCmd_NoDifferences(t *testing.T) {
	base := writeTempDiffEnv(t, "FOO=bar\n")
	other := writeTempDiffEnv(t, "FOO=bar\n")
	out := runDiffCmd(t, base, other)
	if !strings.Contains(out, "No differences found.") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestDiffCmd_MissingFile(t *testing.T) {
	_, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"diff", "nonexistent.env", "also_missing.env"})
		return rootCmd.Execute()
	})
	if err == nil {
		t.Error("expected error for missing files, got nil")
	}
}
