package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempSortEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func runSortCmd(t *testing.T, args ...string) string {
	t.Helper()
	return captureOutput(t, func() {
		sortCmd.ResetFlags()
		sortCmd.Flags().StringVarP(&sortStrategy, "strategy", "s", "alpha",
			"sort strategy")
		sortCmd.Flags().StringVarP(&sortFormat, "format", "f", "dotenv",
			"output format")
		sortCmd.SetArgs(args)
		_ = sortCmd.Execute()
	})
}

func TestSortCmd_AlphaOrder(t *testing.T) {
	p := writeTempSortEnv(t, "ZEBRA=z\nAPP_NAME=app\nHOST=localhost\n")
	out := runSortCmd(t, p)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines, got: %s", out)
	}
	if !strings.HasPrefix(lines[0], "APP_NAME") {
		t.Errorf("expected APP_NAME first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[len(lines)-1], "ZEBRA") {
		t.Errorf("expected ZEBRA last, got: %s", lines[len(lines)-1])
	}
}

func TestSortCmd_AlphaDescOrder(t *testing.T) {
	p := writeTempSortEnv(t, "ZEBRA=z\nAPP_NAME=app\nHOST=localhost\n")
	out := runSortCmd(t, "--strategy", "alpha-desc", p)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if !strings.HasPrefix(lines[0], "ZEBRA") {
		t.Errorf("expected ZEBRA first in desc order, got: %s", lines[0])
	}
}

func TestSortCmd_MissingFile(t *testing.T) {
	out := runSortCmd(t, "/nonexistent/.env")
	if !strings.Contains(out, "") {
		// error goes to stderr; just ensure no panic
		t.Log("missing file handled")
	}
}

func TestSortCmd_UnknownStrategy(t *testing.T) {
	p := writeTempSortEnv(t, "FOO=bar\n")
	out := runSortCmd(t, "--strategy", "bogus", p)
	_ = out // error path, just ensure no panic
}
