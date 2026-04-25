package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/envfile"
)

func writeTempSchemaFiles(t *testing.T, envContent string, schema envfile.Schema) (envPath, schemaPath string) {
	t.Helper()
	dir := t.TempDir()

	envPath = filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(schema)
	schemaPath = filepath.Join(dir, "schema.json")
	if err := os.WriteFile(schemaPath, data, 0644); err != nil {
		t.Fatal(err)
	}
	return
}

func runSchemaCmd(t *testing.T, envPath, schemaPath string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	schemaCmd.SetOut(&buf)
	schemaCmd.SetArgs([]string{"--file", envPath, "--schema", schemaPath})
	err := schemaCmd.Execute()
	return buf.String(), err
}

func TestSchemaCmd_PassesValidation(t *testing.T) {
	schema := envfile.Schema{
		Fields: []envfile.SchemaField{{Key: "PORT", Required: true, Pattern: `^\d+$`}},
	}
	envPath, schemaPath := writeTempSchemaFiles(t, "PORT=8080\n", schema)
	out, err := runSchemaCmd(t, envPath, schemaPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if out == "" {
		t.Error("expected success message")
	}
}

func TestSchemaCmd_FailsOnMissingRequired(t *testing.T) {
	schema := envfile.Schema{
		Fields: []envfile.SchemaField{{Key: "DB_URL", Required: true}},
	}
	envPath, schemaPath := writeTempSchemaFiles(t, "PORT=8080\n", schema)
	_, err := runSchemaCmd(t, envPath, schemaPath)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestSchemaCmd_MissingEnvFile(t *testing.T) {
	schema := envfile.Schema{}
	_, schemaPath := writeTempSchemaFiles(t, "", schema)
	_, err := runSchemaCmd(t, "/nonexistent/.env", schemaPath)
	if err == nil {
		t.Fatal("expected error for missing env file")
	}
}
