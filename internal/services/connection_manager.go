package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/network"
)

// Events
type NewConnection struct {
	// TODO: Define fields for new connection event
}

type FialedNewConnection struct {
	// TODO: Define fields for failed new connection event
}

type ConnectionClosed struct {
	// TODO: Define fields for connection closed event
}

// The Service
type ConnectionManager struct {
	AtomicRunningStatus
	connections  map[string]network.Conn
	mu           sync.RWMutex
	eventManager *core.EventManager
	server       *Server
}

func NewConnectionManager(eventManager *core.EventManager, server *Server) *ConnectionManager {
	return &ConnectionManager{
		AtomicRunningStatus{},
		make(map[string]network.Conn),
		sync.RWMutex{},
		eventManager,
		server,
	}
}

// Init implements core.Service.
func (cm *ConnectionManager) Init() error {
	cm.server.Init()

	return nil
}

// Name implements core.Service.
func (cm *ConnectionManager) Name() string {
	return "ConnectionManager"
}

// Run implements core.Service.
func (cm *ConnectionManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	if cm.isRunning() {
		// TODO: Handle error
		return
	}

	listener := cm.eventManager.Register(ctx)

	cm.setRunningStatus(true)

	go cm.handleApplicationEvents(ctx, wg, listener)

	for cm.isRunning() {
		select {
		case <-ctx.Done():
			// Context is done, stop the connection manager
			cm.setRunningStatus(false)
		default:
			conn, err := cm.server.Accept()
			if err != nil {
				log.Printf("failed to accept connection: %v\n", err)
				cm.emitEvent(FialedNewConnection{})

				continue
			}

			cm.emitEvent(NewConnection{})

			go cm.handleConnection(conn)

			// TODO: Add connection to the map
		}
	}
}

// Close implements core.Service.
func (cm *ConnectionManager) Close() error {
	cm.setRunningStatus(false)
	cm.server.Close()

	return nil
}

func (cm *ConnectionManager) emitEvent(event core.Event) {
	cm.eventManager.Emit(event)
}

func (cm *ConnectionManager) handleApplicationEvents(ctx context.Context, wg *sync.WaitGroup, listener core.EventListener) {
	wg.Add(1)
	defer wg.Done()

	for cm.isRunning() {
		select {
		case <-ctx.Done():
			cm.setRunningStatus(false)
		default:
			e, err := listener.Next(ctx)
			if err != nil {
				// TODO: Handle error
				continue
			}

			switch e.(type) {
			// TODO: Handle specific events
			}
		}
	}
}

func (cm *ConnectionManager) handleConnection(conn *network.Conn) {
	defer conn.Close()
	defer cm.emitEvent(ConnectionClosed{})

	// Handle the connection
	// TODO: Implement connection handling logic
}

// The Network Server
type Server struct {
	address  string
	listener *network.Listener
}

func NewServer(address string) *Server {
	return &Server{
		address,
		nil,
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

func (s *Server) Accept() (*network.Conn, error) {
	if s.listener == nil {
		return nil, fmt.Errorf("server is not initialized")
	}

	conn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}

	return s.listener.Close()
}
