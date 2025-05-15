package services

import (
	"fmt"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/network"
)

type Server struct {
	address   string
	listener  *network.Listener
	em        *core.EventManager
	isRunning bool
}

func NewServer(address string, em *core.EventManager) *Server {
	return &Server{
		address,
		nil,
		em,
		false,
	}
}

func (s *Server) Init() error {
	// Start the server
	if s.listener != nil {
		return fmt.Errorf("server is already running")
	}
	t := network.NewTcpTransport()

	listener, err := t.Listen(s.address)
	if err != nil {
		return err
	}

	s.listener = listener

	return nil
}

func (s *Server) Run() {
	if s.isRunning {
		// TODO: Handle error
		return
	}

	if s.listener == nil {
		// TODO: Handle error
		return
	}

	s.isRunning = true

	for s.isRunning {
		conn, err := s.listener.Accept()
		if err != nil {
			// TODO: Handle error
			continue
		}

		go handleConnection(conn)
	}
}

func (s *Server) Close() error {
	s.isRunning = false

	if s.listener == nil {
		return nil
	}

	return s.listener.Close()
}

func (s *Server) Name() string {
	return "Server"
}

func handleConnection(conn *network.Conn) {
	defer conn.Close()

	// Handle the connection
	// TODO: Implement connection handling logic
}
