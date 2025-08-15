package network

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
	// Read the message frame
	encryptedMessageData, err := c.conn.readFrame()
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

	// Write the message frame
	return c.conn.WriteFrame(encryptedMessageData)
}

func (c *SecureConn) Close() error {
	return c.conn.Close()
}
