package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hop-/gotchat/internal/core"
)

type ChatViewModel struct {
	// Screen component
	Screen

	// Focusable component
	*FocusContainer

	// Components
	chats       *ItemList
	chatHistory *ChatHistory
	chatInput   *ChatInput

	// Stack
	stack *Stack

	// Repos
	channelRepo core.Repository[core.Channel]
}

func newChatViewModel(channelRepo core.Repository[core.Channel]) *ChatViewModel {
	// Initialize chat list
	chats := newItemList([]list.Item{})
	chats.Title = "Chats"

	// Initialize chat history
	chatHistory := newChatHistory("You", "TheOtherOne")
	chatHistory.SetMessages([]ChatMessage{
		{"TheOtherOne", "Hello! How are you?", time.Now()},
		{"You", "I'm good, thanks!", time.Now()},
		{"You", "What about you?", time.Now()},
		{"TheOtherOne", "I'm doing well, just working on some projects.", time.Now()},
		{"You", "That's great to hear!", time.Now()},
		{"TheOtherOne", "What about you?", time.Now()},
		{"You", "Just the usual, you know.", time.Now()},
		{"TheOtherOne", "Yeah, I get that.", time.Now()},
		{"TheOtherOne", "Have you seen the latest news?", time.Now()},
		{"You", "No, I haven't. What's going on?", time.Now()},
	})

	// Initialize chat input
	chatInput := newChatInput()
	chatInput.Placeholder = "Type a message..."
	chatInput.SetActive(true)
	chatInput.SetWidth(20)

	return &ChatViewModel{
		Screen{},
		&FocusContainer{[]FocusableModel{chatInput, chats, chatHistory}, 0},

		chats,
		chatHistory,
		chatInput,
		newStack(Horizontal, 3, chats, newStackWithPosition(lipgloss.Left, Vertical, 2, chatHistory, chatInput)),
		channelRepo,
	}
}

func (m *ChatViewModel) Init() tea.Cmd {
	m.syncComponentSizes()

	return tea.Batch(
		m.chats.Init(),
		m.chatHistory.Init(),
		m.chatInput.Init(),
		m.FocusContainer.Init(),
	)
}

func (m *ChatViewModel) syncComponentSizes() {
	m.chats.SetSize(m.Screen.Width()/5, m.Screen.Height())

	// TODO: validate the stack component [1]
	gapHeight := lipgloss.Height(m.stack.Components()[1].(*Stack).Gap())
	gapWidth := lipgloss.Width(m.stack.Gap())
	reminingWidth := m.Screen.Width() - m.chats.Width() - gapWidth
	m.chatHistory.SetSize(reminingWidth, m.Screen.Height()-gapHeight-m.chatInput.Height())

	m.chatInput.SetWidth(reminingWidth)
}

func (m *ChatViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// handle screen updates
	screenCmd := m.Screen.Update(msg)

	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.syncComponentSizes()
	}

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(screenCmd, cmd)
}

func (m *ChatViewModel) View() string {
	return m.Screen.View(m.stack.View())
}
