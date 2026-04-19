package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempSyncEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(p, []byte(content), 0644))
	return p
}

func runSyncCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	return captureOutput(func() error {
		rootCmd.SetArgs(append([]string{"sync"}, args...))
		return rootCmd.Execute()
	})
}

func TestSyncCmd_NewKey(t *testing.T) {
	base := writeTempSyncEnv(t, "A=1\n")
	src := writeTempSyncEnv(t, "B=2\n")
	out, err := runSyncCmd(t, base, src)
	require.NoError(t, err)
	assert.Contains(t, out, "applied")
	assert.Contains(t, out, "B")
}

func TestSyncCmd_ConflictSkipDefault(t *testing.T) {
	base := writeTempSyncEnv(t, "A=old\n")
	src := writeTempSyncEnv(t, "A=new\n")
	out, err := runSyncCmd(t, base, src)
	require.NoError(t, err)
	assert.True(t, strings.Contains(out, "skipped") || strings.Contains(out, "conflict"))
}

func TestSyncCmd_MissingFile(t *testing.T) {
	_, err := runSyncCmd(t, "/no/such/base.env", "/no/such/src.env")
	assert.Error(t, err)
}
