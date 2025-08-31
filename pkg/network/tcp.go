package network

import "net"

// basicListenerWrapper adapts net.Listener to BasicListener.
type basicListenerWrapper struct {
	net.Listener
}

// Accept wraps the Accept method to return BasicConn.
func (b *basicListenerWrapper) Accept() (BasicConn, error) {
	return b.Listener.Accept()
}

type tcpTransport struct{}

func NewTcpTransport() Transport {
	return &tcpTransport{}
}

// Connect implements Transport.
func (t *tcpTransport) Connect(address string) (*Conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return NewConn(c), nil
}

// Listen implements Transport.
func (t *tcpTransport) Listen(address string) (*Listener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return NewListener(&basicListenerWrapper{Listener: l}), nil
}
