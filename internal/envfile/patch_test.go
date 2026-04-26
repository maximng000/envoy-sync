package envfile

import (
	"testing"
)

func patchEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestPatch_SetExistingKey(t *testing.T) {
	out, results, err := Patch(patchEntries(), []PatchOp{{Op: "set", Key: "PORT", Value: "9090"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied != true {
		t.Errorf("expected applied=true")
	}
	for _, e := range out {
		if e.Key == "PORT" && e.Value != "9090" {
			t.Errorf("expected PORT=9090, got %s", e.Value)
		}
	}
}

func TestPatch_SetNewKey(t *testing.T) {
	out, _, err := Patch(patchEntries(), []PatchOp{{Op: "set", Key: "NEW_KEY", Value: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, e := range out {
		if e.Key == "NEW_KEY" && e.Value == "hello" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected NEW_KEY to be added")
	}
}

func TestPatch_DeleteExistingKey(t *testing.T) {
	out, results, err := Patch(patchEntries(), []PatchOp{{Op: "delete", Key: "APP_ENV"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Errorf("expected applied=true for delete")
	}
	for _, e := range out {
		if e.Key == "APP_ENV" {
			t.Errorf("expected APP_ENV to be removed")
		}
	}
}

func TestPatch_DeleteMissingKey(t *testing.T) {
	_, results, err := Patch(patchEntries(), []PatchOp{{Op: "delete", Key: "MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Errorf("expected applied=false for missing key delete")
	}
	if results[0].Reason == "" {
		t.Errorf("expected a reason for unapplied delete")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	out, results, err := Patch(patchEntries(), []PatchOp{{Op: "rename", Key: "PORT", NewKey: "HTTP_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Errorf("expected applied=true for rename")
	}
	for _, e := range out {
		if e.Key == "PORT" {
			t.Errorf("old key PORT should not exist after rename")
		}
		if e.Key == "HTTP_PORT" {
			return
		}
	}
	t.Errorf("expected HTTP_PORT to exist after rename")
}

func TestPatch_RenameEmptyNewKey(t *testing.T) {
	_, _, err := Patch(patchEntries(), []PatchOp{{Op: "rename", Key: "PORT", NewKey: ""}})
	if err == nil {
		t.Errorf("expected error for rename with empty new_key")
	}
}

func TestPatch_UnknownOp(t *testing.T) {
	_, _, err := Patch(patchEntries(), []PatchOp{{Op: "upsert", Key: "X"}})
	if err == nil {
		t.Errorf("expected error for unknown op")
	}
}
