package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/services"
	"github.com/hop-/gotchat/internal/ui/tui/commands"
	"github.com/hop-/gotchat/internal/ui/tui/components"
)

type SentMessageToChatMsg struct {
	ChatId  string
	Message string
}

func SendMessageToChat(chatId, message string) tea.Cmd {
	return func() tea.Msg {
		return SentMessageToChatMsg{
			ChatId:  chatId,
			Message: message,
		}
	}
}

type Chat struct {
	services.Chat
}

// FilterValue implements list.Item.
func (c Chat) FilterValue() string {
	return c.Name
}

type ChatViewModel struct {
	// Frame component
	components.Frame

	// Focusable component
	*components.FocusContainer

	// Components
	chats               *components.ItemList
	chatHistory         *components.ChatHistory
	chatInput           *components.ChatInput
	newConnectionButton *components.Button

	// Stack
	stack *components.Stack

	// Services
	userManager *services.UserManager
	chatManager *services.ChatManager

	// User entity
	user *core.User
}

func newChatViewModel(
	user *core.User,
	userManager *services.UserManager,
	chatManager *services.ChatManager,
) *ChatViewModel {
	// Initialize chat list
	chats := components.NewItemList([]list.Item{})
	chats.Title = "Chats"

	// Initialize chat history
	chatHistory := components.NewChatHistory("You", "TheOtherOne")

	// Initialize chat input
	chatInput := components.NewChatInput()
	chatInput.Placeholder = "Type a message..."
	chatInput.SetActive(true)
	chatInput.SetWidth(20)

	newConnectionButton := components.NewButton("New Connection")

	return &ChatViewModel{
		components.Frame{},
		components.NewFocusContainer(chatInput, chats, newConnectionButton, chatHistory),
		chats,
		chatHistory,
		chatInput,
		newConnectionButton,
		components.NewStack(
			components.Horizontal, 3,
			components.NewStack(components.Vertical, 1, chats, newConnectionButton, components.NewLabel("")),
			components.NewStack(components.Vertical, 2, chatHistory, chatInput),
		),
		userManager,
		chatManager,
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
	gapHeight1 := lipgloss.Height(m.stack.Components()[0].(*components.Stack).Gap())
	m.chats.SetSize(m.Frame.Width()/5, m.Frame.Height()-2*gapHeight1-2)

	// TODO: validate the stack component [1]
	gapHeight2 := lipgloss.Height(m.stack.Components()[1].(*components.Stack).Gap())
	gapWidth := lipgloss.Width(m.stack.Gap())
	reminingWidth := m.Frame.Width() - m.chats.Width() - gapWidth
	m.chatHistory.SetSize(reminingWidth, m.Frame.Height()-gapHeight2-m.chatInput.Height())

	m.chatInput.SetWidth(reminingWidth)
}

func (m *ChatViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle updates on frame
	cmds := make([]tea.Cmd, 0)
	frameCmd := m.Frame.Update(msg)
	cmds = append(cmds, frameCmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.syncComponentSizes()
	case components.ChatInputMessageSentMsg:
		currnetChat := m.chats.SelectedItem()
		if currnetChat == nil {
			cmds = append(cmds, commands.Error("No chat selected"))
		} else {
			cmds = append(cmds, SendMessageToChat(currnetChat.(Chat).Id, msg.Message))
		}
	}

	fc, cmd := m.FocusContainer.Update(msg)
	cmds = append(cmds, cmd)
	m.FocusContainer = fc

	return m, tea.Batch(cmds...)
}

func (m *ChatViewModel) View() string {
	return m.Frame.View(m.stack.View())
}

func (m *ChatViewModel) showAllChats() tea.Cmd {
	chats, err := m.chatManager.GetChatsByUserId(m.user.Id)
	if err != nil {
		return commands.Error(err.Error())
	}

	chatItems := make([]list.Item, len(chats))
	for i, chat := range chats {
		chatItems[i] = Chat{chat}
	}
	m.chats.SetItems(chatItems)

	return nil
}

func (m *ChatViewModel) showChatHistory(chat Chat) tea.Cmd {
	m.chatHistory.Title = chat.Name

	chatMessages, err := m.chatManager.GetChatMessagesByChatId(chat.Id)
	if err != nil {
		return commands.Error(err.Error())
	}

	chatMessageItems := make([]components.ChatMessage, len(chatMessages))
	for _, message := range chatMessages {
		chatMessageItems = append(chatMessageItems, components.ChatMessage{
			Member: message.Member,
			Text:   message.Text,
			At:     message.At,
		})
	}

	m.chatHistory.SetMessages(chatMessageItems)

	return nil
}
