package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempScopeEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runScopeCmd(t *testing.T, args ...string) string {
	t.Helper()
	out, err := captureOutput(func() error {
		scopeCmd.ResetFlags()
		scopeCmd.Flags().StringVar(&scopeName, "scope", "", "")
		scopeCmd.Flags().BoolVar(&scopeList, "list", false, "")
		scopeCmd.Flags().BoolVar(&scopeSummary, "summary", false, "")
		return scopeCmd.RunE(scopeCmd, args[len(args)-1:])
	})
	if err != nil {
		t.Logf("scope cmd error: %v", err)
	}
	_ = out
	// Re-run properly via Execute path for simplicity
	return out
}

func TestScopeCmd_ListScopes(t *testing.T) {
	f := writeTempScopeEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAWS_KEY=abc\nAPP_ENV=prod\n")
	out, _ := captureOutput(func() error {
		scopeName = ""
		scopeList = true
		scopeSummary = false
		return scopeCmd.RunE(scopeCmd, []string{f})
	})
	if !strings.Contains(out, "DB") {
		t.Errorf("expected DB in output, got: %s", out)
	}
	if !strings.Contains(out, "AWS") {
		t.Errorf("expected AWS in output, got: %s", out)
	}
}

func TestScopeCmd_FilterByScope(t *testing.T) {
	f := writeTempScopeEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAWS_KEY=abc\n")
	out, _ := captureOutput(func() error {
		scopeName = "DB"
		scopeList = false
		scopeSummary = false
		return scopeCmd.RunE(scopeCmd, []string{f})
	})
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST, got: %s", out)
	}
	if strings.Contains(out, "AWS_KEY") {
		t.Errorf("did not expect AWS_KEY, got: %s", out)
	}
}

func TestScopeCmd_SummaryOutput(t *testing.T) {
	f := writeTempScopeEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAWS_KEY=abc\n")
	out, _ := captureOutput(func() error {
		scopeName = ""
		scopeList = false
		scopeSummary = true
		return scopeCmd.RunE(scopeCmd, []string{f})
	})
	if !strings.Contains(out, "DB") {
		t.Errorf("expected DB in summary, got: %s", out)
	}
}

func TestScopeCmd_MissingFile(t *testing.T) {
	_, err := captureOutput(func() error {
		scopeName = "DB"
		scopeList = false
		scopeSummary = false
		return scopeCmd.RunE(scopeCmd, []string{"/nonexistent/.env"})
	})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
