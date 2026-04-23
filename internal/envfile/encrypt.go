package envfile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptValue encrypts a plaintext string using AES-GCM with the provided key.
// The key must be 16, 24, or 32 bytes long.
func EncryptValue(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptValue decrypts a base64-encoded AES-GCM ciphertext using the provided key.
func DecryptValue(encoded, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncryptSecrets returns a new Entry slice where all secret values are encrypted.
func EncryptSecrets(entries []Entry, key string) ([]Entry, error) {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		result[i] = e
		if IsSecret(e.Key) {
			enc, err := EncryptValue(e.Value, key)
			if err != nil {
				return nil, err
			}
			result[i].Value = enc
		}
	}
	return result, nil
}

// DecryptSecrets returns a new Entry slice where all secret values are decrypted.
func DecryptSecrets(entries []Entry, key string) ([]Entry, error) {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		result[i] = e
		if IsSecret(e.Key) {
			dec, err := DecryptValue(e.Value, key)
			if err != nil {
				return nil, err
			}
			result[i].Value = dec
		}
	}
	return result, nil
}
