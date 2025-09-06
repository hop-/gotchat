package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/pkg/log"
	"github.com/hop-/gotchat/pkg/network"
)

// Connection states
const (
	ConnectionStateKnown   = "KNOWN"
	ConnectionStateUnknown = "UNKNOWN"
)

// ConnectionManager Service
type ConnectionManager struct {
	AtomicRunningStatus
	// Mutex
	mu sync.RWMutex

	eventEmitter core.EventEmitter
	server       *Server

	// User controller
	userController *UserController

	// User manager
	userManager *UserManager

	// Connection details manager
	connectionDetailsManager *ConnectionDetailsManager
}

func NewConnectionManager(eventEmitter core.EventEmitter, server *Server, userManager *UserManager, connectionDetailsManager *ConnectionDetailsManager) *ConnectionManager {
	return &ConnectionManager{
		AtomicRunningStatus{},
		sync.RWMutex{},
		eventEmitter,
		server,
		nil,
		userManager,
		connectionDetailsManager,
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

	cm.setRunningStatus(true)

	// If no server is configured, run in client-only mode
	if cm.server == nil {
		log.Infof("No server configured, running in client-only mode")

		return
	}

	cm.runServer(ctx, wg)
}

// Close implements core.Service.
func (cm *ConnectionManager) Close() error {
	cm.setRunningStatus(false)
	if cm.server != nil {
		cm.server.Close()
	}

	cm.mu.RLock()
	defer cm.mu.RUnlock()
	// Close the user controller if it exists
	if cm.userController != nil {
		if err := cm.userController.Close(); err != nil {
			log.Errorf("Failed to close user controller: %v", err)
		}
	}

	return nil
}

// MapEventToCommands maps incoming events to their corresponding commands for the ConnectionManager.
func (cm *ConnectionManager) MapEventToCommands(event core.Event) []core.Command {
	var commands []core.Command
	switch e := event.(type) {
	case core.ConnectEvent:
		address := fmt.Sprintf("%s:%s", e.Host, e.Port)
		commands = append(commands, &Connect{cm, address})
		// TODO: utilize the returned connection ID
	case core.UserLoggedInEvent:
		commands = append(commands, &ChangeUserController{cm, e.User})
	case core.UserLoggedOutEvent:
		commands = append(commands, &RemoveUserController{cm})
	}

	return commands
}

func (cm *ConnectionManager) Connect(address string) (string, error) {
	cm.mu.RLock()
	if !cm.isRunning() || cm.userController == nil {
		log.Errorf("ConnectionManager is not running or user controller is not initialized")

		cm.mu.RUnlock()
		return "", fmt.Errorf("connection manager is not running or user controller is not initialized")
	}
	cm.mu.RUnlock()

	client := NewClient(address)
	conn, err := client.Connect()
	if err != nil {
		cm.emitEvent(ConnectionFailed{err})

		return "", err
	}

	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if cm.userController == nil {
		conn.Close()

		return "", fmt.Errorf("user controller is not initialized")
	}
	connId := cm.userController.Register(conn, true)

	return connId, nil
}

func (cm *ConnectionManager) emitEvent(event core.Event) {
	cm.eventEmitter.Emit(event)
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

			cm.mu.RLock()
			defer cm.mu.RUnlock()
			if cm.userController != nil {
				cm.userController.Register(conn, false)
			} else {
				// TODO: Handle connection without user controller
				log.Infof("No UserController initialized, closing connection")

				conn.Close()
			}
		}
	}
}

func (cm *ConnectionManager) changeUserController(user *core.User) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.userController != nil {
		log.Warnf("UserController is already initialized, removing previous user controller")
		if err := cm.userController.Close(); err != nil {
			log.Errorf("Failed to close previous user controller: %v", err)
		}
	}
	cm.userController = NewUserController(user, cm.eventEmitter, cm.userManager, cm.connectionDetailsManager)
	cm.userController.setRunningStatus(true)
	log.Infof("UserController initialized for user %s", user.Name)
}

func (cm *ConnectionManager) removeUserController() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.userController != nil {
		log.Infof("UserController is closing for user %s", cm.userController.user.Name)
		if err := cm.userController.Close(); err != nil {
			log.Errorf("Failed to close user controller: %v", err)
		}
		cm.userController = nil
	} else {
		log.Warnf("No UserController initialized, nothing to close")
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

	log.Infof("Connecting to %s", c.address)
	conn, err := t.Connect(c.address)
	if err != nil {
		return nil, err
	}
	log.Infof("Connected scuccessfully to %s", c.address)

	return conn, nil
}
