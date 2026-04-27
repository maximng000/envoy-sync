package cmd

import (
	"strings"
	"testing"
)

func TestTagIntegration_SpecificKeys(t *testing.T) {
	f := writeTempTagEnv(t, "APP_NAME=myapp\nDB_PASSWORD=secret\nPORT=8080\n")

	out, err := captureOutput(func() error {
		tagCmd.ResetFlags()
		tagCmd.Flags().String("tag", "", "")
		tagCmd.Flags().String("keys", "", "")
		tagCmd.Flags().Bool("summary", false, "")
		_ = tagCmd.Flags().Set("tag", "monitored")
		_ = tagCmd.Flags().Set("keys", "APP_NAME,PORT")
		return tagCmd.RunE(tagCmd, []string{f})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME in output, got: %s", out)
	}
	if !strings.Contains(out, "[monitored]") {
		t.Errorf("expected tag label in output, got: %s", out)
	}
	if !strings.Contains(out, "skipped: DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in skipped, got: %s", out)
	}
}

func TestTagIntegration_SummaryGroups(t *testing.T) {
	f := writeTempTagEnv(t, "HOST=localhost\nPORT=5432\n")

	out, err := captureOutput(func() error {
		tagCmd.ResetFlags()
		tagCmd.Flags().String("tag", "", "")
		tagCmd.Flags().String("keys", "", "")
		tagCmd.Flags().Bool("summary", false, "")
		_ = tagCmd.Flags().Set("tag", "db")
		_ = tagCmd.Flags().Set("summary", "true")
		return tagCmd.RunE(tagCmd, []string{f})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[db]") {
		t.Errorf("expected '[db]' group header, got: %s", out)
	}
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("expected HOST entry in summary, got: %s", out)
	}
}
