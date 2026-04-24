package envfile

import (
	"testing"
)

func TestInterpolate_SimpleReference(t *testing.T) {
	entries := []Entry{
		{Key: "BASE", Value: "/home/user"},
		{Key: "PATH", Value: "${BASE}/bin"},
	}

	got, err := Interpolate(entries, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[1].Value != "/home/user/bin" {
		t.Errorf("expected /home/user/bin, got %q", got[1].Value)
	}
}

func TestInterpolate_DollarSyntax(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "URL", Value: "http://$HOST:8080"},
	}

	got, err := Interpolate(entries, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[1].Value != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %q", got[1].Value)
	}
}

func TestInterpolate_MissingVarSilent(t *testing.T) {
	entries := []Entry{
		{Key: "URL", Value: "http://${MISSING_HOST}:8080"},
	}

	got, err := Interpolate(entries, InterpolateOptions{FailOnMissing: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0].Value != "http://:8080" {
		t.Errorf("expected empty substitution, got %q", got[0].Value)
	}
}

func TestInterpolate_MissingVarError(t *testing.T) {
	entries := []Entry{
		{Key: "URL", Value: "http://${MISSING_HOST}:8080"},
	}

	_, err := Interpolate(entries, InterpolateOptions{FailOnMissing: true})
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestInterpolate_ChainedReferences(t *testing.T) {
	entries := []Entry{
		{Key: "PROTO", Value: "https"},
		{Key: "HOST", Value: "example.com"},
		{Key: "BASE_URL", Value: "${PROTO}://${HOST}"},
		{Key: "API_URL", Value: "${BASE_URL}/api/v1"},
	}

	got, err := Interpolate(entries, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[3].Value != "https://example.com/api/v1" {
		t.Errorf("expected chained interpolation, got %q", got[3].Value)
	}
}

func TestInterpolate_NoReferences(t *testing.T) {
	entries := []Entry{
		{Key: "PLAIN", Value: "just a plain value"},
	}

	got, err := Interpolate(entries, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0].Value != "just a plain value" {
		t.Errorf("expected unchanged value, got %q", got[0].Value)
	}
}
