package services

import "fmt"

var (
	ErrFieldNotExist        = fmt.Errorf("field does not exist in the entity")
	ErrNotFound             = fmt.Errorf("entity not found")
	ErrorInvalidInput       = fmt.Errorf("invalid input provided")
	ErrorInvalidCredentials = fmt.Errorf("invalid credentials provided")
)
