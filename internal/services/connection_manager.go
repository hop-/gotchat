package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/log"
	"github.com/hop-/gotchat/internal/network"
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

// User Controller
type UserController struct {
	AtomicRunningStatus

	// User associated with this controller
	user *core.User

	mu          sync.RWMutex
	connections map[string]network.AdvancedConn

	// Event emitter for user-related events
	eventEmitter core.EventEmitter

	// User manager
	UserManager *UserManager

	// Connection details manager
	connectionDetailsManager *ConnectionDetailsManager
}

func NewUserController(user *core.User, eventEmitter core.EventEmitter, userManager *UserManager, connectionDetailsManager *ConnectionDetailsManager) *UserController {
	return &UserController{
		AtomicRunningStatus{},
		user,
		sync.RWMutex{},
		make(map[string]network.AdvancedConn),
		eventEmitter,
		userManager,
		connectionDetailsManager,
	}
}

func (uc *UserController) Register(conn *network.Conn, isInitator bool) string {
	log.Debugf("Registering new connection for user %s", uc.user.Name)
	connId := uc.addConnection(conn)
	go uc.handleConnection(connId, conn, isInitator)

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
	uc.connections[id] = conn

	uc.emitEvent(NewConnection{id, conn})

	return id
}

func (uc *UserController) upgradeConnection(connId string, conn network.AdvancedConn) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if _, ok := uc.connections[connId]; ok {
		uc.connections[connId] = conn

		// TODO: Emit connection upgraded event
	}
}

func (uc *UserController) removeConnection(id string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	delete(uc.connections, id)

	uc.emitEvent(ConnectionClosed{id})
}

func (uc *UserController) handleConnection(connId string, conn *network.Conn, isInitiator bool) {
	// Ensure the connection is removed when done
	defer uc.removeConnection(connId)

	// Handshake
	secureConn, err := uc.handshake(connId, conn, isInitiator)
	if err != nil {
		// Close the original connection
		conn.Close()

		log.Errorf("Handshake failed for connection %s: %v", connId, err)
		uc.emitEvent(ConnectionFailed{err})

		return
	}

	// Ensure the connection is upgraded
	defer secureConn.Close()

	// Upgrade the connection
	uc.upgradeConnection(connId, secureConn)

	// Read messages from the secure connection
	for uc.isRunning() {
		m, err := secureConn.Read()
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

func (uc *UserController) handshake(connId string, conn *network.Conn, isInitiator bool) (network.AdvancedConn, error) {
	var clientUserId string
	var secureConn *network.SecureConn
	var err error

	if isInitiator {
		clientUserId, err = uc.initiateAuthentication(connId, conn)
		if err != nil {
			return conn, err
		}

		secureConn, err = uc.initiateAndHandleAuthenticationAndUpgrade(clientUserId, connId, conn)
		if err != nil {
			return conn, err
		}
	} else {
		clientUserId, err = uc.acceptAuthentication(connId, conn)
		if err != nil {
			return conn, err
		}

		secureConn, err = uc.acceptAndHandleAuthenticationAndUpgrade(clientUserId, connId, conn)
		if err != nil {
			return conn, err
		}
	}

	return secureConn, nil
}

func (uc *UserController) initiateAuthentication(connId string, conn *network.Conn) (string, error) {
	if !uc.isRunning() {
		log.Errorf("UserController is not running")

		return "", fmt.Errorf("user controller is not running")
	}
	if uc.user == nil {
		log.Errorf("User is not set")

		return "", fmt.Errorf("user is not set")
	}

	log.Infof("Initiating handshake for connection %s with user %s", connId, uc.user.Name)

	// Send a handshake message to the peer
	err := uc.sendHandshakeUserInfo(connId, conn)
	if err != nil {
		return "", err
	}

	// Receive the handshake response
	userId, err := uc.receiveHandshakeUserInfo(conn)
	if err != nil {
		return "", err
	}
	log.Debugf("Handshake user ID: %s", userId)

	return userId, nil
}

func (uc *UserController) acceptAuthentication(connId string, conn *network.Conn) (string, error) {
	if !uc.isRunning() {
		log.Errorf("UserController is not running")

		return "", fmt.Errorf("user controller is not running")
	}

	if uc.user == nil {
		log.Errorf("User is not set")

		return "", fmt.Errorf("user is not set")
	}

	log.Infof("Accepting handshake for connection %s", connId)

	// Receive the handshake response
	userId, err := uc.receiveHandshakeUserInfo(conn)
	if err != nil {
		return "", err
	}
	log.Debugf("Handshake user ID: %s", userId)

	// Send a handshake message to the peer
	err = uc.sendHandshakeUserInfo(connId, conn)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (uc *UserController) initiateAndHandleAuthenticationAndUpgrade(clientUserId string, connId string, conn *network.Conn) (*network.SecureConn, error) {
	// Get the connection details for the client user id
	connectionDetails, err := uc.connectionDetailsManager.GetConnectionDetails(uc.user.UniqueId, clientUserId)
	if err != nil {
		return nil, err
	}

	var connState string
	if connectionDetails == nil {
		connState = ConnectionStateUnknown
	} else {
		connState = ConnectionStateKnown
	}

	// Send state of the connection to the peer
	log.Debugf("Sending connection state %s to peer for connection %s", connState, connId)
	err = conn.Write(network.NewMessage(map[string]string{
		"action": "connection_state",
		"state":  connState,
	}, nil))
	if err != nil {
		return nil, err
	}

	// Read the response from the peer about the connection state
	msg, err := conn.Read()
	if err != nil {
		return nil, err
	}

	peerConnState := msg.Headers()["state"]

	if connState == ConnectionStateUnknown && peerConnState != ConnectionStateUnknown {
		return nil, fmt.Errorf("peer connection state is not unknown, expected %s, got %s", ConnectionStateUnknown, peerConnState)
	}

	if peerConnState == ConnectionStateUnknown {
		// Handle unknown connection
		encryptionKey, decryptionKey, err := uc.generateAndExchangeKeys(connId, conn)
		if err != nil {
			return nil, err
		}

		connectionDetails, err = uc.connectionDetailsManager.UpsertConnectionDetails(uc.user.UniqueId, clientUserId, encryptionKey, decryptionKey)
		if err != nil {
			return nil, err
		}
	}

	if connectionDetails == nil {
		return nil, fmt.Errorf("connection details not found for user %s and client %s", uc.user.UniqueId, clientUserId)
	}

	// Create the secure component for the connection
	secureComponent, err := network.NewEncryption(connectionDetails.EncryptionKey, connectionDetails.DecryptionKey)
	if err != nil {
		return nil, err
	}

	secureConn := network.NewSecureConn(*conn, secureComponent)

	// Generate a random phrase to send to the peer
	randomPhrase := generateRandomString()

	// Send the random phrase to the peer
	log.Debugf("Sending random phrase to peer with connection id %s: %s", connId, randomPhrase)
	err = secureConn.Write(network.NewMessage(map[string]string{
		"action": "send_phrase",
		"phrase": randomPhrase,
	}, nil))
	if err != nil {
		return nil, err
	}

	// Receive the echoed phrase from the peer
	log.Debugf("Waiting for echoed phrase from peer with connection id %s", connId)
	msg, err = secureConn.Read()
	if err != nil {
		return nil, err
	}

	echoedPhrase := msg.Headers()["phrase"]

	// Check if the echoed phrase matches the original random phrase
	if echoedPhrase != randomPhrase {
		log.Debugf("Echoed phrase does not match original phrase: %s != %s", echoedPhrase, randomPhrase)
		err := uc.connectionDetailsManager.RemoveConnectionDetails(uc.user.UniqueId, clientUserId)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("echoed phrase does not match original phrase")
	}

	log.Debugf("Echoed phrase matches original phrase: %s", echoedPhrase)

	return secureConn, nil
}

func (uc *UserController) acceptAndHandleAuthenticationAndUpgrade(clientUserId string, connId string, conn *network.Conn) (*network.SecureConn, error) {
	// Read the response from the peer about the connection state
	msg, err := conn.Read()
	if err != nil {
		return nil, err
	}

	// Get the connection details for the client user id
	connectionDetails, err := uc.connectionDetailsManager.GetConnectionDetails(uc.user.UniqueId, clientUserId)
	if err != nil {
		return nil, err
	}

	peerConnState := msg.Headers()["state"]
	var connState string

	if peerConnState == ConnectionStateUnknown {
		connState = ConnectionStateUnknown
	} else {
		if connectionDetails == nil {
			connState = ConnectionStateUnknown
		} else {
			connState = ConnectionStateKnown
		}
	}

	// Send the connection state to the peer
	log.Debugf("Sending connection state %s to peer for connection %s", connState, connId)
	err = conn.Write(network.NewMessage(map[string]string{
		"action": "connection_state",
		"state":  connState,
	}, nil))
	if err != nil {
		return nil, err
	}

	if connState == ConnectionStateUnknown {
		// handle unknown connection
		encryptionKey, decryptionKey, err := uc.generateAndExchangeKeys(connId, conn)
		if err != nil {
			return nil, err
		}

		connectionDetails, err = uc.connectionDetailsManager.UpsertConnectionDetails(uc.user.UniqueId, clientUserId, encryptionKey, decryptionKey)
		if err != nil {
			return nil, err
		}
	}

	if connectionDetails == nil {
		return nil, fmt.Errorf("connection details not found for user %s and client %s", uc.user.UniqueId, clientUserId)
	}

	// Create the secure component for the connection
	secureComponent, err := network.NewEncryption(connectionDetails.EncryptionKey, connectionDetails.DecryptionKey)
	if err != nil {
		return nil, err
	}

	secureConn := network.NewSecureConn(*conn, secureComponent)

	// Reading first phrase from secure connection
	log.Debugf("Waiting for phrase from peer with connection id %s", connId)
	msg, err = secureConn.Read()
	if err != nil {
		return nil, err
	}

	helloPhrase := msg.Headers()["phrase"]

	// Sending the phrase back to the peer
	log.Debugf("Sending echoed phrase back to peer with connection id %s: %s", connId, helloPhrase)
	err = secureConn.Write(network.NewMessage(map[string]string{
		"action": "echo_phrase",
		"phrase": helloPhrase,
	}, nil))
	if err != nil {
		return nil, err
	}

	log.Debugf("Echoed phrase back to peer: %s", helloPhrase)

	return secureConn, nil
}

func (uc *UserController) generateAndExchangeKeys(connId string, conn *network.Conn) ([]byte, []byte, error) {
	log.Debugf("Generating and exchanging keys for connection %s", connId)
	// Generate a new encryption key and decryption key
	encryptionKey, err := network.GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	// base64 the encryption key
	passphrase := base64.StdEncoding.EncodeToString(encryptionKey)

	// Send the keys to the peer
	log.Debugf("Sending encryption key to peer for connection %s", connId)
	err = conn.Write(network.NewMessage(map[string]string{
		"action":     "exchange_keys",
		"passphrase": passphrase,
	}, nil))
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("Waiting for decryption key from peer for connection %s", connId)
	msg, err := conn.Read()
	if err != nil {
		return nil, nil, err
	}

	// Get the decryption key from the response
	passphrase = msg.Headers()["passphrase"]

	decryptionKey, err := base64.StdEncoding.DecodeString(passphrase)
	if err != nil {
		return nil, nil, err
	}

	if decryptionKey == nil {
		return nil, nil, fmt.Errorf("decryption key is nil")
	}

	log.Debugf("Key exchange complete for connection %s", connId)

	return encryptionKey, decryptionKey, nil
}

func (uc *UserController) sendHandshakeUserInfo(connId string, conn *network.Conn) error {
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

	return nil
}

func (uc *UserController) receiveHandshakeUserInfo(conn *network.Conn) (string, error) {
	msg, err := conn.Read()
	if err != nil {
		if network.IsClosedError(err) {
			log.Infof("Connection closed by peer before the handshake")
		}

		return "", err
	}

	// Checking message
	if action, ok := msg.Headers()["action"]; !ok || action != "authenticate" {
		return "", fmt.Errorf("handshake response missing or invalid action: %s", action)
	}

	var userId string
	var ok bool
	if userId, ok = msg.Headers()["userId"]; !ok {
		return "", fmt.Errorf("handshake response missing userId")
	}

	// Validate user ID
	if userId == uc.user.UniqueId {
		return "", fmt.Errorf("handshake user ID matches the current user: %s", userId)
	}

	return userId, nil
}
