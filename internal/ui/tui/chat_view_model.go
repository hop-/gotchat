package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatViewModel struct {
	// Screen component
	Screen

	// Focusable component
	*FocusContainer

	// Components
	chats       *ItemList
	chatHistory *ItemList
	chatInput   *ChatInput

	// Stack
	stack *Stack
}

func newChatViewModel() *ChatViewModel {
	// Initialize chat list
	chats := newItemList([]list.Item{
		User{"Alice", "Hello!"},
		User{"Bob", "How are you?"},
		User{"Charlie", "Good morning!"},
	})

	// Initialize chat history
	chatHistory := newItemList([]list.Item{
		User{"Alice", "Hello!"},
		User{"Me", "How are you?"},
		User{"Alice", "I'm good, thanks!"},
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
	}
}

func (m *ChatViewModel) Init() tea.Cmd {
	m.syncSizes()

	return m.FocusContainer.Init()
}

func (m *ChatViewModel) syncSizes() {
	m.chats.SetSize(m.GetWidth()/5, m.GetHeight())
	m.chatHistory.SetSize(m.GetWidth()/5*4, m.GetHeight()-3-m.chatInput.Height())
	m.chatInput.SetWidth(m.GetWidth()/5*4 - 2)
}

func (m *ChatViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// handle screen updates
	screenCmd := m.Screen.Update(msg)

	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.syncSizes()
	}

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(screenCmd, cmd)
}

func (m *ChatViewModel) View() string {
	return m.Screen.View(m.stack.View())
}
