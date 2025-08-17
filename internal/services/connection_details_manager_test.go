package services

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyManager_WrapKey(t *testing.T) {
	// Generate a random KEK
	kek := make([]byte, 32)
	_, err := rand.Read(kek)
	require.NoError(t, err)

	km := &KeyManager{kek: kek}

	// Test wrapping a key
	originalKey := []byte("test-key-1234567")
	wrappedKey, err := km.WrapKey(originalKey)

	assert.NoError(t, err)
	assert.NotNil(t, wrappedKey)
	assert.NotEqual(t, originalKey, wrappedKey)
	assert.Greater(t, len(wrappedKey), len(originalKey))
}

func TestKeyManager_UnwrapKey(t *testing.T) {
	// Generate a random KEK
	kek := make([]byte, 32)
	_, err := rand.Read(kek)
	require.NoError(t, err)

	km := &KeyManager{kek: kek}

	// Test wrapping and unwrapping a key
	originalKey := []byte("test-key-1234567")
	wrappedKey, err := km.WrapKey(originalKey)
	require.NoError(t, err)

	unwrappedKey, err := km.UnwrapKey(wrappedKey)

	assert.NoError(t, err)
	assert.Equal(t, originalKey, unwrappedKey)
}

func TestKeyManager_WrapUnwrapKey_RoundTrip(t *testing.T) {
	// Generate a random KEK
	kek := make([]byte, 32)
	_, err := rand.Read(kek)
	require.NoError(t, err)

	km := &KeyManager{kek: kek}

	testCases := []struct {
		name string
		key  []byte
	}{
		{"short key", []byte("abc")},
		{"medium key", []byte("this-is-a-medium-length-key")},
		{"long key", []byte("this-is-a-very-long-key-that-should-still-work-correctly")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wrappedKey, err := km.WrapKey(tc.key)
			require.NoError(t, err)

			unwrappedKey, err := km.UnwrapKey(wrappedKey)
			require.NoError(t, err)

			assert.Equal(t, tc.key, unwrappedKey)
		})
	}
}
