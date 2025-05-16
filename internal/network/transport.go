package network

type Transport interface {
	Connect(address string) (*Conn, error)
	Listen(address string) (*Listener, error)
}
