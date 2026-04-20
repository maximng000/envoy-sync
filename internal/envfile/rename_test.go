package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseEntries() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_PASS":  "secret",
	}
}

func TestRenameEntry_Success(t *testing.T) {
	updated, res := RenameEntry(baseEntries(), "APP_HOST", "SERVICE_HOST", false)
	assert.True(t, res.Renamed)
	assert.Equal(t, "ok", res.Reason)
	_, oldExists := updated["APP_HOST"]
	assert.False(t, oldExists)
	assert.Equal(t, "localhost", updated["SERVICE_HOST"])
}

func TestRenameEntry_OldKeyMissing(t *testing.T) {
	updated, res := RenameEntry(baseEntries(), "MISSING_KEY", "NEW_KEY", false)
	assert.False(t, res.Renamed)
	assert.Contains(t, res.Reason, "not found")
	_, newExists := updated["NEW_KEY"]
	assert.False(t, newExists)
}

func TestRenameEntry_ConflictNoOverwrite(t *testing.T) {
	updated, res := RenameEntry(baseEntries(), "APP_HOST", "APP_PORT", false)
	assert.False(t, res.Renamed)
	assert.Contains(t, res.Reason, "already exists")
	// original values preserved
	assert.Equal(t, "localhost", updated["APP_HOST"])
	assert.Equal(t, "8080", updated["APP_PORT"])
}

func TestRenameEntry_ConflictWithOverwrite(t *testing.T) {
	updated, res := RenameEntry(baseEntries(), "APP_HOST", "APP_PORT", true)
	assert.True(t, res.Renamed)
	_, oldExists := updated["APP_HOST"]
	assert.False(t, oldExists)
	assert.Equal(t, "localhost", updated["APP_PORT"])
}

func TestRenameEntry_SameKey(t *testing.T) {
	updated, res := RenameEntry(baseEntries(), "APP_HOST", "APP_HOST", false)
	assert.False(t, res.Renamed)
	assert.Contains(t, res.Reason, "identical")
	assert.Equal(t, "localhost", updated["APP_HOST"])
}

func TestRenameEntry_OriginalUnmodified(t *testing.T) {
	orig := baseEntries()
	_, _ = RenameEntry(orig, "APP_HOST", "SERVICE_HOST", false)
	// original must not be mutated
	assert.Equal(t, "localhost", orig["APP_HOST"])
}

func TestBulkRename_AllSucceed(t *testing.T) {
	renames := [][2]string{
		{"APP_HOST", "SERVICE_HOST"},
		{"APP_PORT", "SERVICE_PORT"},
	}
	updated, results := BulkRename(baseEntries(), renames, false)
	assert.Len(t, results, 2)
	assert.True(t, results[0].Renamed)
	assert.True(t, results[1].Renamed)
	assert.Equal(t, "localhost", updated["SERVICE_HOST"])
	assert.Equal(t, "8080", updated["SERVICE_PORT"])
}

func TestBulkRename_PartialFailure(t *testing.T) {
	renames := [][2]string{
		{"APP_HOST", "SERVICE_HOST"},
		{"NONEXISTENT", "OTHER"},
	}
	updated, results := BulkRename(baseEntries(), renames, false)
	assert.True(t, results[0].Renamed)
	assert.False(t, results[1].Renamed)
	assert.Equal(t, "localhost", updated["SERVICE_HOST"])
	_, exists := updated["OTHER"]
	assert.False(t, exists)
}
