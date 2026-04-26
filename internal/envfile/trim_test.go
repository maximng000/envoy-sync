package envfile

import (
	"testing"
)

func TestTrim_RemovesLeadingAndTrailingSpaces(t *testing.T) {
	entries := map[string]string{
		"HOST": "  localhost  ",
		"PORT": "8080",
		"NAME": "\t my-app \t",
	}

	result, changes := Trim(entries)

	if result["HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", result["HOST"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", result["PORT"])
	}
	if result["NAME"] != "my-app" {
		t.Errorf("expected 'my-app', got %q", result["NAME"])
	}
	if len(changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(changes))
	}
}

func TestTrim_NoChangesWhenClean(t *testing.T) {
	entries := map[string]string{
		"KEY": "value",
		"OTHER": "clean",
	}

	_, changes := Trim(entries)

	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}

func TestTrim_EmptyValueUnchanged(t *testing.T) {
	entries := map[string]string{
		"EMPTY": "",
	}

	result, changes := Trim(entries)

	if result["EMPTY"] != "" {
		t.Errorf("expected empty string, got %q", result["EMPTY"])
	}
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}

func TestTrim_OriginalUnmodified(t *testing.T) {
	entries := map[string]string{
		"KEY": "  spaced  ",
	}

	Trim(entries)

	if entries["KEY"] != "  spaced  " {
		t.Error("original entries map should not be modified")
	}
}

func TestTrimKeys_OnlySpecifiedKeys(t *testing.T) {
	entries := map[string]string{
		"A": "  hello  ",
		"B": "  world  ",
		"C": "  skip  ",
	}

	result, changes := TrimKeys(entries, []string{"A", "B"})

	if result["A"] != "hello" {
		t.Errorf("expected 'hello', got %q", result["A"])
	}
	if result["B"] != "world" {
		t.Errorf("expected 'world', got %q", result["B"])
	}
	if result["C"] != "  skip  " {
		t.Errorf("expected '  skip  ' unchanged, got %q", result["C"])
	}
	if len(changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(changes))
	}
}

func TestTrimKeys_MissingKeyIgnored(t *testing.T) {
	entries := map[string]string{
		"EXISTS": "  val  ",
	}

	result, changes := TrimKeys(entries, []string{"EXISTS", "MISSING"})

	if result["EXISTS"] != "val" {
		t.Errorf("expected 'val', got %q", result["EXISTS"])
	}
	if len(changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(changes))
	}
}
