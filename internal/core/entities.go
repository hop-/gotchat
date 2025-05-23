package core

import "time"

type Entity interface {
	GetId() string
}

type BaseEntity struct {
	Id string
}

func (e BaseEntity) GetId() string {
	return e.Id
}

type User struct {
	BaseEntity
	Name      string
	UniqueId  string
	LastLogin time.Time
}

func NewUser(name string) *User {
	return &User{
		BaseEntity: BaseEntity{},
		UniqueId:   generateUuid(),
		Name:       name,
		LastLogin:  time.Now(),
	}
}
