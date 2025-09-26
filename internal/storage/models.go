package storage

import (
	"time"

	"github.com/hop-/gotchat/internal/core"
)

type Model[T core.Entity] interface {
	ToEntity() *T
}

func getMappedFieldForType[T Model[E], E core.Entity](field string) string {
	var model T
	switch any(model).(type) {
	case *Message:
		return getMappedFieldForMessage(field)
	// Add more cases here for other models if needed
	default:
		return field
	}
}

// User model and dto
type User struct {
	Id       uint   `gorm:"primaryKey"`
	UniqueId string `gorm:"unique"`
	Name     string
}

func (u User) ToEntity() *core.User {
	return &core.User{
		BaseEntity: core.BaseEntity{Id: u.Id},
		UniqueId:   u.UniqueId,
		Name:       u.Name,
	}
}

func FromEntityToUser(entity *core.User) *User {
	return &User{
		Id:       entity.GetId(),
		UniqueId: entity.UniqueId,
		Name:     entity.Name,
	}
}

// Account model and dto
type Account struct {
	Id        uint `gorm:"primaryKey"`
	UserId    uint
	Password  string
	LastLogin time.Time
	User      User `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (a Account) ToEntity() *core.Account {
	return &core.Account{
		BaseEntity: core.BaseEntity{Id: a.Id},
		UserId:     a.UserId,
		Password:   a.Password,
		LastLogin:  a.LastLogin,
	}
}

func FromEntityToAccount(entity *core.Account) *Account {
	return &Account{
		Id:        entity.GetId(),
		UserId:    entity.UserId,
		Password:  entity.Password,
		LastLogin: entity.LastLogin,
	}
}

// ConnectionDetails model and dto
type ConnectionDetails struct {
	Id                uint   `gorm:"primaryKey"`
	HostUniqueId      string `gorm:"uniqueIndex:uni_connection_details_host_client"`
	ClientUniqueId    string `gorm:"uniqueIndex:uni_connection_details_host_client"`
	EncryptionKey     string
	DecryptionKey     string
	KeyDerivationSalt string
	CreatedAt         time.Time
}

func (c ConnectionDetails) ToEntity() *core.ConnectionDetails {
	return &core.ConnectionDetails{
		BaseEntity:        core.BaseEntity{Id: c.Id},
		HostUniqueId:      c.HostUniqueId,
		ClientUniqueId:    c.ClientUniqueId,
		EncryptionKey:     c.EncryptionKey,
		DecryptionKey:     c.DecryptionKey,
		KeyDerivationSalt: c.KeyDerivationSalt,
		CreatedAt:         c.CreatedAt,
	}
}

func FromEntityToConnectionDetails(entity *core.ConnectionDetails) *ConnectionDetails {
	return &ConnectionDetails{
		Id:                entity.GetId(),
		HostUniqueId:      entity.HostUniqueId,
		ClientUniqueId:    entity.ClientUniqueId,
		EncryptionKey:     entity.EncryptionKey,
		DecryptionKey:     entity.DecryptionKey,
		KeyDerivationSalt: entity.KeyDerivationSalt,
		CreatedAt:         entity.CreatedAt,
	}
}

// Channel model and dto
type Channel struct {
	Id       uint   `gorm:"primaryKey"`
	UniqueId string `gorm:"unique"`
	Name     string
}

func (ch Channel) ToEntity() *core.Channel {
	return &core.Channel{
		BaseEntity: core.BaseEntity{Id: ch.Id},
		UniqueId:   ch.UniqueId,
		Name:       ch.Name,
	}
}

func FromEntityToChannel(entity *core.Channel) *Channel {
	return &Channel{
		Id:       entity.GetId(),
		UniqueId: entity.UniqueId,
		Name:     entity.Name,
	}
}

// Attendance model and dto
type Attendance struct {
	Id        uint `gorm:"primaryKey"`
	UserId    uint `gorm:"uniqueIndex:uni_attendances_user_channel"`
	ChannelId uint `gorm:"uniqueIndex:uni_attendances_user_channel"`
	JoinedAt  time.Time
	User      User    `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Channel   Channel `gorm:"foreignKey:ChannelId;constraint:OnDelete:CASCADE"`
}

func (a Attendance) ToEntity() *core.Attendance {
	return &core.Attendance{
		BaseEntity: core.BaseEntity{Id: a.Id},
		UserId:     a.UserId,
		ChannelId:  a.ChannelId,
		JoinedAt:   a.JoinedAt,
	}
}

func FromEntityToAttendance(entity *core.Attendance) *Attendance {
	return &Attendance{
		Id:        entity.GetId(),
		UserId:    entity.UserId,
		ChannelId: entity.ChannelId,
		JoinedAt:  entity.JoinedAt,
	}
}

// Message model and dto
type Message struct {
	Id        uint `gorm:"primaryKey"`
	Content   string
	SenderId  uint      `gorm:"index:idx_messages_channel_user_created,priority:2"`
	ChannelId uint      `gorm:"index:idx_messages_channel_user_created,priority:1"`
	CreatedAt time.Time `gorm:"index:idx_messages_channel_user_created,priority:3"`
	Sender    User      `gorm:"foreignKey:SenderId;constraint:OnDelete:CASCADE"`
	Channel   Channel   `gorm:"foreignKey:ChannelId;constraint:OnDelete:CASCADE"`
}

func (m Message) ToEntity() *core.Message {
	return &core.Message{
		BaseEntity: core.BaseEntity{Id: m.Id},
		Content:    m.Content,
		UserId:     m.SenderId,
		ChannelId:  m.ChannelId,
		CreatedAt:  m.CreatedAt,
	}
}

func FromEntityToMessage(entity *core.Message) *Message {
	return &Message{
		Id:        entity.GetId(),
		Content:   entity.Content,
		SenderId:  entity.UserId,
		ChannelId: entity.ChannelId,
		CreatedAt: entity.CreatedAt,
	}
}

func getMappedFieldForMessage(field string) string {
	switch field {
	case "UserId":
		return "SenderId"
	default:
		return field
	}
}
