package core

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUuid(t *testing.T) {
	id := generateUuid()
	if id == "" {
		t.Errorf("Expected a non-empty UUID, got an empty string")
	}

	_, err := uuid.Parse(id)
	if err != nil {
		t.Errorf("Expected a valid UUID, got an invalid one: %v", err)
	}
}

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(hash) == 0 {
		t.Errorf("Expected a non-empty hash, got an empty string")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securepassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Errorf("Expected password to match hash, but it did not")
	}

	if CheckPasswordHash("wrongpassword", hash) {
		t.Errorf("Expected password not to match hash, but it did")
	}
}
