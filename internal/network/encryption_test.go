package network

import (
	"bytes"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() failed: %v", err)
	}

	if len(key) != passphraseLength {
		t.Errorf("Expected key length %d, got %d", passphraseLength, len(key))
	}

	// Generate another key and ensure they're different
	key2, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() failed: %v", err)
	}

	if bytes.Equal(key, key2) {
		t.Error("Generated keys should be different")
	}
}

func TestNewEncryption(t *testing.T) {
	validKey1, _ := GenerateKey()
	validKey2, _ := GenerateKey()

	t.Run("valid keys", func(t *testing.T) {
		enc, err := NewEncryption(validKey1, validKey2)
		if err != nil {
			t.Fatalf("NewEncryption() failed: %v", err)
		}
		if enc == nil {
			t.Error("Expected non-nil Encryption")
		}
	})

	t.Run("invalid encryption key length", func(t *testing.T) {
		invalidKey := make([]byte, 16) // Wrong length
		_, err := NewEncryption(invalidKey, validKey2)
		if err != ErrInvalidKeyLength {
			t.Errorf("Expected ErrInvalidKeyLength, got %v", err)
		}
	})

	t.Run("invalid decryption key length", func(t *testing.T) {
		invalidKey := make([]byte, 16) // Wrong length
		_, err := NewEncryption(validKey1, invalidKey)
		if err != ErrInvalidKeyLength {
			t.Errorf("Expected ErrInvalidKeyLength, got %v", err)
		}
	})
}

func TestEncryptDecrypt(t *testing.T) {
	key, _ := GenerateKey()
	enc, err := NewEncryption(key, key)
	if err != nil {
		t.Fatalf("NewEncryption() failed: %v", err)
	}

	testData := []byte("Hello, World!")

	t.Run("encrypt and decrypt", func(t *testing.T) {
		encrypted, err := enc.Encrypt(testData)
		if err != nil {
			t.Fatalf("Encrypt() failed: %v", err)
		}

		if bytes.Equal(encrypted, testData) {
			t.Error("Encrypted data should not equal original data")
		}

		decrypted, err := enc.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decrypt() failed: %v", err)
		}

		if !bytes.Equal(decrypted, testData) {
			t.Errorf("Decrypted data doesn't match original. Expected %s, got %s", testData, decrypted)
		}
	})

	t.Run("empty data", func(t *testing.T) {
		empty := []byte{}
		encrypted, err := enc.Encrypt(empty)
		if err != nil {
			t.Fatalf("Encrypt() failed: %v", err)
		}

		decrypted, err := enc.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decrypt() failed: %v", err)
		}

		if !bytes.Equal(decrypted, empty) {
			t.Error("Decrypted empty data should be empty")
		}
	})

	t.Run("decrypt too short data", func(t *testing.T) {
		shortData := make([]byte, 5) // Too short to contain nonce
		_, err := enc.Decrypt(shortData)
		if err == nil {
			t.Error("Expected error when decrypting too short data")
		}
	})

	t.Run("decrypt corrupted data", func(t *testing.T) {
		encrypted, _ := enc.Encrypt(testData)
		// Corrupt the last byte
		encrypted[len(encrypted)-1] ^= 0x01

		_, err := enc.Decrypt(encrypted)
		if err == nil {
			t.Error("Expected error when decrypting corrupted data")
		}
	})
}

func TestEncryptionDeterminism(t *testing.T) {
	key1, _ := GenerateKey()
	key2, _ := GenerateKey()
	enc, _ := NewEncryption(key1, key2)

	testData := []byte("Same data")

	encrypted1, _ := enc.Encrypt(testData)
	encrypted2, _ := enc.Encrypt(testData)

	if bytes.Equal(encrypted1, encrypted2) {
		t.Error("Same plaintext should produce different ciphertext due to random nonces")
	}
}
