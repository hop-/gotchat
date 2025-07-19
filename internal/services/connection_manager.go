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

type ConnectionFailed struct {
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
	eventEmitter core.EventEmitter
	server       *Server

	// User controller
	userController *UserController
}

func NewConnectionManager(eventEmitter core.EventEmitter, server *Server) *ConnectionManager {
	return &ConnectionManager{
		AtomicRunningStatus{},
		eventEmitter,
		server,
		nil,
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
	if !cm.isRunning() || cm.userController == nil {
		log.Errorf("ConnectionManager is not running or user controller is not initialized")

		return "", fmt.Errorf("connection manager is not running or user controller is not initialized")
	}
	client := NewClient(address)
	conn, err := client.Connect()
	if err != nil {
		cm.emitEvent(ConnectionFailed{err})

		return "", err
	}

	if cm.userController == nil {
		conn.Close()

		return "", fmt.Errorf("user controller is not initialized")
	}
	connId := cm.userController.Register(conn, false)

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

			if cm.userController != nil {
				cm.userController.Register(conn, true)
			} else {
				// TODO: Handle connection without user controller
				log.Infof("No UserController initialized, closing connection")

				conn.Close()
			}
		}
	}
}

func (cm *ConnectionManager) changeUserController(user *core.User) {
	if cm.userController != nil {
		log.Warnf("UserController is already initialized, removing previous user controller")
		if err := cm.userController.Close(); err != nil {
			log.Errorf("Failed to close previous user controller: %v", err)
		}
	}
	cm.userController = NewUserController(cm.eventEmitter, user)
	cm.userController.setRunningStatus(true)
	log.Infof("UserController initialized for user %s", user.Name)
}

func (cm *ConnectionManager) removeUserController() {
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

// User Controller
type UserController struct {
	AtomicRunningStatus
	user         *core.User
	mu           sync.RWMutex
	connections  map[string]network.Conn
	eventEmitter core.EventEmitter
}

func NewUserController(eventEmitter core.EventEmitter, user *core.User) *UserController {
	return &UserController{
		AtomicRunningStatus{},
		user,
		sync.RWMutex{},
		make(map[string]network.Conn),
		eventEmitter,
	}
}

func (uc *UserController) Register(conn *network.Conn, isServer bool) string {
	log.Debugf("Registering new connection for user %s", uc.user.Name)
	connId := uc.addConnection(conn)
	go uc.handleConnection(connId, conn, isServer)

	return connId
}

func (uc *UserController) Close() error {
	uc.setRunningStatus(false)

	for _, conn := range uc.connections {
		if err := conn.Close(); err != nil {
			log.Errorf("Failed to close connection: %v", err)
		}
	}

	return nil
}

func (uc *UserController) emitEvent(event core.Event) {
	uc.eventEmitter.Emit(event)
}

func (uc *UserController) addConnection(conn *network.Conn) string {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	id := generateUuid()
	uc.connections[id] = *conn

	uc.emitEvent(NewConnection{id, conn})

	return id
}

func (uc *UserController) removeConnection(id string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	delete(uc.connections, id)

	uc.emitEvent(ConnectionClosed{id})
}

func (uc *UserController) handleConnection(connId string, conn *network.Conn, isServer bool) {
	// Ensure the connection is closed and removed when done
	defer conn.Close()
	defer uc.removeConnection(connId)

	// Handshake
	if isServer {
		err := uc.acceptHandshake(connId, conn)
		if err != nil {
			log.Errorf("Handshake failed for connection %s: %v", connId, err)
			uc.emitEvent(ConnectionFailed{err})

			return
		}
	} else {
		err := uc.initiateHandshake(connId, conn)
		if err != nil {
			log.Errorf("Handshake initiation failed for connection %s: %v", connId, err)
			uc.emitEvent(ConnectionFailed{err})

			return
		}
	}

	for uc.isRunning() {
		m, err := conn.Read()
		if err != nil {
			// Check connection close
			if network.IsClosedError(err) {
				log.Infof("Connection %s closed", connId)

				// Exit the loop
				break
			}

			uc.emitEvent(MessageReadError{connId, err})
		}

		// TODO: Handle the message and emit an event
		_ = m
		uc.emitEvent(NewMessage{connId})
	}
}

func (uc *UserController) initiateHandshake(connId string, conn *network.Conn) error {
	if !uc.isRunning() {
		log.Errorf("UserController is not running")

		return fmt.Errorf("user controller is not running")
	}
	if uc.user == nil {
		log.Errorf("User is not set")

		return fmt.Errorf("user is not set")
	}

	log.Infof("Initiating handshake for connection %s with user %s", connId, uc.user.Name)

	// Send a handshake message to the server
	err := conn.Write(network.NewMessage(map[string]string{
		"action": "authenticate",
		"user":   uc.user.Name,
		"userId": uc.user.UniqueId,
	}, nil))
	if err != nil {
		if network.IsClosedError(err) {
			log.Infof("Connection %s closed by peer before the handshake", connId)
		}

		return err
	}

	// Receive the handshake response
	msg, err := conn.Read()
	if err != nil {
		return err
	}

	// Checking message
	if action, ok := msg.Headers()["action"]; !ok || action != "authenticate" {
		return fmt.Errorf("handshake response missing or invalid action: %s", action)
	}

	var userId string
	var ok bool
	if userId, ok = msg.Headers()["userId"]; !ok {
		return fmt.Errorf("handshake response missing userId")
	}

	log.Debugf("Handshake user ID: %s", userId)
	// TODO: Implement the handshake logic

	return nil
}

func (uc *UserController) acceptHandshake(connId string, conn *network.Conn) error {
	if !uc.isRunning() {
		log.Errorf("UserController is not running")

		return fmt.Errorf("user controller is not running")
	}

	if uc.user == nil {
		log.Errorf("User is not set")

		return fmt.Errorf("user is not set")
	}

	log.Infof("Accepting handshake for connection %s", connId)

	// Receive the handshake response
	msg, err := conn.Read()
	if err != nil {
		return err
	}

	// Checking message
	if action, ok := msg.Headers()["action"]; !ok || action != "authenticate" {
		return fmt.Errorf("handshake response missing or invalid action: %s", action)
	}

	var userId string
	var ok bool
	if userId, ok = msg.Headers()["userId"]; !ok {
		return fmt.Errorf("handshake response missing userId")
	}

	log.Debugf("Handshake user ID: %s", userId)

	// Send a handshake message to the server
	err = conn.Write(network.NewMessage(map[string]string{
		"action": "authenticate",
		"user":   uc.user.Name,
		"userId": uc.user.UniqueId,
	}, nil))
	if err != nil {
		return err
	}
	// TODO: Implement the actual handshake acceptance logic

	return nil
}
