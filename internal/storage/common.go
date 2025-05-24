package storage

import "fmt"

var (
	ErrFieldNotExist = fmt.Errorf("field does not exist in the entity")
	ErrNotFound      = fmt.Errorf("entity not found")
)
