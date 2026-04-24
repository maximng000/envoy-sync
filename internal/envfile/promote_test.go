package envfile

import (
	"testing"
)

func entries(pairs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func findResult(results []PromoteResult, key string) (PromoteResult, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return PromoteResult{}, false
}

func TestPromote_AddedKey(t *testing.T) {
	src := entries("NEW_KEY", "new_val")
	dst := entries("EXISTING", "val")
	_, results, err := Promote(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := findResult(results, "NEW_KEY")
	if !ok || r.Action != "added" {
		t.Errorf("expected NEW_KEY to be added, got %+v", r)
	}
}

func TestPromote_SkippedByDefault(t *testing.T) {
	src := entries("KEY", "new")
	dst := entries("KEY", "old")
	_, results, _ := Promote(src, dst, PromoteOptions{})
	r, _ := findResult(results, "KEY")
	if r.Action != "skipped" {
		t.Errorf("expected skipped, got %s", r.Action)
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := entries("KEY", "new")
	dst := entries("KEY", "old")
	out, results, _ := Promote(src, dst, PromoteOptions{Overwrite: true})
	r, _ := findResult(results, "KEY")
	if r.Action != "updated" {
		t.Errorf("expected updated, got %s", r.Action)
	}
	for _, e := range out {
		if e.Key == "KEY" && e.Value != "new" {
			t.Errorf("expected value 'new', got %s", e.Value)
		}
	}
}

func TestPromote_DryRunDoesNotMutate(t *testing.T) {
	src := entries("KEY", "new")
	dst := entries("KEY", "old")
	out, _, _ := Promote(src, dst, PromoteOptions{Overwrite: true, DryRun: true})
	for _, e := range out {
		if e.Key == "KEY" && e.Value != "old" {
			t.Errorf("dry run mutated dst: got %s", e.Value)
		}
	}
}

func TestPromote_UnchangedKey(t *testing.T) {
	src := entries("KEY", "same")
	dst := entries("KEY", "same")
	_, results, _ := Promote(src, dst, PromoteOptions{})
	r, _ := findResult(results, "KEY")
	if r.Action != "unchanged" {
		t.Errorf("expected unchanged, got %s", r.Action)
	}
}

func TestPromote_SecretMasked(t *testing.T) {
	src := entries("SECRET_TOKEN", "abc123")
	dst := entries("OTHER", "val")
	_, results, _ := Promote(src, dst, PromoteOptions{MaskSecrets: true})
	r, _ := findResult(results, "SECRET_TOKEN")
	if r.NewValue != "***" {
		t.Errorf("expected masked value, got %s", r.NewValue)
	}
}

func TestPromote_NilSrcError(t *testing.T) {
	_, _, err := Promote(nil, entries("K", "v"), PromoteOptions{})
	if err == nil {
		t.Error("expected error for nil src")
	}
}
