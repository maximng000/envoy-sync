package envfile

import (
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"C": "3"}

	result := Merge(base, override, PreferBase)

	if len(result.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(result.Conflicts))
	}
	if result.Merged["A"] != "1" || result.Merged["B"] != "2" || result.Merged["C"] != "3" {
		t.Errorf("unexpected merged map: %v", result.Merged)
	}
}

func TestMerge_ConflictPreferBase(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	override := map[string]string{"HOST": "prod.example.com"}

	result := Merge(base, override, PreferBase)

	if len(result.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(result.Conflicts))
	}
	if result.Merged["HOST"] != "localhost" {
		t.Errorf("expected base value 'localhost', got %q", result.Merged["HOST"])
	}
}

func TestMerge_ConflictPreferOverride(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	override := map[string]string{"HOST": "prod.example.com"}

	result := Merge(base, override, PreferOverride)

	if len(result.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(result.Conflicts))
	}
	if result.Merged["HOST"] != "prod.example.com" {
		t.Errorf("expected override value 'prod.example.com', got %q", result.Merged["HOST"])
	}
}

func TestMerge_OverrideOnlyKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	override := map[string]string{"A": "1", "B": "2"}

	result := Merge(base, override, PreferBase)

	if len(result.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(result.Conflicts))
	}
	if result.Merged["B"] != "2" {
		t.Errorf("expected B=2 in merged map")
	}
}

func TestMerge_ConflictDetails(t *testing.T) {
	base := map[string]string{"SECRET_KEY": "old"}
	override := map[string]string{"SECRET_KEY": "new"}

	result := Merge(base, override, PreferBase)

	if len(result.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict")
	}
	c := result.Conflicts[0]
	if c.Key != "SECRET_KEY" || c.BaseValue != "old" || c.OverrideValue != "new" {
		t.Errorf("unexpected conflict details: %+v", c)
	}
}
