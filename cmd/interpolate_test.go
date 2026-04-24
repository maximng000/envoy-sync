package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempInterpolateEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return path
}

func runInterpolateCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	return captureOutput(func() error {
		rootCmd.SetArgs(append([]string{"interpolate"}, args...))
		return rootCmd.Execute()
	})
}

func TestInterpolateCmd_BasicResolution(t *testing.T) {
	path := writeTempInterpolateEnv(t, "HOST=localhost\nURL=http://${HOST}:9000\n")

	out, err := runInterpolateCmd(t, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "URL=http://localhost:9000") {
		t.Errorf("expected resolved URL, got:\n%s", out)
	}
}

func TestInterpolateCmd_MissingFile(t *testing.T) {
	_, err := runInterpolateCmd(t, "/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestInterpolateCmd_FailOnMissing(t *testing.T) {
	path := writeTempInterpolateEnv(t, "URL=http://${UNDEFINED_HOST}:9000\n")

	_, err := runInterpolateCmd(t, "--fail-on-missing", path)
	if err == nil {
		t.Fatal("expected error when --fail-on-missing is set")
	}
}

func TestInterpolateCmd_MaskSecrets(t *testing.T) {
	path := writeTempInterpolateEnv(t, "BASE=http://api\nAPI_SECRET=topsecret\nURL=${BASE}/endpoint\n")

	out, err := runInterpolateCmd(t, "--mask", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "topsecret") {
		t.Errorf("expected secret to be masked, got:\n%s", out)
	}
	if !strings.Contains(out, "****") {
		t.Errorf("expected masked placeholder in output, got:\n%s", out)
	}
}
