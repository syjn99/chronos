package aes

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := generateRandomKey() // AES-256
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}
	plaintext := []byte("Hello, AES!")

	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if bytes.Equal(plaintext, encrypted) {
		t.Fatal("Encrypted data is the same as plaintext")
	}

	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("Decrypted text doesn't match original. Expected %s but got %s", plaintext, decrypted)
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key, err := generateRandomKey() // AES-256
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}
	wrongKey, err := generateRandomKey() // AES-256
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}
	plaintext := []byte("Hello, AES!")

	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	_, err = Decrypt(wrongKey, encrypted)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
