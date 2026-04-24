package envfile

import "sort"

// CompareResult holds the result of comparing two env maps.
type CompareResult struct {
	OnlyInA    []string          // keys present only in A
	OnlyInB    []string          // keys present only in B
	InBoth     []string          // keys present in both with identical values
	Different  []string          // keys present in both but with different values
	Summary    map[string]string // human-readable summary per key
}

// Compare performs a detailed comparison between two env entry maps.
// Secret values are masked in the Summary output.
func Compare(a, b map[string]Entry) CompareResult {
	seen := make(map[string]bool)
	result := CompareResult{
		Summary: make(map[string]string),
	}

	for k, ae := range a {
		seen[k] = true
		be, ok := b[k]
		if !ok {
			result.OnlyInA = append(result.OnlyInA, k)
			result.Summary[k] = "only in A: " + maskIfSecret(k, ae.Value)
			continue
		}
		if ae.Value == be.Value {
			result.InBoth = append(result.InBoth, k)
			result.Summary[k] = "same: " + maskIfSecret(k, ae.Value)
		} else {
			result.Different = append(result.Different, k)
			result.Summary[k] = "changed: " + maskIfSecret(k, ae.Value) + " → " + maskIfSecret(k, be.Value)
		}
	}

	for k, be := range b {
		if !seen[k] {
			result.OnlyInB = append(result.OnlyInB, k)
			result.Summary[k] = "only in B: " + maskIfSecret(k, be.Value)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.InBoth)
	sort.Strings(result.Different)

	return result
}
