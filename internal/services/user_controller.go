package services

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/pkg/log"
	"github.com/hop-/gotchat/pkg/network"
)

type ConnectionInfo struct {
	Conn          network.AdvancedConn
	Authenticated bool
	peerUserId    string
}

// User Controller
type UserController struct {
	AtomicRunningStatus

	mu              sync.RWMutex
	connectionInfos map[string]*ConnectionInfo

	// User associated with this controller
	user *core.User

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
		sync.RWMutex{},
		make(map[string]*ConnectionInfo),
		user,
		eventEmitter,
		userManager,
		connectionDetailsManager,
	}
}

func (uc *UserController) Register(conn *network.Conn, isInitator bool) string {
	log.Debugf("Registering new connection for user %s", uc.user.Name)
	connId := uc.addUnauthenticatiedConnection(conn)
	go uc.handleConnection(connId, conn, isInitator)

	return connId
}

func (uc *UserController) Close() error {
	uc.setRunningStatus(false)

	for _, connInfo := range uc.connectionInfos {
		if err := connInfo.Conn.Close(); err != nil {
			log.Errorf("Failed to close connection: %v", err)
		}
	}

	return nil
}

func (uc *UserController) emitEvent(event core.Event) {
	uc.eventEmitter.Emit(event)
}

func (uc *UserController) addUnauthenticatiedConnection(conn *network.Conn) string {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	id := generateUuid()
	uc.connectionInfos[id] = &ConnectionInfo{conn, false, ""}

	uc.emitEvent(NewUnauthenticatedConnection{id, conn})

	return id
}

func (uc *UserController) upgradeConnection(connId string, conn network.AdvancedConn, peerUserId string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if connInfo, ok := uc.connectionInfos[connId]; ok {
		connInfo.Conn = conn
		connInfo.Authenticated = true
		connInfo.peerUserId = peerUserId
	}

	// Emit connection established event
	uc.emitEvent(ConnectionEstablished{connId, conn, peerUserId})
}

func (uc *UserController) removeConnection(id string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	delete(uc.connectionInfos, id)

	uc.emitEvent(ConnectionClosed{id})
}

func (uc *UserController) handleConnection(connId string, conn *network.Conn, isInitiator bool) {
	// Ensure the connection is removed when done
	defer uc.removeConnection(connId)

	// Handshake
	secureConn, peerUserId, err := uc.handshake(connId, conn, isInitiator)
	if err != nil {
		// Close the original connection
		conn.Close()

		log.Errorf("Handshake failed for connection %s: %v", connId, err)
		uc.emitEvent(ConnectionFailed{err})

		return
	}

	log.Infof("Handshake was successful, connection accepted %s", connId)

	// Ensure the connection is upgraded
	defer secureConn.Close()

	// Upgrade the connection
	uc.upgradeConnection(connId, secureConn, peerUserId)

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

func (uc *UserController) handshake(connId string, conn *network.Conn, isInitiator bool) (network.AdvancedConn, string, error) {
	var clientUserId string
	var secureConn *network.SecureConn
	var peerUserId string
	var err error

	if isInitiator {
		clientUserId, err = uc.initiateAuthentication(connId, conn)
		if err != nil {
			return conn, "", err
		}

		secureConn, peerUserId, err = uc.initiateAndHandleAuthenticationAndUpgrade(clientUserId, connId, conn)
		if err != nil {
			return conn, "", err
		}
	} else {
		clientUserId, err = uc.acceptAuthentication(connId, conn)
		if err != nil {
			return conn, "", err
		}

		secureConn, peerUserId, err = uc.acceptAndHandleAuthenticationAndUpgrade(clientUserId, connId, conn)
		if err != nil {
			return conn, "", err
		}
	}

	return secureConn, peerUserId, nil
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

func (uc *UserController) initiateAndHandleAuthenticationAndUpgrade(clientUserId string, connId string, conn *network.Conn) (*network.SecureConn, string, error) {
	// Get the connection details for the client user id
	connectionDetails, err := uc.connectionDetailsManager.GetConnectionDetails(uc.user.UniqueId, clientUserId)
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	// Read the response from the peer about the connection state
	msg, err := conn.Read()
	if err != nil {
		return nil, "", err
	}

	peerConnState := msg.Headers()["state"]

	if connState == ConnectionStateUnknown && peerConnState != ConnectionStateUnknown {
		return nil, "", fmt.Errorf("peer connection state is not unknown, expected %s, got %s", ConnectionStateUnknown, peerConnState)
	}

	if peerConnState == ConnectionStateUnknown {
		// Handle unknown connection
		encryptionKey, decryptionKey, err := uc.generateAndExchangeKeys(connId, conn)
		if err != nil {
			return nil, "", err
		}

		connectionDetails, err = uc.connectionDetailsManager.UpsertConnectionDetails(uc.user.UniqueId, clientUserId, encryptionKey, decryptionKey)
		if err != nil {
			return nil, "", err
		}
	}

	if connectionDetails == nil {
		return nil, "", fmt.Errorf("connection details not found for user %s and client %s", uc.user.UniqueId, clientUserId)
	}

	// Create the secure component for the connection
	secureComponent, err := network.NewEncryption(connectionDetails.EncryptionKey, connectionDetails.DecryptionKey)
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	// Receive the echoed phrase from the peer
	log.Debugf("Waiting for echoed phrase from peer with connection id %s", connId)
	msg, err = secureConn.Read()
	if err != nil {
		return nil, "", err
	}

	echoedPhrase := msg.Headers()["phrase"]

	// Check if the echoed phrase matches the original random phrase
	if echoedPhrase != randomPhrase {
		log.Debugf("Echoed phrase does not match original phrase: %s != %s", echoedPhrase, randomPhrase)
		err := uc.connectionDetailsManager.RemoveConnectionDetails(uc.user.UniqueId, clientUserId)
		if err != nil {
			return nil, "", err
		}

		return nil, "", fmt.Errorf("echoed phrase does not match original phrase")
	}

	log.Debugf("Echoed phrase matches original phrase: %s", echoedPhrase)

	return secureConn, clientUserId, nil
}

func (uc *UserController) acceptAndHandleAuthenticationAndUpgrade(clientUserId string, connId string, conn *network.Conn) (*network.SecureConn, string, error) {
	// Read the response from the peer about the connection state
	msg, err := conn.Read()
	if err != nil {
		return nil, "", err
	}

	// Get the connection details for the client user id
	connectionDetails, err := uc.connectionDetailsManager.GetConnectionDetails(uc.user.UniqueId, clientUserId)
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	if connState == ConnectionStateUnknown {
		// handle unknown connection
		encryptionKey, decryptionKey, err := uc.generateAndExchangeKeys(connId, conn)
		if err != nil {
			return nil, "", err
		}

		connectionDetails, err = uc.connectionDetailsManager.UpsertConnectionDetails(uc.user.UniqueId, clientUserId, encryptionKey, decryptionKey)
		if err != nil {
			return nil, "", err
		}
	}

	if connectionDetails == nil {
		return nil, "", fmt.Errorf("connection details not found for user %s and client %s", uc.user.UniqueId, clientUserId)
	}

	// Create the secure component for the connection
	secureComponent, err := network.NewEncryption(connectionDetails.EncryptionKey, connectionDetails.DecryptionKey)
	if err != nil {
		return nil, "", err
	}

	secureConn := network.NewSecureConn(*conn, secureComponent)

	// Reading first phrase from secure connection
	log.Debugf("Waiting for phrase from peer with connection id %s", connId)
	msg, err = secureConn.Read()
	if err != nil {
		return nil, "", err
	}

	helloPhrase := msg.Headers()["phrase"]

	// Sending the phrase back to the peer
	log.Debugf("Sending echoed phrase back to peer with connection id %s: %s", connId, helloPhrase)
	err = secureConn.Write(network.NewMessage(map[string]string{
		"action": "echo_phrase",
		"phrase": helloPhrase,
	}, nil))
	if err != nil {
		return nil, "", err
	}

	log.Debugf("Echoed phrase back to peer: %s", helloPhrase)

	return secureConn, clientUserId, nil
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
