package cmd

import (
	"os"
	"strings"
	"testing"
)

func writeTempStatsEnv(t *testing.T, content string) string {
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

func runStatsCmd(t *testing.T, args ...string) string {
	t.Helper()
	out, err := captureOutput(func() error {
		rootCmd.SetArgs(append([]string{"stats"}, args...))
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("stats command failed: %v", err)
	}
	return out
}

func TestStatsCmd_BasicOutput(t *testing.T) {
	path := writeTempStatsEnv(t, "APP_NAME=myapp\nDB_PASSWORD=secret\nAPP_ENV=prod\n")
	out := runStatsCmd(t, path)

	if !strings.Contains(out, "Total:") {
		t.Error("expected Total in output")
	}
	if !strings.Contains(out, "Secrets:") {
		t.Error("expected Secrets in output")
	}
}

func TestStatsCmd_ShowsTopPrefixes(t *testing.T) {
	path := writeTempStatsEnv(t, "APP_NAME=myapp\nAPP_ENV=prod\nDB_HOST=localhost\n")
	out := runStatsCmd(t, "--top-prefixes", "2", path)

	if !strings.Contains(out, "Top prefixes:") {
		t.Error("expected Top prefixes in output")
	}
	if !strings.Contains(out, "APP") {
		t.Error("expected APP prefix in output")
	}
}

func TestStatsCmd_MissingFile(t *testing.T) {
	_, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"stats", "/nonexistent/.env"})
		return rootCmd.Execute()
	})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
