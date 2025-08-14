package network

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const (
	passphraseLength = 32
	keyLength        = 32
)

var (
	ErrInvalidKeyLength = fmt.Errorf("invalid key length")
)

func GenerateKey() ([]byte, error) {
	key := make([]byte, passphraseLength)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return key, nil
}

type Encryption struct {
	encryptionKey []byte
	decryptionKey []byte
	gcmEncrypt    cipher.AEAD
	gcmDecrypt    cipher.AEAD
}

func NewEncryption(encryptionKey, decryptionKey []byte) (*Encryption, error) {
	// Ensure the keys are of the correct length
	if len(encryptionKey) != keyLength || len(decryptionKey) != keyLength {
		return nil, ErrInvalidKeyLength
	}

	// Create AES cipher blocks
	blockEncrypt, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	blockDecrypt, err := aes.NewCipher(decryptionKey)
	if err != nil {
		return nil, err
	}

	gcmEncrypt, err := cipher.NewGCM(blockEncrypt)
	if err != nil {
		return nil, err
	}
	gcmDecrypt, err := cipher.NewGCM(blockDecrypt)
	if err != nil {
		return nil, err
	}

	return &Encryption{
		encryptionKey,
		decryptionKey,
		gcmEncrypt,
		gcmDecrypt,
	}, nil
}

func (e *Encryption) Encrypt(data []byte) ([]byte, error) {
	// Generate a random nonce
	nonce := make([]byte, e.gcmEncrypt.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	cipherData := e.gcmEncrypt.Seal(nil, nonce, data, nil)

	// Prepend nonce to ciphertext for transmission
	encryptedData := make([]byte, len(nonce)+len(cipherData))
	copy(encryptedData[:len(nonce)], nonce)
	copy(encryptedData[len(nonce):], cipherData)

	return encryptedData, nil
}

func (e *Encryption) Decrypt(encryptedData []byte) ([]byte, error) {
	nonceSize := e.gcmDecrypt.NonceSize()

	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	// Extract nonce and ciphertext
	nonce := encryptedData[:nonceSize]
	cipherData := encryptedData[nonceSize:]

	// Decrypt the data
	data, err := e.gcmDecrypt.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return data, nil
}
