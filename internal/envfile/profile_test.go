package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProfile_DotenvSuffix(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dev.env")
	_ = os.WriteFile(path, []byte("APP_ENV=dev\nDEBUG=true\n"), 0644)

	p, err := LoadProfile(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "dev" {
		t.Errorf("expected name 'dev', got %q", p.Name)
	}
	if p.Entries["APP_ENV"] != "dev" {
		t.Errorf("expected APP_ENV=dev, got %q", p.Entries["APP_ENV"])
	}
}

func TestLoadProfile_EnvDotPrefix(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.staging")
	_ = os.WriteFile(path, []byte("APP_ENV=staging\n"), 0644)

	p, err := LoadProfile(dir, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Entries["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", p.Entries["APP_ENV"])
	}
}

func TestLoadProfile_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadProfile(dir, "prod")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestListProfiles_MultipleProfiles(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "dev.env"), []byte("X=1"), 0644)
	_ = os.WriteFile(filepath.Join(dir, ".env.staging"), []byte("X=2"), 0644)
	_ = os.WriteFile(filepath.Join(dir, ".env"), []byte("X=0"), 0644) // should be ignored
	_ = os.WriteFile(filepath.Join(dir, "README.md"), []byte(""), 0644) // should be ignored

	profiles, err := ListProfiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d: %v", len(profiles), profiles)
	}
}

func TestDiffProfiles_DetectsChanges(t *testing.T) {
	a := &Profile{Name: "dev", Entries: map[string]string{"HOST": "localhost", "PORT": "8080"}}
	b := &Profile{Name: "prod", Entries: map[string]string{"HOST": "example.com", "PORT": "8080"}}

	diffs := DiffProfiles(a, b)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Key != "HOST" {
		t.Errorf("expected diff on HOST, got %q", diffs[0].Key)
	}
}
