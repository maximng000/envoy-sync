package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeEntries(kvs map[string]string) map[string]Entry {
	m := make(map[string]Entry, len(kvs))
	for k, v := range kvs {
		m[k] = Entry{Key: k, Value: v}
	}
	return m
}

func TestCompare_OnlyInA(t *testing.T) {
	a := makeEntries(map[string]string{"FOO": "bar", "ONLY_A": "yes"})
	b := makeEntries(map[string]string{"FOO": "bar"})
	r := Compare(a, b)
	assert.Equal(t, []string{"ONLY_A"}, r.OnlyInA)
	assert.Empty(t, r.OnlyInB)
}

func TestCompare_OnlyInB(t *testing.T) {
	a := makeEntries(map[string]string{"FOO": "bar"})
	b := makeEntries(map[string]string{"FOO": "bar", "ONLY_B": "yes"})
	r := Compare(a, b)
	assert.Equal(t, []string{"ONLY_B"}, r.OnlyInB)
	assert.Empty(t, r.OnlyInA)
}

func TestCompare_InBoth(t *testing.T) {
	a := makeEntries(map[string]string{"FOO": "bar"})
	b := makeEntries(map[string]string{"FOO": "bar"})
	r := Compare(a, b)
	assert.Equal(t, []string{"FOO"}, r.InBoth)
	assert.Empty(t, r.Different)
}

func TestCompare_Different(t *testing.T) {
	a := makeEntries(map[string]string{"FOO": "old"})
	b := makeEntries(map[string]string{"FOO": "new"})
	r := Compare(a, b)
	assert.Equal(t, []string{"FOO"}, r.Different)
	assert.Empty(t, r.InBoth)
}

func TestCompare_SecretMaskedInSummary(t *testing.T) {
	a := makeEntries(map[string]string{"SECRET_KEY": "mysecret"})
	b := makeEntries(map[string]string{"SECRET_KEY": "othersecret"})
	r := Compare(a, b)
	assert.Contains(t, r.Summary["SECRET_KEY"], "***")
	assert.NotContains(t, r.Summary["SECRET_KEY"], "mysecret")
	assert.NotContains(t, r.Summary["SECRET_KEY"], "othersecret")
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := Compare(map[string]Entry{}, map[string]Entry{})
	assert.Empty(t, r.OnlyInA)
	assert.Empty(t, r.OnlyInB)
	assert.Empty(t, r.InBoth)
	assert.Empty(t, r.Different)
}
