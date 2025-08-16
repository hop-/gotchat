package core

import "errors"

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityAlreadyExists = errors.New("entity already exists")
	ErrEntityFieldNotExist = errors.New("entity field does not exist")
)
