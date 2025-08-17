package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"

	"github.com/hop-/gotchat/internal/core"
)

// Master key manager
type KeyManager struct {
	kek []byte
}

func (k *KeyManager) WrapKey(key []byte) ([]byte, error) {
	block, err := aes.NewCipher(k.kek)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, key, nil)
	return ciphertext, nil
}

func (k *KeyManager) UnwrapKey(wrappedKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(k.kek)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := wrappedKey[:nonceSize], wrappedKey[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

// ConnectionDetails
type ConnectionDetails struct {
	HostUniqueId   string
	ClientUniqueId string
	EncryptionKey  []byte
	DecryptionKey  []byte
}

func newConnectionDetailsFromEntity(mk *KeyManager, entity *core.ConnectionDetails) (*ConnectionDetails, error) {
	encryptionKey, decryptionKey, err := retrieveKeysFromConnectionDetails(mk, entity)
	if err != nil {
		return nil, err
	}

	connDetails := &ConnectionDetails{
		HostUniqueId:   entity.HostUniqueId,
		ClientUniqueId: entity.ClientUniqueId,
		EncryptionKey:  encryptionKey,
		DecryptionKey:  decryptionKey,
	}

	return connDetails, nil
}

func updateKeysOfConnectionDetails(mk *KeyManager, details *core.ConnectionDetails, encryptionKey []byte, decryptionKey []byte) error {
	wrappedEncryptionKey, err := mk.WrapKey(encryptionKey)
	if err != nil {
		return err
	}

	details.EncryptionKey = base64.StdEncoding.EncodeToString(wrappedEncryptionKey)

	wrappedDecryptionKey, err := mk.WrapKey(decryptionKey)
	if err != nil {
		return err
	}

	details.DecryptionKey = base64.StdEncoding.EncodeToString(wrappedDecryptionKey)

	return nil
}

func retrieveKeysFromConnectionDetails(mk *KeyManager, details *core.ConnectionDetails) ([]byte, []byte, error) {
	wrappedEncryptionKey, err := base64.StdEncoding.DecodeString(details.EncryptionKey)
	if err != nil {
		return nil, nil, err
	}

	wrappedDecryptionKey, err := base64.StdEncoding.DecodeString(details.DecryptionKey)
	if err != nil {
		return nil, nil, err
	}

	encryptionKey, err := mk.UnwrapKey(wrappedEncryptionKey)
	if err != nil {
		return nil, nil, err
	}

	decryptionKey, err := mk.UnwrapKey(wrappedDecryptionKey)
	if err != nil {
		return nil, nil, err
	}

	return encryptionKey, decryptionKey, nil
}

// ConnectionDetailsManager manages connection details for hosts and clients
type ConnectionDetailsManager struct {
	eventEmitter core.EventEmitter
	repo         core.Repository[core.ConnectionDetails]
	mk           *KeyManager
}

func NewConnectionDetailsManager(eventEmitter core.EventEmitter, connectionDetailsRepo core.Repository[core.ConnectionDetails]) *ConnectionDetailsManager {
	return &ConnectionDetailsManager{
		eventEmitter,
		connectionDetailsRepo,
		&KeyManager{
			kek: []byte("this-is-a-very-secure-key-------"),
		}, // TODO: Hardcoded
	}
}

// Init implements core.Service.
func (m *ConnectionDetailsManager) Init() error {
	// Nothing to do
	return nil
}

// MapEventToCommands implements core.Service.
func (m *ConnectionDetailsManager) MapEventToCommands(event core.Event) []core.Command {
	return nil
}

// Name implements core.Service.
func (m *ConnectionDetailsManager) Name() string {
	return "ConnectionDetailsManager"
}

// Run implements core.Service.
func (m *ConnectionDetailsManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	// This service does not run any background tasks.
}

// Close implements core.Service.
func (m *ConnectionDetailsManager) Close() error {
	// Nothing to do
	return nil
}

func (m *ConnectionDetailsManager) GetConnectionDetails(host string, client string) (*ConnectionDetails, error) {
	details, err := m.getConnectionDetails(host, client)
	if details == nil {
		return nil, err
	}

	return newConnectionDetailsFromEntity(m.mk, details)
}

func (m *ConnectionDetailsManager) UpsertConnectionDetails(host string, client string, encryptionKey []byte, decryptionKey []byte) (*ConnectionDetails, error) {
	details, err := m.getConnectionDetails(host, client)
	if err != nil {
		return nil, err
	}

	// Check if details is already created
	if details == nil {
		// Create details if not exist
		details = &core.ConnectionDetails{
			HostUniqueId:   host,
			ClientUniqueId: client,
		}

		err = updateKeysOfConnectionDetails(m.mk, details, encryptionKey, decryptionKey)
		if err != nil {
			return nil, err
		}

		err = m.repo.Create(details)
		if err != nil {
			return nil, err
		}

		return newConnectionDetailsFromEntity(m.mk, details)
	}

	// Update existing details
	err = updateKeysOfConnectionDetails(m.mk, details, encryptionKey, decryptionKey)
	if err != nil {
		return nil, err
	}

	err = m.repo.Update(details)
	if err != nil {
		return nil, err
	}

	return newConnectionDetailsFromEntity(m.mk, details)
}

func (m *ConnectionDetailsManager) getConnectionDetails(host string, client string) (*core.ConnectionDetails, error) {
	// TODO: use better approach when available (GetOneWhere, GetOneByMany, etc.)
	detailsList, err := m.repo.GetAllBy("host_unique_id", host)
	if err != nil {
		if err == core.ErrEntityNotFound {
			return nil, nil
		}

		return nil, err
	}

	for _, details := range detailsList {
		if details.ClientUniqueId == client {
			return details, nil
		}
	}

	return nil, nil
}
