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
		fieldName := field.Tag.Get("name")
		fieldNames = append(fieldNames, fieldName)
	}

	return fieldNames
}

type BaseEntity struct {
	Id int
}

func (e BaseEntity) GetId() int {
	return e.Id
}

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

type Message struct {
	BaseEntity
	UserId    string    `name:"user_id"`
	ChannelId string    `name:"channel_id"`
	Text      string    `name:"text"`
	CreatedAt time.Time `name:"created_at"`
}

func NewMessage(userId, channelId, text string) *Message {
	return &Message{
		BaseEntity: BaseEntity{},
		UserId:     userId,
		ChannelId:  channelId,
		Text:       text,
		CreatedAt:  time.Now(),
	}
}

type Channel struct {
	BaseEntity
	UniqueId string `name:"unique_id"`
	Name     string `name:"name"`
}

func NewChannel(name string) *Channel {
	return &Channel{
		BaseEntity: BaseEntity{},
		Name:       name,
	}
}

type Attendance struct {
	BaseEntity
	UserId    string    `name:"user_id"`
	ChannelId string    `name:"channel_id"`
	JoinedAt  time.Time `name:"joined_at"`
}
