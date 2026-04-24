package cmd

import (
	"os"
	"strings"
	"testing"
)

func writeTempCloneEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func runCloneCmd(t *testing.T, args ...string) string {
	t.Helper()
	rootCmd.SetArgs(append([]string{"clone"}, args...))
	out, err := captureOutput(func() error {
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("clone command failed: %v", err)
	}
	return out
}

func TestCloneCmd_NewKey(t *testing.T) {
	src := writeTempCloneEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	dst := writeTempCloneEnv(t, "")
	out := runCloneCmd(t, src, dst)
	if !strings.Contains(out, "CLONE APP_NAME") {
		t.Errorf("expected CLONE APP_NAME in output, got: %s", out)
	}
}

func TestCloneCmd_SkipConflictDefault(t *testing.T) {
	src := writeTempCloneEnv(t, "PORT=9090\n")
	dst := writeTempCloneEnv(t, "PORT=8080\n")
	out := runCloneCmd(t, src, dst)
	if !strings.Contains(out, "SKIP  PORT") {
		t.Errorf("expected SKIP PORT in output, got: %s", out)
	}
}

func TestCloneCmd_OverwriteFlag(t *testing.T) {
	src := writeTempCloneEnv(t, "PORT=9090\n")
	dst := writeTempCloneEnv(t, "PORT=8080\n")
	out := runCloneCmd(t, "--overwrite", src, dst)
	if !strings.Contains(out, "CLONE PORT") {
		t.Errorf("expected CLONE PORT in output, got: %s", out)
	}
}

func TestCloneCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"clone", "/nonexistent/src.env", "/tmp/dst.env"})
	_, err := captureOutput(func() error {
		return rootCmd.Execute()
	})
	if err == nil {
		t.Error("expected error for missing src file")
	}
}
