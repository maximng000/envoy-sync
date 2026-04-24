package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPromoteIntegration_OverwriteUpdatesFile(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "src.env")
	dstPath := filepath.Join(dir, "dst.env")

	if err := os.WriteFile(srcPath, []byte("KEY=updated\nNEW=added\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dstPath, []byte("KEY=original\n"), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := captureOutput(func() error {
		promoteCmd.ResetFlags()
		promoteCmd.Flags().Bool("overwrite", false, "")
		promoteCmd.Flags().Bool("dry-run", false, "")
		promoteCmd.Flags().Bool("mask-secrets", true, "")
		promoteCmd.SetArgs([]string{"--overwrite", srcPath, dstPath})
		return promoteCmd.Execute()
	})
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}

	if !strings.Contains(out, "~ KEY") {
		t.Errorf("expected update marker for KEY, got: %s", out)
	}
	if !strings.Contains(out, "+ NEW") {
		t.Errorf("expected add marker for NEW, got: %s", out)
	}

	data, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "KEY=updated") {
		t.Errorf("expected updated value in dst, got:\n%s", content)
	}
	if !strings.Contains(content, "NEW=added") {
		t.Errorf("expected new key in dst, got:\n%s", content)
	}
}
