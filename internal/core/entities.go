package core

import (
	"reflect"
	"sync"
	"time"
)

// Entity fieldName cache
var (
	mutex            sync.Mutex
	entityFieldNames = make(map[reflect.Type]map[string]bool)
)

type Entity interface {
	GetId() uint
}

func ListFieldNames[T Entity]() map[string]bool {
	mutex.Lock()
	defer mutex.Unlock()

	// Check cache
	var zero T
	t := reflect.TypeOf(zero)
	if fields, ok := entityFieldNames[t]; ok {
		// TODO: concider to return a copy to avoid external modification
		return fields
	}

	// Not in cache, compute and store
	var fields = make(map[string]bool)

	// Recursively walk through embedded structs to get all fields
	var walk func(reflect.Type)
	walk = func(tt reflect.Type) {
		for i := 0; i < tt.NumField(); i++ {
			f := tt.Field(i)
			if f.Anonymous {
				// recurse into embedded struct
				walk(f.Type)
			} else {
				fields[f.Name] = true
			}
		}
	}
	walk(t)

	// Cache result
	entityFieldNames[t] = fields

	// TODO: concider to return a copy to avoid external modification
	return fields
}

func IsFieldExist[T Entity](field string) bool {
	_, ok := ListFieldNames[T]()[field]

	return ok
}

type BaseEntity struct {
	Id uint
}

func (e BaseEntity) GetId() uint {
	return e.Id
}

type Account struct {
	BaseEntity
	UserId    uint
	Password  string
	LastLogin time.Time
}

func NewAccount(userId uint, password string, lastLogin time.Time) *Account {
	return &Account{
		BaseEntity: BaseEntity{},
		UserId:     userId,
		Password:   password,
		LastLogin:  lastLogin,
	}
}

// User entity
type User struct {
	BaseEntity
	UniqueId string `dbname:"unique_id"`
	Name     string `dbname:"name"`
}

func NewUser(name string) *User {
	return &User{
		BaseEntity: BaseEntity{},
		UniqueId:   generateUuid(),
		Name:       name,
	}
}

// ConnectionDetails entity
type ConnectionDetails struct {
	BaseEntity
	HostUniqueId      string    `dbname:"host_unique_id"`
	ClientUniqueId    string    `dbname:"client_unique_id"`
	EncryptionKey     string    `dbname:"encryption_key"`
	DecryptionKey     string    `dbname:"decryption_key"`
	KeyDerivationSalt string    `dbname:"key_derivation_salt"`
	CreatedAt         time.Time `dbname:"created_at"`
}

func NewConnectionDetails(hostUniqueId string, clientUniqueId string, encryptionKey string, decryptionKey string, keyDerivationSalt string) *ConnectionDetails {
	return &ConnectionDetails{
		BaseEntity:        BaseEntity{},
		HostUniqueId:      hostUniqueId,
		ClientUniqueId:    clientUniqueId,
		EncryptionKey:     encryptionKey,
		DecryptionKey:     decryptionKey,
		KeyDerivationSalt: keyDerivationSalt,
		CreatedAt:         time.Now(),
	}
}

// Message entity
type Message struct {
	BaseEntity
	UserId    uint      `dbname:"user_id"`
	ChannelId uint      `dbname:"channel_id"`
	Content   string    `dbname:"text"`
	CreatedAt time.Time `dbname:"created_at"`
}

func NewMessage(userId uint, channelId uint, text string) *Message {
	return &Message{
		BaseEntity: BaseEntity{},
		UserId:     userId,
		ChannelId:  channelId,
		Content:    text,
		CreatedAt:  time.Now(),
	}
}

// Channel entity
type Channel struct {
	BaseEntity
	UniqueId string `dbname:"unique_id"`
	Name     string `dbname:"name"`
}

func NewChannel(name string) *Channel {
	return &Channel{
		BaseEntity: BaseEntity{},
		UniqueId:   generateUuid(),
		Name:       name,
	}
}

// Attendance entity
type Attendance struct {
	BaseEntity
	UserId    uint      `dbname:"user_id"`
	ChannelId uint      `dbname:"channel_id"`
	JoinedAt  time.Time `dbname:"joined_at"`
}

func NewAttendance(userId uint, channelId uint) *Attendance {
	return &Attendance{
		BaseEntity: BaseEntity{},
		UserId:     userId,
		ChannelId:  channelId,
		JoinedAt:   time.Now(),
	}
}
