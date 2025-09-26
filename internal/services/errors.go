package services

import "fmt"

var (
	ErrFieldNotExist        = fmt.Errorf("field does not exist in the entity")
	ErrorInvalidInput       = fmt.Errorf("invalid input provided")
	ErrorInvalidCredentials = fmt.Errorf("invalid credentials provided")
	ErrorEntityExists       = fmt.Errorf("entity already exists")
)
