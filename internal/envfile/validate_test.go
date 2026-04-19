package envfile

import (
	"testing"
)

func TestValidateAgainstSchema_AllPresent(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	schema := map[string]string{"DB_HOST": "", "DB_PORT": ""}
	result := ValidateAgainstSchema(env, schema)
	if !result.OK() {
		t.Fatalf("expected no errors, got %v", result.Errors)
	}
}

func TestValidateAgainstSchema_MissingKey(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	schema := map[string]string{"DB_HOST": "", "DB_PORT": ""}
	result := ValidateAgainstSchema(env, schema)
	if result.OK() {
		t.Fatal("expected errors for missing key")
	}
	if len(result.Errors) != 1 || result.Errors[0].Key != "DB_PORT" {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
}

func TestValidateAgainstSchema_EmptyValue(t *testing.T) {
	env := map[string]string{"DB_HOST": "", "DB_PORT": "5432"}
	schema := map[string]string{"DB_HOST": "", "DB_PORT": ""}
	result := ValidateAgainstSchema(env, schema)
	if result.OK() {
		t.Fatal("expected error for empty value")
	}
}

func TestValidateKeys_NoEmpty(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := ValidateKeys(env)
	if !result.OK() {
		t.Fatalf("expected no errors, got %v", result.Errors)
	}
}

func TestValidateKeys_WithEmpty(t *testing.T) {
	env := map[string]string{"FOO": "", "BAR": "value"}
	result := ValidateKeys(env)
	if result.OK() {
		t.Fatal("expected error for empty value")
	}
	if result.Errors[0].Key != "FOO" {
		t.Fatalf("expected FOO, got %s", result.Errors[0].Key)
	}
}

func TestValidationError_Error(t *testing.T) {
	e := ValidationError{Key: "MY_KEY", Message: "missing required key"}
	got := e.Error()
	expected := `key "MY_KEY": missing required key`
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
