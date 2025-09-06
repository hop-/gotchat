package services

import "github.com/hop-/gotchat/pkg/network"

type NewUnauthenticatedConnection struct {
	Id   string
	Conn network.AdvancedConn
}

type ConnectionEstablished struct {
	Id         string
	Conn       network.AdvancedConn
	PeerUserId string
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
