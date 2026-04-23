package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-sync/internal/envfile"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	original := "APP_NAME=myapp\nDB_PASSWORD=hunter2\nAPI_SECRET=topsecret\n"
	src := writeTempEncryptEnv(t, original)
	dir := t.TempDir()
	encFile := filepath.Join(dir, "enc.env")
	decFile := filepath.Join(dir, "dec.env")

	// Encrypt
	rootCmd.SetArgs([]string{"encrypt", src, "--key", encTestKey, "--output", encFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("encrypt step failed: %v", err)
	}

	// Verify encrypted file differs for secrets
	encEntries, err := envfile.Parse(encFile)
	if err != nil {
		t.Fatalf("parse enc file: %v", err)
	}
	for _, e := range encEntries {
		switch e.Key {
		case "DB_PASSWORD":
			if e.Value == "hunter2" {
				t.Error("DB_PASSWORD not encrypted")
			}
		case "API_SECRET":
			if e.Value == "topsecret" {
				t.Error("API_SECRET not encrypted")
			}
		case "APP_NAME":
			if e.Value != "myapp" {
				t.Errorf("APP_NAME should be unchanged, got %q", e.Value)
			}
		}
	}

	// Decrypt
	rootCmd.SetArgs([]string{"encrypt", encFile, "--key", encTestKey, "--decrypt", "--output", decFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("decrypt step failed: %v", err)
	}

	// Compare decrypted to original
	decEntries, err := envfile.Parse(decFile)
	if err != nil {
		t.Fatalf("parse dec file: %v", err)
	}
	origEntries, _ := envfile.Parse(src)
	origMap := map[string]string{}
	for _, e := range origEntries {
		origMap[e.Key] = e.Value
	}
	for _, e := range decEntries {
		if origMap[e.Key] != e.Value {
			t.Errorf("key %q: expected %q, got %q", e.Key, origMap[e.Key], e.Value)
		}
	}

	_ = os.Remove(encFile)
	_ = os.Remove(decFile)
}
