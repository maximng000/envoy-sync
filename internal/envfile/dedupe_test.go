package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func dedupeEntries(pairs ...string) []Entry {
	entries := make([]Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestDedupe_NoDuplicates(t *testing.T) {
	entries := dedupeEntries("A", "1", "B", "2", "C", "3")
	res, err := Dedupe(entries, DedupeKeepFirst)
	require.NoError(t, err)
	assert.Len(t, res.Entries, 3)
	assert.Empty(t, res.Duplicates)
}

func TestDedupe_KeepFirst(t *testing.T) {
	entries := dedupeEntries("KEY", "first", "OTHER", "x", "KEY", "second")
	res, err := Dedupe(entries, DedupeKeepFirst)
	require.NoError(t, err)
	assert.Len(t, res.Entries, 2)
	require.Len(t, res.Duplicates, 1)
	assert.Equal(t, "first", res.Duplicates[0].Kept)
	assert.Equal(t, []string{"first", "second"}, res.Duplicates[0].Values)
}

func TestDedupe_KeepLast(t *testing.T) {
	entries := dedupeEntries("KEY", "first", "KEY", "second", "KEY", "third")
	res, err := Dedupe(entries, DedupeKeepLast)
	require.NoError(t, err)
	assert.Len(t, res.Entries, 1)
	require.Len(t, res.Duplicates, 1)
	assert.Equal(t, "third", res.Duplicates[0].Kept)
	assert.Equal(t, "KEY", res.Duplicates[0].Key)
}

func TestDedupe_PreservesOrder(t *testing.T) {
	entries := dedupeEntries("B", "b1", "A", "a1", "B", "b2", "C", "c1")
	res, err := Dedupe(entries, DedupeKeepFirst)
	require.NoError(t, err)
	keys := make([]string, len(res.Entries))
	for i, e := range res.Entries {
		keys[i] = e.Key
	}
	assert.Equal(t, []string{"B", "A", "C"}, keys)
}

func TestDedupe_UnknownStrategy(t *testing.T) {
	entries := dedupeEntries("KEY", "val")
	_, err := Dedupe(entries, DedupeStrategy("unknown"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown dedupe strategy")
}

func TestDedupe_MultipleDuplicateKeys(t *testing.T) {
	entries := dedupeEntries("A", "a1", "B", "b1", "A", "a2", "B", "b2")
	res, err := Dedupe(entries, DedupeKeepLast)
	require.NoError(t, err)
	assert.Len(t, res.Entries, 2)
	assert.Len(t, res.Duplicates, 2)
}
