package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hop-/gotchat/internal/core"
)

type Chat struct {
	Name string
	Id   string
}

// FilterValue implements list.Item.
func (c Chat) FilterValue() string {
	return c.Name
}

type ChatViewModel struct {
	// Frame component
	Frame

	// Focusable component
	*FocusContainer

	// Components
	chats       *ItemList
	chatHistory *ChatHistory
	chatInput   *ChatInput

	// Stack
	stack *Stack

	// Repos
	userRepo       core.Repository[core.User]
	channelRepo    core.Repository[core.Channel]
	attendanceRepo core.Repository[core.Attendance]
	messageRepo    core.Repository[core.Message]

	// User entity
	user *core.User
}

func newChatViewModel(
	user *core.User,
	userRepo core.Repository[core.User],
	channelRepo core.Repository[core.Channel],
	attendanceRepo core.Repository[core.Attendance],
	messageRepo core.Repository[core.Message],
) *ChatViewModel {
	// Initialize chat list
	chats := newItemList([]list.Item{})
	chats.Title = "Chats"

	// Initialize chat history
	chatHistory := newChatHistory("You", "TheOtherOne")

	// Initialize chat input
	chatInput := newChatInput()
	chatInput.Placeholder = "Type a message..."
	chatInput.SetActive(true)
	chatInput.SetWidth(20)

	return &ChatViewModel{
		Frame{},
		&FocusContainer{[]FocusableModel{chatInput, chats, chatHistory}, 0},

		chats,
		chatHistory,
		chatInput,
		newStack(Horizontal, 3, chats, newStackWithPosition(lipgloss.Left, Vertical, 2, chatHistory, chatInput)),
		userRepo,
		channelRepo,
		attendanceRepo,
		messageRepo,
		user,
	}
}

func (m *ChatViewModel) Init() tea.Cmd {
	m.syncComponentSizes()

	return tea.Batch(
		m.chats.Init(),
		m.chatHistory.Init(),
		m.chatInput.Init(),
		m.FocusContainer.Init(),
		m.showAllChats(),
	)
}

func (m *ChatViewModel) syncComponentSizes() {
	m.chats.SetSize(m.Frame.Width()/5, m.Frame.Height())

	// TODO: validate the stack component [1]
	gapHeight := lipgloss.Height(m.stack.Components()[1].(*Stack).Gap())
	gapWidth := lipgloss.Width(m.stack.Gap())
	reminingWidth := m.Frame.Width() - m.chats.Width() - gapWidth
	m.chatHistory.SetSize(reminingWidth, m.Frame.Height()-gapHeight-m.chatInput.Height())

	m.chatInput.SetWidth(reminingWidth)
}

func (m *ChatViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// handle updates on frame
	frameCmd := m.Frame.Update(msg)

	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.syncComponentSizes()
	}

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(frameCmd, cmd)
}

func (m *ChatViewModel) View() string {
	return m.Frame.View(m.stack.View())
}

func (m *ChatViewModel) showAllChats() tea.Cmd {
	attendatnces, err := m.attendanceRepo.GetAllBy("user_id", m.user.UniqueId)
	if err != nil {
		return Error("Failed to get attendances: " + err.Error())
	}

	errorCmds := []tea.Cmd{}
	chats := make([]list.Item, 0, len(attendatnces))
	for _, attendance := range attendatnces {
		channel, err := m.channelRepo.GetOneBy("unique_id", attendance.ChannelId)
		if err != nil {
			errorCmds = append(errorCmds, Error("Failed to get channel: "+err.Error()))
		}

		chats = append(chats, Chat{channel.Name, channel.UniqueId})
	}

	m.chats.SetItems(chats)

	return tea.Batch(errorCmds...)
}

func (m *ChatViewModel) showChatHistory(chatId string) tea.Cmd {
	chat, err := m.channelRepo.GetOneBy("unique_id", chatId)
	if err != nil {
		return Error("Failed to get channel: " + err.Error())
	}

	m.chatHistory.Title = chat.Name

	messages, err := m.messageRepo.GetAllBy("channel_id", chat.Id)
	if err != nil {
		return Error("Failed to get messages: " + err.Error())
	}

	chatMessages := make([]ChatMessage, 0, len(messages))

	for _, message := range messages {
		user, err := m.userRepo.GetOne(message.UserId)
		if err != nil {
			return Error("Failed to get user: " + err.Error())
		}
		chatMessages = append(chatMessages, ChatMessage{
			Member: user.Name,
			Text:   message.Text,
			At:     message.CreatedAt,
		})
	}

	m.chatHistory.SetMessages(chatMessages)

	return nil
}
