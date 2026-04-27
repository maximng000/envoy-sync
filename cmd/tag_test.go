package cmd

import (
	"os"
	"strings"
	"testing"
)

func writeTempTagEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "tag-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runTagCmd(args ...string) (string, error) {
	return captureOutput(func() error {
		tagCmd.ResetFlags()
		tagCmd.Flags().String("tag", "", "")
		tagCmd.Flags().String("keys", "", "")
		tagCmd.Flags().Bool("summary", false, "")
		return tagCmd.RunE(tagCmd, args)
	})
}

func TestTagCmd_AllEntries(t *testing.T) {
	f := writeTempTagEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	out, err := runTagCmd(f)
	if err != nil {
		// missing --tag flag should error
		if !strings.Contains(err.Error(), "--tag is required") {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	_ = out
}

func TestTagCmd_WithTagFlag(t *testing.T) {
	f := writeTempTagEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	out, err := captureOutput(func() error {
		tagCmd.ResetFlags()
		tagCmd.Flags().String("tag", "", "")
		tagCmd.Flags().String("keys", "", "")
		tagCmd.Flags().Bool("summary", false, "")
		_ = tagCmd.Flags().Set("tag", "v2")
		return tagCmd.RunE(tagCmd, []string{f})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[v2]") {
		t.Errorf("expected '[v2]' in output, got: %s", out)
	}
}

func TestTagCmd_MissingFile(t *testing.T) {
	_, err := captureOutput(func() error {
		tagCmd.ResetFlags()
		tagCmd.Flags().String("tag", "", "")
		tagCmd.Flags().String("keys", "", "")
		tagCmd.Flags().Bool("summary", false, "")
		_ = tagCmd.Flags().Set("tag", "v1")
		return tagCmd.RunE(tagCmd, []string{"/nonexistent/file.env"})
	})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestTagCmd_SummaryFlag(t *testing.T) {
	f := writeTempTagEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	out, err := captureOutput(func() error {
		tagCmd.ResetFlags()
		tagCmd.Flags().String("tag", "", "")
		tagCmd.Flags().String("keys", "", "")
		tagCmd.Flags().Bool("summary", false, "")
		_ = tagCmd.Flags().Set("tag", "infra")
		_ = tagCmd.Flags().Set("summary", "true")
		return tagCmd.RunE(tagCmd, []string{f})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[infra]") {
		t.Errorf("expected '[infra]' in summary output, got: %s", out)
	}
}
