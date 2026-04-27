package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// TagEntry represents a key-value entry annotated with a tag.
type TagEntry struct {
	Key   string
	Value string
	Tag   string
}

// TagResult holds the output of a Tag operation.
type TagResult struct {
	Tagged   []TagEntry
	Skipped  []string
}

// Tag annotates entries whose keys match the given keys slice with the
// provided tag string. If keys is empty, all entries are tagged.
// Returns a TagResult describing what was tagged and what was skipped.
func Tag(entries []Entry, keys []string, tag string) (TagResult, error) {
	if tag == "" {
		return TagResult{}, fmt.Errorf("tag must not be empty")
	}

	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	result := TagResult{}

	for _, e := range entries {
		if len(keys) == 0 || keySet[e.Key] {
			result.Tagged = append(result.Tagged, TagEntry{
				Key:   e.Key,
				Value: e.Value,
				Tag:   tag,
			})
		} else {
			result.Skipped = append(result.Skipped, e.Key)
		}
	}

	return result, nil
}

// GroupByTag groups TagEntry items by their tag value.
func GroupByTag(tagged []TagEntry) map[string][]TagEntry {
	groups := make(map[string][]TagEntry)
	for _, t := range tagged {
		groups[t.Tag] = append(groups[t.Tag], t)
	}
	return groups
}

// TagSummary returns a human-readable summary of tag groups.
func TagSummary(tagged []TagEntry) string {
	groups := GroupByTag(tagged)
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, tag := range keys {
		sb.WriteString(fmt.Sprintf("[%s]\n", tag))
		for _, e := range groups[tag] {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, e.Value))
		}
	}
	return sb.String()
}
