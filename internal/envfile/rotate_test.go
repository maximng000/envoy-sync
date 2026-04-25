package envfile

import (
	"errors"
	"strings"
	"testing"
)

func TestRotate_SecretKeyRotated(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PASSWORD", Value: "old-pass"},
	}
	fn := func(key, old string) (string, error) { return "new-pass", nil }

	updated, results, err := Rotate(entries, fn, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated[0].Value != "new-pass" {
		t.Errorf("expected new-pass, got %s", updated[0].Value)
	}
	if !results[0].Rotated {
		t.Error("expected Rotated=true")
	}
}

func TestRotate_NonSecretSkipped(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
	}
	fn := func(key, old string) (string, error) { return "should-not-be-used", nil }

	updated, results, err := Rotate(entries, fn, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated[0].Value != "myapp" {
		t.Errorf("value should be unchanged, got %s", updated[0].Value)
	}
	if !results[0].Skipped {
		t.Error("expected Skipped=true")
	}
}

func TestRotate_ForceKeyRotatesNonSecret(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "old"},
	}
	fn := func(key, old string) (string, error) { return "new", nil }

	updated, results, err := Rotate(entries, fn, []string{"APP_NAME"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated[0].Value != "new" {
		t.Errorf("expected new, got %s", updated[0].Value)
	}
	if !results[0].Rotated {
		t.Error("expected Rotated=true for forced key")
	}
}

func TestRotate_FnErrorPropagates(t *testing.T) {
	entries := []Entry{
		{Key: "API_SECRET", Value: "val"},
	}
	fn := func(key, old string) (string, error) {
		return "", errors.New("generation failed")
	}

	_, _, err := Rotate(entries, fn, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "API_SECRET") {
		t.Errorf("error should mention key, got: %v", err)
	}
}

func TestRotate_MixedEntries(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "prod"},
		{Key: "DB_SECRET", Value: "old-secret"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
	fn := func(key, old string) (string, error) { return old + "-rotated", nil }

	updated, results, err := Rotate(entries, fn, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updated) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(updated))
	}
	if updated[1].Value != "old-secret-rotated" {
		t.Errorf("expected old-secret-rotated, got %s", updated[1].Value)
	}
	rotatedCount := 0
	for _, r := range results {
		if r.Rotated {
			rotatedCount++
		}
	}
	if rotatedCount != 1 {
		t.Errorf("expected 1 rotated, got %d", rotatedCount)
	}
}
