package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeTempConvertEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func runConvertCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "envoy-sync"}
	convertCmd := &cobra.Command{}
	*convertCmd = *convertCmdRef
	root.AddCommand(convertCmd)
	return captureOutput(func() error {
		root.SetArgs(args)
		return root.Execute()
	})
}

// convertCmdRef is set in convert.go's init via a package-level var so tests can clone it.
// We reference the registered subcommand directly here instead.
func runConvertCmdDirect(t *testing.T, args ...string) string {
	t.Helper()
	out, err := captureOutput(func() error {
		RootCmd.SetArgs(append([]string{"convert"}, args...))
		return RootCmd.Execute()
	})
	if err != nil {
		t.Logf("convert cmd error (may be expected): %v", err)
	}
	return out
}

func TestConvertCmd_DotenvOutput(t *testing.T) {
	path := writeTempConvertEnv(t, "APP_ENV=production\nDB_PASSWORD=secret123\n")
	out := runConvertCmdDirect(t, "--file", path, "--from", "dotenv", "--to", "shell")
	if !strings.Contains(out, "export APP_ENV") {
		t.Errorf("expected shell export syntax, got: %s", out)
	}
}

func TestConvertCmd_MissingFile(t *testing.T) {
	out := runConvertCmdDirect(t, "--file", "/nonexistent/.env", "--from", "dotenv", "--to", "shell")
	if !strings.Contains(strings.ToLower(out), "error") &&
		!strings.Contains(strings.ToLower(out), "no such file") &&
		!strings.Contains(strings.ToLower(out), "failed") {
		t.Errorf("expected error message for missing file, got: %s", out)
	}
}

func TestConvertCmd_UnsupportedFormat(t *testing.T) {
	path := writeTempConvertEnv(t, "KEY=value\n")
	out := runConvertCmdDirect(t, "--file", path, "--from", "dotenv", "--to", "yaml")
	if !strings.Contains(strings.ToLower(out), "unsupported") &&
		!strings.Contains(strings.ToLower(out), "error") {
		t.Errorf("expected unsupported format error, got: %s", out)
	}
}

func TestConvertCmd_SecretMasked(t *testing.T) {
	path := writeTempConvertEnv(t, "API_SECRET=topsecret\nAPP_NAME=myapp\n")
	out := runConvertCmdDirect(t, "--file", path, "--from", "dotenv", "--to", "dotenv", "--mask-secrets")
	if strings.Contains(out, "topsecret") {
		t.Errorf("expected secret to be masked, but raw value appeared in output: %s", out)
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected non-secret key to appear in output: %s", out)
	}
}
