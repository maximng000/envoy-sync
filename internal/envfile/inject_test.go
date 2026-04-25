package envfile

import (
	"os"
	"testing"
)

func TestInject_SetsEnvVars(t *testing.T) {
	entries := map[string]string{"INJECT_FOO": "bar", "INJECT_BAZ": "qux"}
	t.Cleanup(func() {
		os.Unsetenv("INJECT_FOO")
		os.Unsetenv("INJECT_BAZ")
	})

	result, err := Inject(entries, InjectOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(result.Injected))
	}
	if os.Getenv("INJECT_FOO") != "bar" {
		t.Errorf("expected INJECT_FOO=bar")
	}
}

func TestInject_SkipsExistingWithoutOverwrite(t *testing.T) {
	os.Setenv("INJECT_EXISTING", "original")
	t.Cleanup(func() { os.Unsetenv("INJECT_EXISTING") })

	entries := map[string]string{"INJECT_EXISTING": "new"}
	result, err := Inject(entries, InjectOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "INJECT_EXISTING" {
		t.Errorf("expected INJECT_EXISTING to be skipped")
	}
	if os.Getenv("INJECT_EXISTING") != "original" {
		t.Errorf("expected original value to be preserved")
	}
}

func TestInject_OverwriteExisting(t *testing.T) {
	os.Setenv("INJECT_OVR", "old")
	t.Cleanup(func() { os.Unsetenv("INJECT_OVR") })

	entries := map[string]string{"INJECT_OVR": "new"}
	result, err := Inject(entries, InjectOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 1 {
		t.Errorf("expected 1 injected")
	}
	if os.Getenv("INJECT_OVR") != "new" {
		t.Errorf("expected overwritten value")
	}
}

func TestInject_FilterKeys(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("INJECT_A")
		os.Unsetenv("INJECT_B")
	})
	entries := map[string]string{"INJECT_A": "1", "INJECT_B": "2"}
	result, err := Inject(entries, InjectOptions{Keys: []string{"INJECT_A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 1 || result.Injected[0] != "INJECT_A" {
		t.Errorf("expected only INJECT_A injected, got %v", result.Injected)
	}
	if os.Getenv("INJECT_B") != "" {
		t.Errorf("expected INJECT_B to remain unset")
	}
}

func TestInject_EmptyEntries(t *testing.T) {
	result, err := Inject(map[string]string{}, InjectOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 0 {
		t.Errorf("expected nothing injected")
	}
}
