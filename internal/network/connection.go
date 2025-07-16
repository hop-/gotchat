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
	var headersSize, bodySize int64
	err := binary.Read(c.conn, binary.LittleEndian, &headersSize)
	if err != nil {
		return nil, err
	}

	err = binary.Read(c.conn, binary.LittleEndian, &bodySize)
	if err != nil {
		return nil, err
	}

	headersData := make([]byte, headersSize)
	// Read headers
	err = c.readAll(headersData)
	if err != nil {
		return nil, err
	}

	body := make([]byte, bodySize)
	// Read body
	err = c.readAll(body)
	if err != nil {
		return nil, err
	}

	return newMessageFromBytes(headersData, body)
}

func (c *Conn) Write(m *Message) error {

	headersData, body := m.toBytes()

	// Write the message sizes
	headerSize := len(headersData)
	err := binary.Write(c.conn, binary.LittleEndian, int64(headerSize))
	if err != nil {
		return err
	}
	bodySize := len(body)
	err = binary.Write(c.conn, binary.LittleEndian, int64(bodySize))
	if err != nil {
		return err
	}

	// Write whole message
	err = c.writeAll(headersData)
	if err != nil {
		return err
	}

	return c.writeAll(body)
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
