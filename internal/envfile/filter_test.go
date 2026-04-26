package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func filterEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result := FilterByPrefix(filterEntries(), "APP_")
	assert.Len(t, result, 2)
	assert.Equal(t, "APP_NAME", result[0].Key)
	assert.Equal(t, "APP_PORT", result[1].Key)
}

func TestFilter_BySuffix(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{Suffix: "_KEY"})
	assert.Len(t, result, 1)
	assert.Equal(t, "API_KEY", result[0].Key)
}

func TestFilter_ByKeys(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{Keys: []string{"DB_HOST", "LOG_LEVEL"}})
	assert.Len(t, result, 2)
	keys := []string{result[0].Key, result[1].Key}
	assert.Contains(t, keys, "DB_HOST")
	assert.Contains(t, keys, "LOG_LEVEL")
}

func TestFilter_SecretsOnly(t *testing.T) {
	result := FilterSecrets(filterEntries())
	for _, e := range result {
		assert.True(t, IsSecret(e.Key), "expected %s to be a secret", e.Key)
	}
	assert.NotEmpty(t, result)
}

func TestFilter_NonSecretsOnly(t *testing.T) {
	result := FilterNonSecrets(filterEntries())
	for _, e := range result {
		assert.False(t, IsSecret(e.Key), "expected %s to not be a secret", e.Key)
	}
	assert.NotEmpty(t, result)
}

func TestFilter_EmptyOptions(t *testing.T) {
	entries := filterEntries()
	result := Filter(entries, FilterOptions{})
	assert.Equal(t, entries, result)
}

func TestFilter_NoMatch(t *testing.T) {
	result := FilterByPrefix(filterEntries(), "NONEXISTENT_")
	assert.Empty(t, result)
}

func TestFilter_PrefixAndSecretsOnly(t *testing.T) {
	result := Filter(filterEntries(), FilterOptions{Prefix: "DB_", SecretsOnly: true})
	for _, e := range result {
		assert.True(t, strings.HasPrefix(e.Key, "DB_"))
		assert.True(t, IsSecret(e.Key))
	}
}
