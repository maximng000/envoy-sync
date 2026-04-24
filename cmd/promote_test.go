package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempPromoteEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runPromoteCmd(t *testing.T, args ...string) string {
	t.Helper()
	out, err := captureOutput(func() error {
		promoteCmd.ResetFlags()
		promoteCmd.Flags().Bool("overwrite", false, "")
		promoteCmd.Flags().Bool("dry-run", false, "")
		promoteCmd.Flags().Bool("mask-secrets", true, "")
		promoteCmd.SetArgs(args)
		return promoteCmd.Execute()
	})
	if err != nil {
		t.Logf("promoteCmd output: %s", out)
	}
	return out
}

func TestPromoteCmd_AddedKey(t *testing.T) {
	dir := t.TempDir()
	src := writeTempPromoteEnv(t, dir, "src.env", "NEW_KEY=hello\n")
	dst := writeTempPromoteEnv(t, dir, "dst.env", "EXISTING=world\n")
	out := runPromoteCmd(t, src, dst)
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected added key in output, got: %s", out)
	}
}

func TestPromoteCmd_SkipConflict(t *testing.T) {
	dir := t.TempDir()
	src := writeTempPromoteEnv(t, dir, "src.env", "KEY=new\n")
	dst := writeTempPromoteEnv(t, dir, "dst.env", "KEY=old\n")
	out := runPromoteCmd(t, src, dst)
	if !strings.Contains(out, "! KEY skipped") {
		t.Errorf("expected skip message, got: %s", out)
	}
}

func TestPromoteCmd_DryRun(t *testing.T) {
	dir := t.TempDir()
	src := writeTempPromoteEnv(t, dir, "src.env", "NEW=val\n")
	dst := writeTempPromoteEnv(t, dir, "dst.env", "")
	out := runPromoteCmd(t, "--dry-run", src, dst)
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run notice, got: %s", out)
	}
}

func TestPromoteCmd_MissingFile(t *testing.T) {
	_, err := captureOutput(func() error {
		promoteCmd.ResetFlags()
		promoteCmd.Flags().Bool("overwrite", false, "")
		promoteCmd.Flags().Bool("dry-run", false, "")
		promoteCmd.Flags().Bool("mask-secrets", true, "")
		promoteCmd.SetArgs([]string{"/nonexistent/src.env", "/nonexistent/dst.env"})
		return promoteCmd.Execute()
	})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
