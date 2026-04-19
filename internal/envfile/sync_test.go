package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSync_NewKeysApplied(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}
	out, sr := Sync(dst, src, StrategySkip)
	assert.Equal(t, "1", out["A"])
	assert.Equal(t, "2", out["B"])
	assert.Contains(t, sr.Applied, "B")
}

func TestSync_ConflictSkip(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"A": "new"}
	out, sr := Sync(dst, src, StrategySkip)
	assert.Equal(t, "old", out["A"])
	assert.Contains(t, sr.Skipped, "A")
	assert.Len(t, sr.Conflicts, 1)
}

func TestSync_ConflictOverride(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"A": "new"}
	out, sr := Sync(dst, src, StrategyOverride)
	assert.Equal(t, "new", out["A"])
	assert.Contains(t, sr.Applied, "A")
	assert.Len(t, sr.Conflicts, 1)
}

func TestSync_IdenticalValueSkipped(t *testing.T) {
	dst := map[string]string{"A": "same"}
	src := map[string]string{"A": "same"}
	_, sr := Sync(dst, src, StrategyOverride)
	assert.Contains(t, sr.Skipped, "A")
	assert.Empty(t, sr.Conflicts)
}

func TestSync_DstUnchangedWhenNoSrc(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{}
	out, sr := Sync(dst, src, StrategySkip)
	assert.Equal(t, dst, out)
	assert.Empty(t, sr.Applied)
}
