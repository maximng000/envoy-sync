package cmd

import (
	"os"
	"strings"
	"testing"
)

func writeTempNormalizeEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func runNormalizeCmd(t *testing.T, args ...string) string {
	t.Helper()
	out, err := captureOutput(func() error {
		normalizeCmd.ResetFlags()
		normalizeCmd.Flags().BoolVar(&normalizeUppercase, "uppercase", false, "")
		normalizeCmd.Flags().BoolVar(&normalizeTrim, "trim", false, "")
		normalizeCmd.Flags().BoolVar(&normalizeRemoveEmpty, "remove-empty", false, "")
		normalizeCmd.Flags().BoolVar(&normalizeSort, "sort", false, "")
		normalizeCmd.Flags().BoolVar(&normalizeWrite, "write", false, "")
		normalizeCmd.SetArgs(args)
		return normalizeCmd.Execute()
	})
	if err != nil {
		t.Fatalf("command error: %v", err)
	}
	return out
}

func TestNormalizeCmd_ShowsSummary(t *testing.T) {
	path := writeTempNormalizeEnv(t, "app_name=myapp\ndb_host=localhost\n")
	out := runNormalizeCmd(t, path)
	if !strings.Contains(out, "kept") {
		t.Errorf("expected summary in output, got: %s", out)
	}
}

func TestNormalizeCmd_UppercaseFlag(t *testing.T) {
	path := writeTempNormalizeEnv(t, "app_name=myapp\n")
	out := runNormalizeCmd(t, "--uppercase", path)
	if !strings.Contains(out, "modified") {
		t.Errorf("expected modified key in output, got: %s", out)
	}
}

func TestNormalizeCmd_RemoveEmptyFlag(t *testing.T) {
	path := writeTempNormalizeEnv(t, "PRESENT=yes\nEMPTY=\n")
	out := runNormalizeCmd(t, "--remove-empty", path)
	if !strings.Contains(out, "removed") {
		t.Errorf("expected removed key in output, got: %s", out)
	}
}

func TestNormalizeCmd_MissingFile(t *testing.T) {
	_, err := captureOutput(func() error {
		normalizeCmd.SetArgs([]string{"/nonexistent/.env"})
		return normalizeCmd.Execute()
	})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
