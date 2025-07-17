package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hop-/gotchat/internal/core"
)

type Chat struct {
	Id   string
	Name string
}

type ChatMessage struct {
	Member string
	Text   string
	At     time.Time
}

type ChatManager struct {
	// Services
	userManager *UserManager

	// Repos
	channelRepo    core.Repository[core.Channel]
	attendanceRepo core.Repository[core.Attendance]
	messageRepo    core.Repository[core.Message]
}

func NewChatManager(
	userManager *UserManager,
	channelRepo core.Repository[core.Channel],
	attendanceRepo core.Repository[core.Attendance],
	messageRepo core.Repository[core.Message],
) *ChatManager {
	return &ChatManager{
		userManager,
		channelRepo,
		attendanceRepo,
		messageRepo,
	}
}

// Init implements core.Service.
func (cm *ChatManager) Init() error {
	return nil
}

// Name implements core.Service.
func (cm *ChatManager) Name() string {
	return "ChatManager"
}

// Run implements core.Service.
func (cm *ChatManager) Run(ctx context.Context, wg *sync.WaitGroup) {
}

// Close implements core.Service.
func (cm *ChatManager) Close() error {
	return nil
}

func (cm *ChatManager) MapEventToCommands(event core.Event) []core.Command {
	// TODO
	return nil
}

func (cm *ChatManager) GetChatsByUserId(userId int) ([]Chat, error) {
	attendatnces, err := cm.attendanceRepo.GetAllBy("user_id", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %s", err.Error())
	}

	chats := make([]Chat, 0, len(attendatnces))
	for _, attendance := range attendatnces {
		channel, err := cm.channelRepo.GetOneBy("unique_id", attendance.ChannelId)
		if err != nil {
			return nil, fmt.Errorf("failed to get channel: %s", err.Error())
		}

		chats = append(chats, Chat{channel.UniqueId, channel.Name})
	}

	return chats, nil
}

func (cm *ChatManager) GetChatMessagesByChatId(chatId string) ([]ChatMessage, error) {
	chat, err := cm.channelRepo.GetOneBy("unique_id", chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %s", err.Error())
	}

	messages, err := cm.messageRepo.GetAllBy("channel_id", chat.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %s", err.Error())
	}

	chatMessages := make([]ChatMessage, 0, len(messages))

	for _, message := range messages {
		user, err := cm.userManager.GetUserById(message.UserId)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %s", err.Error())
		}

		chatMessages = append(chatMessages, ChatMessage{
			Member: user.Name,
			Text:   message.Text,
			At:     message.CreatedAt,
		})
	}

	return chatMessages, nil
}
