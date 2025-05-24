package core

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func generateUuid() string {
	id, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return ""
	}
	fmt.Println("Generated UUID:", id.String())
	return id.String()
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
