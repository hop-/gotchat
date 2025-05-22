package core

type Entity struct {
	Id string
}

type User struct {
	Entity
	Name      string
	LastLogin string
}
