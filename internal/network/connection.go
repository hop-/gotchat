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

type AdvancedConn interface {
	Conn() BasicConn
	Read() (*Message, error)
	Write(m *Message) error
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
	// Read the message frame
	messageData, err := c.readFrame()
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

	// Write the message frame
	return c.WriteFrame(messageData)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) readFrame() ([]byte, error) {
	// Read the frame size
	var frameSize uint64
	err := binary.Read(c.conn, binary.LittleEndian, &frameSize)
	if err != nil {
		return nil, err
	}

	frameData := make([]byte, frameSize)
	// Read the frame data
	err = c.readAll(frameData)
	if err != nil {
		return nil, err
	}

	return frameData, nil
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

func (c *Conn) WriteFrame(frame []byte) error {
	// Write the frame size
	err := binary.Write(c.conn, binary.LittleEndian, uint64(len(frame)))
	if err != nil {
		return err
	}

	// Write the frame data
	return c.writeAll(frame)
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
