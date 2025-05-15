package network

type BasicListener interface {
	Accept() (BasicConn, error)
	Close() error
}

type Listener struct {
	Listener BasicListener
}

func NewListener(listener BasicListener) *Listener {
	return &Listener{listener}
}

func (l *Listener) Accept() (*Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return NewConn(conn), nil
}

func (l *Listener) Close() error {
	return l.Listener.Close()
}
