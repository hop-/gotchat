package network

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"syscall"
)

type BasicConn interface {
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
	Close() error
}

type Conn struct {
	conn BasicConn
}

func NewConn(conn BasicConn) *Conn {
	return &Conn{conn}
}

func (c *Conn) Conn() BasicConn {
	return c.conn
}

func (c *Conn) Read() (*Message, error) {
	// Read the message size
	var messageSize uint64
	err := binary.Read(c.conn, binary.LittleEndian, &messageSize)
	if err != nil {
		return nil, err
	}

	messageData := make([]byte, messageSize)
	// Read the message data
	err = c.readAll(messageData)
	if err != nil {
		return nil, err
	}

	// Deserialize the message
	return DeserializeMessage(messageData)
}

func (c *Conn) Write(m *Message) error {
	// Serialize the message
	messageData, err := SerializeMessage(m)
	if err != nil {
		return err
	}

	// Write the message size
	err = binary.Write(c.conn, binary.LittleEndian, uint64(len(messageData)))
	if err != nil {
		return err
	}

	// Write the message data
	return c.writeAll(messageData)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) readAll(b []byte) error {
	offset := 0

	// Read whole message
	for offset < len(b) {
		size, err := c.conn.Read(b[offset:])
		if err != nil {
			return err
		}

		offset += size
		// TODO: check size == 0 case
	}

	return nil
}

func (c *Conn) writeAll(b []byte) error {
	offset := 0

	// Write whole message
	for offset < len(b) {
		size, err := c.conn.Write(b[offset:])
		if err != nil {
			return err
		}

		offset += size
		// TODO: check size == 0 case
	}

	return nil
}

func IsClosedError(err error) bool {
	switch {
	case
		errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	default:
		return false
	}
}
