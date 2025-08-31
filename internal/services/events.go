package services

import "github.com/hop-/gotchat/pkg/network"

type NewConnection struct {
	Id   string
	Conn *network.Conn
}

type ConnectionClosed struct {
	Id string
}

type ConnectionAcceptError struct {
	Err error
}

type ConnectionFailed struct {
	Err error
}

type NewMessage struct {
	ConnId string
	// TODO: add message payload
}

type MessageReadError struct {
	ConnId string
	Err    error
}
