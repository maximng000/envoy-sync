package cmd

import (
	"os"
	"strings"
	"testing"

	"envoy-sync/internal/envfile"
)

func TestCloneIntegration_WritesToDstFile(t *testing.T) {
	src := writeTempCloneEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	dst := writeTempCloneEnv(t, "EXISTING=yes\n")

	runCloneCmd(t, src, dst)

	result, err := envfile.Parse(dst)
	if err != nil {
		t.Fatalf("failed to parse dst after clone: %v", err)
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp in dst, got %q", result["APP_NAME"])
	}
	if result["EXISTING"] != "yes" {
		t.Errorf("expected EXISTING key preserved, got %q", result["EXISTING"])
	}
}

func TestCloneIntegration_FilterKeys(t *testing.T) {
	src := writeTempCloneEnv(t, "APP_NAME=myapp\nPORT=8080\nDEBUG=true\n")
	dst := writeTempCloneEnv(t, "")

	runCloneCmd(t, "--keys", "PORT", src, dst)

	result, err := envfile.Parse(dst)
	if err != nil {
		t.Fatalf("failed to parse dst: %v", err)
	}
	if _, ok := result["APP_NAME"]; ok {
		t.Error("APP_NAME should not have been cloned")
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", result["PORT"])
	}
}

func TestCloneIntegration_CreatesNewDst(t *testing.T) {
	src := writeTempCloneEnv(t, "KEY=value\n")
	newDst := src + ".new.env"
	t.Cleanup(func() { os.Remove(newDst) })

	rootCmd.SetArgs([]string{"clone", src, newDst})
	_, err := captureOutput(func() error { return rootCmd.Execute() })
	if err != nil {
		t.Fatalf("clone to new file failed: %v", err)
	}

	b, err := os.ReadFile(newDst)
	if err != nil {
		t.Fatalf("new dst file not created: %v", err)
	}
	if !strings.Contains(string(b), "KEY=value") {
		t.Errorf("expected KEY=value in new dst, got: %s", string(b))
	}
}
