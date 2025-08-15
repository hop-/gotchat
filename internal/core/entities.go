package core

import (
	"reflect"
	"time"
)

type Entity interface {
	GetId() int
}

func GetFieldNamesOfEntity[T Entity]() []string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	fieldNames := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if fieldName, ok := field.Tag.Lookup("name"); ok {
			fieldNames = append(fieldNames, fieldName)
		}
	}

	return fieldNames
}

type BaseEntity struct {
	Id int
}

func (e BaseEntity) GetId() int {
	return e.Id
}

// User entity
type User struct {
	BaseEntity
	UniqueId  string    `name:"unique_id"`
	Name      string    `name:"name"`
	Password  string    `name:"password"`
	LastLogin time.Time `name:"last_login"`
}

func NewUser(name string, password string) *User {
	return &User{
		BaseEntity: BaseEntity{},
		UniqueId:   generateUuid(),
		Name:       name,
		Password:   password,
		LastLogin:  time.Now(),
	}
}

// ConnectionDetails entity
type ConnectionDetails struct {
	BaseEntity
	HostUniqueId      string    `name:"host_unique_id"`
	ClientUniqueId    string    `name:"client_unique_id"`
	EncryptionKey     string    `name:"encryption_key"`
	DecryptionKey     string    `name:"decryption_key"`
	KeyDerivationSalt string    `name:"key_derivation_salt"`
	CreatedAt         time.Time `name:"created_at"`
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
	UserId    int       `name:"user_id"`
	ChannelId int       `name:"channel_id"`
	Text      string    `name:"text"`
	CreatedAt time.Time `name:"created_at"`
}

func NewMessage(userId int, channelId int, text string) *Message {
	return &Message{
		BaseEntity: BaseEntity{},
		UserId:     userId,
		ChannelId:  channelId,
		Text:       text,
		CreatedAt:  time.Now(),
	}
}

// Channel entity
type Channel struct {
	BaseEntity
	UniqueId string `name:"unique_id"`
	Name     string `name:"name"`
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
	UserId    int       `name:"user_id"`
	ChannelId int       `name:"channel_id"`
	JoinedAt  time.Time `name:"joined_at"`
}

func NewAttendance(userId int, channelId int) *Attendance {
	return &Attendance{
		BaseEntity: BaseEntity{},
		UserId:     userId,
		ChannelId:  channelId,
		JoinedAt:   time.Now(),
	}
}
