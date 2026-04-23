package envfile

import (
	"testing"
)

const testKey = "0123456789abcdef" // 16-byte AES key

func TestEncryptDecryptValue_RoundTrip(t *testing.T) {
	plaintext := "super-secret-value"
	enc, err := EncryptValue(plaintext, testKey)
	if err != nil {
		t.Fatalf("EncryptValue error: %v", err)
	}
	if enc == plaintext {
		t.Error("encrypted value should differ from plaintext")
	}
	dec, err := DecryptValue(enc, testKey)
	if err != nil {
		t.Fatalf("DecryptValue error: %v", err)
	}
	if dec != plaintext {
		t.Errorf("expected %q, got %q", plaintext, dec)
	}
}

func TestEncryptValue_DifferentEachTime(t *testing.T) {
	plaintext := "value"
	enc1, _ := EncryptValue(plaintext, testKey)
	enc2, _ := EncryptValue(plaintext, testKey)
	if enc1 == enc2 {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestDecryptValue_InvalidBase64(t *testing.T) {
	_, err := DecryptValue("not-valid-base64!!!", testKey)
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestEncryptSecrets_OnlySecretsEncrypted(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "hunter2"},
		{Key: "API_SECRET", Value: "topsecret"},
	}
	result, err := EncryptSecrets(entries, testKey)
	if err != nil {
		t.Fatalf("EncryptSecrets error: %v", err)
	}
	if result[0].Value != "myapp" {
		t.Errorf("non-secret should be unchanged, got %q", result[0].Value)
	}
	if result[1].Value == "hunter2" {
		t.Error("DB_PASSWORD should be encrypted")
	}
	if result[2].Value == "topsecret" {
		t.Error("API_SECRET should be encrypted")
	}
}

func TestDecryptSecrets_RoundTrip(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "hunter2"},
	}
	encrypted, err := EncryptSecrets(entries, testKey)
	if err != nil {
		t.Fatalf("EncryptSecrets error: %v", err)
	}
	decrypted, err := DecryptSecrets(encrypted, testKey)
	if err != nil {
		t.Fatalf("DecryptSecrets error: %v", err)
	}
	for i, e := range entries {
		if decrypted[i].Value != e.Value {
			t.Errorf("key %q: expected %q, got %q", e.Key, e.Value, decrypted[i].Value)
		}
	}
}
