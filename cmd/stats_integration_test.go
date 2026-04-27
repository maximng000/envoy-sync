package cmd

import (
	"strings"
	"testing"
)

func TestStatsIntegration_DuplicateKeys(t *testing.T) {
	path := writeTempStatsEnv(t,
		"APP_NAME=first\nAPP_NAME=second\nDB_HOST=localhost\n")

	out := runStatsCmd(t, path)

	if !strings.Contains(out, "Duplicates:") {
		t.Error("expected Duplicates field in output")
	}
	if !strings.Contains(out, "Total:") {
		t.Error("expected Total field in output")
	}
}

func TestStatsIntegration_EmptyValues(t *testing.T) {
	path := writeTempStatsEnv(t,
		"APP_NAME=\nDB_HOST=localhost\nAPI_KEY=\n")

	out := runStatsCmd(t, path)

	if !strings.Contains(out, "Empty:") {
		t.Error("expected Empty field in output")
	}
}

func TestStatsIntegration_SecretsCount(t *testing.T) {
	path := writeTempStatsEnv(t,
		"APP_NAME=myapp\nDB_PASSWORD=secret\nAPI_SECRET=topsecret\nAPP_ENV=prod\n")

	out := runStatsCmd(t, path)

	if !strings.Contains(out, "Secrets:") {
		t.Error("expected Secrets field in output")
	}
	if !strings.Contains(out, "Non-secrets:") {
		t.Error("expected Non-secrets field in output")
	}
}
