package envfile

import (
	"strings"
	"testing"
)

func baseResolveEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
		{Key: "DB_PASSWORD", Value: ""},
		{Key: "MISSING_KEY", Value: ""},
	}
}

func TestResolve_FileValueUsed(t *testing.T) {
	entries := []Entry{{Key: "APP_NAME", Value: "myapp"}}
	results, err := Resolve(entries, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != "myapp" || results[0].Source != SourceFile {
		t.Errorf("expected file source with value 'myapp', got %+v", results[0])
	}
}

func TestResolve_EnvFallback(t *testing.T) {
	t.Setenv("PORT", "9090")
	entries := []Entry{{Key: "PORT", Value: ""}}
	results, err := Resolve(entries, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != "9090" || results[0].Source != SourceEnv {
		t.Errorf("expected env source with value '9090', got %+v", results[0])
	}
}

func TestResolve_DefaultFallback(t *testing.T) {
	entries := []Entry{{Key: "TIMEOUT", Value: ""}}
	opts := ResolveOptions{Defaults: map[string]string{"TIMEOUT": "30s"}}
	results, err := Resolve(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != "30s" || results[0].Source != SourceDefault {
		t.Errorf("expected default source, got %+v", results[0])
	}
}

func TestResolve_PreferEnvOverFile(t *testing.T) {
	t.Setenv("APP_NAME", "from-env")
	entries := []Entry{{Key: "APP_NAME", Value: "from-file"}}
	opts := ResolveOptions{PreferEnv: true}
	results, err := Resolve(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != "from-env" || results[0].Source != SourceEnv {
		t.Errorf("expected env to win over file, got %+v", results[0])
	}
}

func TestResolve_MissingSource(t *testing.T) {
	entries := []Entry{{Key: "GHOST", Value: ""}}
	results, err := Resolve(entries, ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Source != SourceMissing {
		t.Errorf("expected missing source, got %+v", results[0])
	}
}

func TestResolve_FailOnMissing(t *testing.T) {
	entries := []Entry{{Key: "REQUIRED", Value: ""}}
	_, err := Resolve(entries, ResolveOptions{FailOnMissing: true})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestResolvedToEntries_SkipsMissing(t *testing.T) {
	resolved := []ResolvedEntry{
		{Key: "A", Value: "1", Source: SourceFile},
		{Key: "B", Value: "", Source: SourceMissing},
	}
	out := ResolvedToEntries(resolved)
	if len(out) != 1 || out[0].Key != "A" {
		t.Errorf("expected only entry A, got %+v", out)
	}
}

func TestResolvedSummary_MasksSecrets(t *testing.T) {
	resolved := []ResolvedEntry{
		{Key: "DB_PASSWORD", Value: "supersecret", Source: SourceFile},
		{Key: "APP_NAME", Value: "myapp", Source: SourceFile},
	}
	summary := ResolvedSummary(resolved)
	if strings.Contains(summary, "supersecret") {
		t.Error("summary should not contain raw secret value")
	}
	if !strings.Contains(summary, "myapp") {
		t.Error("summary should contain non-secret value")
	}
}
