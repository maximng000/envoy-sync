package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempInjectEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempInjectEnv: %v", err)
	}
	return p
}

func runInjectCmd(t *testing.T, args ...string) string {
	t.Helper()
	out, err := captureOutput(func() error {
		rootCmd.SetArgs(append([]string{"inject"}, args...))
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("inject cmd error: %v", err)
	}
	return out
}

func TestInjectCmd_InjectsKey(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("CMD_INJECT_KEY") })
	f := writeTempInjectEnv(t, "CMD_INJECT_KEY=hello\n")
	out := runInjectCmd(t, f)
	if !strings.Contains(out, "CMD_INJECT_KEY") {
		t.Errorf("expected CMD_INJECT_KEY in output, got: %s", out)
	}
	if os.Getenv("CMD_INJECT_KEY") != "hello" {
		t.Errorf("expected CMD_INJECT_KEY=hello in env")
	}
}

func TestInjectCmd_MissingFile(t *testing.T) {
	_, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"inject", "/nonexistent/.env"})
		return rootCmd.Execute()
	})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestInjectCmd_SkipExisting(t *testing.T) {
	os.Setenv("CMD_SKIP_KEY", "original")
	t.Cleanup(func() { os.Unsetenv("CMD_SKIP_KEY") })
	f := writeTempInjectEnv(t, "CMD_SKIP_KEY=new\n")
	out := runInjectCmd(t, f)
	if !strings.Contains(out, "Skipped") {
		t.Errorf("expected Skipped in output, got: %s", out)
	}
	if os.Getenv("CMD_SKIP_KEY") != "original" {
		t.Errorf("expected original value preserved")
	}
}
