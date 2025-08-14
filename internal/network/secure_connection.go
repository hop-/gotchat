package network

import "encoding/binary"

type SecureComponent interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type SecureConn struct {
	conn Conn
	sc   SecureComponent
}

func NewSecureConn(conn Conn, sc SecureComponent) *SecureConn {
	return &SecureConn{conn, sc}
}

func (c *SecureConn) Conn() BasicConn {
	return c.conn.Conn()
}

func (c *SecureConn) Read() (*Message, error) {
	// Read the message size
	var messageSize uint64
	err := binary.Read(c.conn.Conn(), binary.LittleEndian, &messageSize)
	if err != nil {
		return nil, err
	}

	encryptedMessageData := make([]byte, messageSize)
	// Read the message data
	err = c.conn.readAll(encryptedMessageData)
	if err != nil {
		return nil, err
	}

	// Decrypt the message
	messageData, err := c.sc.Decrypt(encryptedMessageData)
	if err != nil {
		return nil, err
	}

	// Deserialize the message
	return DeserializeMessage(messageData)
}

func (c *SecureConn) Write(m *Message) error {
	// Serialize the message
	messageData, err := SerializeMessage(m)
	if err != nil {
		return err
	}

	// Encrypt the message
	encryptedMessageData, err := c.sc.Encrypt(messageData)
	if err != nil {
		return err
	}

	// Write the message size
	err = binary.Write(c.conn.Conn(), binary.LittleEndian, uint64(len(encryptedMessageData)))
	if err != nil {
		return err
	}

	// Write the message data
	return c.conn.writeAll(encryptedMessageData)
}

func (c *SecureConn) Close() error {
	return c.conn.Close()
}
