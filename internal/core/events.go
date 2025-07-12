package core

type Event any

type QuitEvent struct{}

type NewMessageEvent struct {
	Message string
	// TODO: Add more fields as needed
}

type ConnectEvent struct {
	Host string
	Port string
}

type UserCreatedEvent struct {
	User *User
}

type UserLoggedInEvent struct {
	User *User
}

type UserLoggedOutEvent struct {
	User *User
}

type UserUpdatedEvent struct {
	User *User
}
