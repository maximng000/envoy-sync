package envfile

import "sort"

// Stats holds summary statistics about a set of env entries.
type Stats struct {
	Total      int
	Secrets    int
	NonSecrets int
	Empty      int
	Unique     int
	Duplicates int
	Prefixes   map[string]int
}

// GatherStats computes statistics over a slice of Entry values.
func GatherStats(entries []Entry) Stats {
	seen := make(map[string]int)
	prefixCount := make(map[string]int)

	s := Stats{
		Prefixes: make(map[string]int),
	}

	for _, e := range entries {
		s.Total++
		seen[e.Key]++

		if IsSecret(e.Key) {
			s.Secrets++
		} else {
			s.NonSecrets++
		}

		if e.Value == "" {
			s.Empty++
		}

		if prefix := keyPrefix(e.Key); prefix != "" {
			prefixCount[prefix]++
		}
	}

	for k, count := range seen {
		_ = k
		if count == 1 {
			s.Unique++
		} else {
			s.Duplicates++
		}
	}

	for prefix, count := range prefixCount {
		s.Prefixes[prefix] = count
	}

	return s
}

// TopPrefixes returns the top-n prefixes by key count.
func TopPrefixes(s Stats, n int) []string {
	type kv struct {
		Key   string
		Count int
	}
	var pairs []kv
	for k, v := range s.Prefixes {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Key < pairs[j].Key
	})
	result := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, p.Key)
	}
	return result
}

// keyPrefix returns the underscore-delimited prefix of a key, or empty string.
func keyPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' && i > 0 {
			return key[:i]
		}
	}
	return ""
}
