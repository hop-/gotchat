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
