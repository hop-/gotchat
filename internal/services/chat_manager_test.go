package services

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hop-/gotchat/internal/core"
	"github.com/stretchr/testify/mock"
)

func TestNewChatManager(t *testing.T) {
	userManager := &UserManager{}
	channelRepo := core.NewMockRepository[core.Channel](t)
	attendanceRepo := core.NewMockRepository[core.Attendance](t)
	messageRepo := core.NewMockRepository[core.Message](t)

	cm := NewChatManager(userManager, channelRepo, attendanceRepo, messageRepo)

	if cm == nil {
		t.Error("Expected ChatManager to be created")
		return
	}
	if cm.userManager != userManager {
		t.Error("Expected userManager to be set")
	}
	if cm.channelRepo != channelRepo {
		t.Error("Expected channelRepo to be set")
	}
	if cm.attendanceRepo != attendanceRepo {
		t.Error("Expected attendanceRepo to be set")
	}
	if cm.messageRepo != messageRepo {
		t.Error("Expected messageRepo to be set")
	}
}

func TestChatManager_Name(t *testing.T) {
	cm := &ChatManager{}
	if cm.Name() != "ChatManager" {
		t.Errorf("Expected Name() to return 'ChatManager', got %s", cm.Name())
	}
}

func TestChatManager_Init(t *testing.T) {
	cm := &ChatManager{}
	if err := cm.Init(); err != nil {
		t.Errorf("Expected Init() to return nil, got %v", err)
	}
}

func TestChatManager_Close(t *testing.T) {
	cm := &ChatManager{}
	if err := cm.Close(); err != nil {
		t.Errorf("Expected Close() to return nil, got %v", err)
	}
}

func TestChatManager_Run(t *testing.T) {
	cm := &ChatManager{}
	ctx := context.Background()
	var wg sync.WaitGroup

	// This should not panic or block
	cm.Run(ctx, &wg)
}

func TestChatManager_MapEventToCommands(t *testing.T) {
	cm := &ChatManager{}
	mockEvent := core.NewMockEvent(t)

	commands := cm.MapEventToCommands(mockEvent)
	if commands != nil {
		t.Errorf("Expected MapEventToCommands to return nil, got %v", commands)
	}
}

func TestChatManager_GetChatsByUserId_Success(t *testing.T) {
	userManager := &UserManager{}
	channelRepo := core.NewMockRepository[core.Channel](t)
	attendanceRepo := core.NewMockRepository[core.Attendance](t)
	messageRepo := core.NewMockRepository[core.Message](t)

	cm := NewChatManager(userManager, channelRepo, attendanceRepo, messageRepo)

	userId := 1
	attendances := []*core.Attendance{
		{
			BaseEntity: core.BaseEntity{Id: 1},
			UserId:     userId,
			ChannelId:  10,
			JoinedAt:   time.Now(),
		},
		{
			BaseEntity: core.BaseEntity{Id: 2},
			UserId:     userId,
			ChannelId:  20,
			JoinedAt:   time.Now(),
		},
	}

	channels := []*core.Channel{
		{
			BaseEntity: core.BaseEntity{Id: 10},
			UniqueId:   "channel-1",
			Name:       "General",
		},
		{
			BaseEntity: core.BaseEntity{Id: 20},
			UniqueId:   "channel-2",
			Name:       "Random",
		},
	}

	// Setup mock expectations
	attendanceRepo.On("GetAllBy", "user_id", userId).Return(attendances, nil)
	channelRepo.On("GetOneBy", "unique_id", mock.Anything).Return(channels[0], nil).Once()
	channelRepo.On("GetOneBy", "unique_id", mock.Anything).Return(channels[1], nil).Once()

	chats, err := cm.GetChatsByUserId(userId)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(chats) != 2 {
		t.Errorf("Expected 2 chats, got %d", len(chats))
	}
	if chats[0].Id != "channel-1" {
		t.Errorf("Expected first chat ID to be 'channel-1', got %s", chats[0].Id)
	}
	if chats[0].Name != "General" {
		t.Errorf("Expected first chat name to be 'General', got %s", chats[0].Name)
	}
	if chats[1].Id != "channel-2" {
		t.Errorf("Expected second chat ID to be 'channel-2', got %s", chats[1].Id)
	}
	if chats[1].Name != "Random" {
		t.Errorf("Expected second chat name to be 'Random', got %s", chats[1].Name)
	}
}

func TestChatManager_GetChatsByUserId_AttendanceError(t *testing.T) {
	userManager := &UserManager{}
	channelRepo := core.NewMockRepository[core.Channel](t)
	attendanceRepo := core.NewMockRepository[core.Attendance](t)
	messageRepo := core.NewMockRepository[core.Message](t)

	cm := NewChatManager(userManager, channelRepo, attendanceRepo, messageRepo)

	userId := 1
	expectedError := fmt.Errorf("database error")

	// Setup mock expectations
	attendanceRepo.On("GetAllBy", "user_id", userId).Return(nil, expectedError)

	chats, err := cm.GetChatsByUserId(userId)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if chats != nil {
		t.Errorf("Expected chats to be nil, got %v", chats)
	}
	if err.Error() != "failed to get attendances: database error" {
		t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

func TestChatManager_GetChatsByUserId_ChannelError(t *testing.T) {
	userManager := &UserManager{}
	channelRepo := core.NewMockRepository[core.Channel](t)
	attendanceRepo := core.NewMockRepository[core.Attendance](t)
	messageRepo := core.NewMockRepository[core.Message](t)

	cm := NewChatManager(userManager, channelRepo, attendanceRepo, messageRepo)

	userId := 1
	attendances := []*core.Attendance{
		{
			BaseEntity: core.BaseEntity{Id: 1},
			UserId:     userId,
			ChannelId:  10,
			JoinedAt:   time.Now(),
		},
	}
	expectedError := fmt.Errorf("channel not found")

	// Setup mock expectations
	attendanceRepo.On("GetAllBy", "user_id", userId).Return(attendances, nil)
	channelRepo.On("GetOneBy", "unique_id", mock.Anything).Return(nil, expectedError)

	chats, err := cm.GetChatsByUserId(userId)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if chats != nil {
		t.Errorf("Expected chats to be nil, got %v", chats)
	}
	if err.Error() != "failed to get channel: channel not found" {
		t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

func TestChatManager_GetChatMessagesByChatId_Success(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)
	userManager := NewUserManager(eventEmitter, userRepo)

	channelRepo := core.NewMockRepository[core.Channel](t)
	attendanceRepo := core.NewMockRepository[core.Attendance](t)
	messageRepo := core.NewMockRepository[core.Message](t)

	cm := NewChatManager(userManager, channelRepo, attendanceRepo, messageRepo)

	chatId := "channel-1"
	channel := &core.Channel{
		BaseEntity: core.BaseEntity{Id: 10},
		UniqueId:   chatId,
		Name:       "General",
	}

	messages := []*core.Message{
		{
			BaseEntity: core.BaseEntity{Id: 1},
			UserId:     1,
			ChannelId:  10,
			Text:       "Hello world!",
			CreatedAt:  time.Now().Add(-2 * time.Hour),
		},
		{
			BaseEntity: core.BaseEntity{Id: 2},
			UserId:     2,
			ChannelId:  10,
			Text:       "How are you?",
			CreatedAt:  time.Now().Add(-1 * time.Hour),
		},
	}

	users := []*core.User{
		{
			BaseEntity: core.BaseEntity{Id: 1},
			Name:       "Alice",
		},
		{
			BaseEntity: core.BaseEntity{Id: 2},
			Name:       "Bob",
		},
	}

	// Setup mock expectations
	channelRepo.On("GetOneBy", "unique_id", chatId).Return(channel, nil)
	messageRepo.On("GetAllBy", "channel_id", 10).Return(messages, nil)
	userRepo.On("GetOne", 1).Return(users[0], nil)
	userRepo.On("GetOne", 2).Return(users[1], nil)

	chatMessages, err := cm.GetChatMessagesByChatId(chatId)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(chatMessages) != 2 {
		t.Errorf("Expected 2 chat messages, got %d", len(chatMessages))
	}
	if chatMessages[0].Member != "Alice" {
		t.Errorf("Expected first message member to be 'Alice', got %s", chatMessages[0].Member)
	}
	if chatMessages[0].Text != "Hello world!" {
		t.Errorf("Expected first message text to be 'Hello world!', got %s", chatMessages[0].Text)
	}
	if chatMessages[1].Member != "Bob" {
		t.Errorf("Expected second message member to be 'Bob', got %s", chatMessages[1].Member)
	}
	if chatMessages[1].Text != "How are you?" {
		t.Errorf("Expected second message text to be 'How are you?', got %s", chatMessages[1].Text)
	}
}

func TestChatManager_GetChatMessagesByChatId_ChannelError(t *testing.T) {
	userManager := &UserManager{}
	channelRepo := core.NewMockRepository[core.Channel](t)
	attendanceRepo := core.NewMockRepository[core.Attendance](t)
	messageRepo := core.NewMockRepository[core.Message](t)

	cm := NewChatManager(userManager, channelRepo, attendanceRepo, messageRepo)

	chatId := "non-existent"
	expectedError := fmt.Errorf("channel not found")

	// Setup mock expectations
	channelRepo.On("GetOneBy", "unique_id", chatId).Return(nil, expectedError)

	chatMessages, err := cm.GetChatMessagesByChatId(chatId)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if chatMessages != nil {
		t.Errorf("Expected chat messages to be nil, got %v", chatMessages)
	}
	if err.Error() != "failed to get channel: channel not found" {
		t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

func TestChatManager_GetChatMessagesByChatId_MessageError(t *testing.T) {
	userManager := &UserManager{}
	channelRepo := core.NewMockRepository[core.Channel](t)
	attendanceRepo := core.NewMockRepository[core.Attendance](t)
	messageRepo := core.NewMockRepository[core.Message](t)

	cm := NewChatManager(userManager, channelRepo, attendanceRepo, messageRepo)

	chatId := "channel-1"
	channel := &core.Channel{
		BaseEntity: core.BaseEntity{Id: 10},
		UniqueId:   chatId,
		Name:       "General",
	}
	expectedError := fmt.Errorf("database error")

	// Setup mock expectations
	channelRepo.On("GetOneBy", "unique_id", chatId).Return(channel, nil)
	messageRepo.On("GetAllBy", "channel_id", 10).Return(nil, expectedError)

	chatMessages, err := cm.GetChatMessagesByChatId(chatId)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if chatMessages != nil {
		t.Errorf("Expected chat messages to be nil, got %v", chatMessages)
	}
	if err.Error() != "failed to get messages: database error" {
		t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

func TestChatManager_GetChatMessagesByChatId_UserError(t *testing.T) {
	eventEmitter := core.NewMockEventEmitter(t)
	userRepo := core.NewMockRepository[core.User](t)
	userManager := NewUserManager(eventEmitter, userRepo)

	channelRepo := core.NewMockRepository[core.Channel](t)
	attendanceRepo := core.NewMockRepository[core.Attendance](t)
	messageRepo := core.NewMockRepository[core.Message](t)

	cm := NewChatManager(userManager, channelRepo, attendanceRepo, messageRepo)

	chatId := "channel-1"
	channel := &core.Channel{
		BaseEntity: core.BaseEntity{Id: 10},
		UniqueId:   chatId,
		Name:       "General",
	}

	messages := []*core.Message{
		{
			BaseEntity: core.BaseEntity{Id: 1},
			UserId:     999, // Non-existent user
			ChannelId:  10,
			Text:       "Hello world!",
			CreatedAt:  time.Now(),
		},
	}

	// Setup mock expectations
	channelRepo.On("GetOneBy", "unique_id", chatId).Return(channel, nil)
	messageRepo.On("GetAllBy", "channel_id", 10).Return(messages, nil)
	userRepo.On("GetOne", 999).Return((*core.User)(nil), nil)

	chatMessages, err := cm.GetChatMessagesByChatId(chatId)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if chatMessages != nil {
		t.Errorf("Expected chat messages to be nil, got %v", chatMessages)
	}
	if err.Error() != "failed to get user: entity not found" {
		t.Errorf("Expected specific error message, got %s", err.Error())
	}
}
