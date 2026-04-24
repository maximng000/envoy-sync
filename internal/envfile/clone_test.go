package envfile

import (
	"testing"
)

func TestClone_AllKeys(t *testing.T) {
	src := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	dst := map[string]string{}
	out, results := Clone(src, dst, CloneOptions{})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestClone_SkipExisting(t *testing.T) {
	src := map[string]string{"PORT": "9090"}
	dst := map[string]string{"PORT": "8080"}
	out, results := Clone(src, dst, CloneOptions{Overwrite: false})
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT to remain 8080, got %s", out["PORT"])
	}
	if !results[0].Skipped {
		t.Error("expected result to be marked skipped")
	}
}

func TestClone_OverwriteExisting(t *testing.T) {
	src := map[string]string{"PORT": "9090"}
	dst := map[string]string{"PORT": "8080"}
	out, results := Clone(src, dst, CloneOptions{Overwrite: true})
	if out["PORT"] != "9090" {
		t.Errorf("expected PORT to be 9090, got %s", out["PORT"])
	}
	if results[0].Skipped {
		t.Error("expected result not to be skipped")
	}
}

func TestClone_FilterKeys(t *testing.T) {
	src := map[string]string{"APP_NAME": "myapp", "PORT": "8080", "DEBUG": "true"}
	dst := map[string]string{}
	out, results := Clone(src, dst, CloneOptions{FilterKeys: []string{"PORT"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", out["PORT"])
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestClone_SecretMaskedInResult(t *testing.T) {
	src := map[string]string{"SECRET_KEY": "supersecret"}
	dst := map[string]string{}
	_, results := Clone(src, dst, CloneOptions{MaskSecrets: true})
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if results[0].OldEnv == "supersecret" {
		t.Error("expected secret value to be masked in result")
	}
}

func TestClone_DoesNotMutateDst(t *testing.T) {
	src := map[string]string{"NEW_KEY": "value"}
	dst := map[string]string{"EXISTING": "yes"}
	Clone(src, dst, CloneOptions{})
	if _, ok := dst["NEW_KEY"]; ok {
		t.Error("Clone must not mutate the original dst map")
	}
}
