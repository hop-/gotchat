package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/log"
	"github.com/hop-/gotchat/internal/network"
)

// Events
type NewConnection struct {
	Id   string
	Conn *network.Conn
}

type ConnectionClosed struct {
	Id string
}

type ConnectionAcceptError struct {
	Err error
}

type ConnectionError struct {
	Err error
}

type NewMessage struct {
	ConnId string
	// TODO: add message payload
}

type MessageReadError struct {
	ConnId string
	Err    error
}

// ConnectionManager Service
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
	if cm.server != nil {
		err := cm.server.Init()
		if err != nil {
			return fmt.Errorf("failed to initialize server: %w", err)
		}
	}

	return nil
}

// Name implements core.Service.
func (cm *ConnectionManager) Name() string {
	return "ConnectionManager"
}

// Run implements core.Service.
func (cm *ConnectionManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	if cm.isRunning() {
		log.Errorf("ConnectionManager is already running")

		return
	}

	listener := cm.eventManager.Register(ctx)

	cm.setRunningStatus(true)

	// If no server is configured, run in client-only mode
	if cm.server == nil {
		log.Infof("No server configured, running in client-only mode")
		cm.handleApplicationEvents(ctx, wg, listener)

		return
	}

	go cm.handleApplicationEvents(ctx, wg, listener)
	cm.runServer(ctx, wg)
}

// Close implements core.Service.
func (cm *ConnectionManager) Close() error {
	cm.setRunningStatus(false)
	if cm.server != nil {
		cm.server.Close()
	}

	for _, conn := range cm.connections {
		if err := conn.Close(); err != nil {
			log.Errorf("Failed to close connection: %v", err)
		}
	}

	return nil
}

func (cm *ConnectionManager) Connect(address string) (string, error) {
	log.Infof("Connecting to %s", address)
	client := NewClient(address)
	conn, err := client.Connect()
	if err != nil {
		cm.emitEvent(ConnectionError{err})

		return "", err
	}
	log.Infof("Connected successfully")

	connId := cm.addConnection(conn)
	go cm.handleConnection(connId, conn)

	return connId, nil
}

func (cm *ConnectionManager) emitEvent(event core.Event) {
	cm.eventManager.Emit(event)
}

func (cm *ConnectionManager) addConnection(conn *network.Conn) string {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	id := generateUuid()
	cm.connections[id] = *conn

	cm.emitEvent(NewConnection{id, conn})

	return id
}

func (cm *ConnectionManager) removeConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.connections, id)

	cm.emitEvent(ConnectionClosed{id})
}

func (cm *ConnectionManager) handleApplicationEvents(
	ctx context.Context,
	wg *sync.WaitGroup,
	listener core.EventListener,
) {
	wg.Add(1)
	defer wg.Done()

	for cm.isRunning() {
		select {
		case <-ctx.Done():
			cm.setRunningStatus(false)
		default:
			e, err := listener.Next(ctx)
			if err != nil {
				log.Errorf("Failed to get next event: %v", err)
				continue
			}

			switch e := e.(type) {
			case core.ConnectEvent:
				address := fmt.Sprintf("%s:%s", e.Host, e.Port)
				cm.Connect(address)
				// TODO: utilize the returned connection ID
			}
		}
	}
}

func (cm *ConnectionManager) runServer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for cm.isRunning() && cm.server.IsInitialized() {
		select {
		case <-ctx.Done():
			// Context is done, stop the connection manager
			cm.setRunningStatus(false)
		default:
			conn, err := cm.server.Accept()
			if err != nil {
				log.Errorf("Failed to accept connection: %v", err)
				cm.emitEvent(ConnectionAcceptError{err})

				continue
			}

			connId := cm.addConnection(conn)

			go cm.handleConnection(connId, conn)
		}
	}
}

func (cm *ConnectionManager) handleConnection(connId string, conn *network.Conn) {
	defer conn.Close()
	defer cm.removeConnection(connId)

	// TODO: Handle handshake or any initial setup for the connection

	for cm.isRunning() {
		m, err := conn.Read()
		if err != nil {
			cm.emitEvent(MessageReadError{connId, err})
		}

		// TODO: Handle the message and emit an event
		_ = m
		cm.emitEvent(NewMessage{connId})
	}
}

// Network Server
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

	log.Infof("Starting server on %s", s.address)
	listener, err := t.Listen(s.address)
	if err != nil {
		return err
	}

	log.Infof("Server started successfully")
	s.listener = listener

	return nil
}

func (s *Server) IsInitialized() bool {
	return s.listener != nil
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

// Network Client
type Client struct {
	address string
}

func NewClient(address string) *Client {
	return &Client{
		address,
	}
}

func (c *Client) Connect() (*network.Conn, error) {
	t := network.NewTcpTransport()

	log.Infof("Connecting to server at %s", c.address)
	conn, err := t.Connect(c.address)
	if err != nil {
		return nil, err
	}
	log.Infof("Connected")

	return conn, nil
}
