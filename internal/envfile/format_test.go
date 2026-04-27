package envfile

import (
	"strings"
	"testing"
)

func formatEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestFormat_CompactStyle(t *testing.T) {
	result := Format(formatEntries(), FormatOptions{Style: StyleCompact})
	if len(result.Lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result.Lines))
	}
	if result.Lines[0] != "APP_NAME=myapp" {
		t.Errorf("unexpected line: %s", result.Lines[0])
	}
}

func TestFormat_SpacedStyle(t *testing.T) {
	result := Format(formatEntries(), FormatOptions{Style: StyleSpaced})
	if result.Lines[0] != "APP_NAME = myapp" {
		t.Errorf("expected spaced format, got: %s", result.Lines[0])
	}
}

func TestFormat_AlignedStyle(t *testing.T) {
	result := Format(formatEntries(), FormatOptions{Style: StyleAligned})
	// All '=' signs should be at the same column
	positions := make([]int, len(result.Lines))
	for i, line := range result.Lines {
		positions[i] = strings.Index(line, "=")
	}
	for i := 1; i < len(positions); i++ {
		if positions[i] != positions[0] {
			t.Errorf("misaligned '=' at line %d: pos %d vs %d", i, positions[i], positions[0])
		}
	}
}

func TestFormat_SortKeys(t *testing.T) {
	result := Format(formatEntries(), FormatOptions{Style: StyleCompact, SortKeys: true})
	keys := make([]string, len(result.Lines))
	for i, line := range result.Lines {
		keys[i] = strings.Split(line, "=")[0]
	}
	if keys[0] != "APP_NAME" || keys[1] != "DB_PASSWORD" || keys[2] != "PORT" {
		t.Errorf("unexpected sort order: %v", keys)
	}
}

func TestFormat_MaskSecret(t *testing.T) {
	result := Format(formatEntries(), FormatOptions{Style: StyleCompact, MaskSecret: true})
	for _, line := range result.Lines {
		if strings.HasPrefix(line, "DB_PASSWORD=") {
			if !strings.Contains(line, "***") {
				t.Errorf("expected secret to be masked, got: %s", line)
			}
			return
		}
	}
	t.Error("DB_PASSWORD line not found")
}

func TestFormat_DefaultStyleIsCompact(t *testing.T) {
	result := Format(formatEntries(), FormatOptions{})
	if !strings.Contains(result.Lines[0], "=") || strings.Contains(result.Lines[0], " = ") {
		t.Errorf("expected compact style, got: %s", result.Lines[0])
	}
}

func TestFormat_ModifiedCount(t *testing.T) {
	result := Format(formatEntries(), FormatOptions{Style: StyleSpaced})
	if result.Modified != 3 {
		t.Errorf("expected 3 modified, got %d", result.Modified)
	}
}
