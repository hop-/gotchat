package tui

import tea "github.com/charmbracelet/bubbletea"

type ChatViewModel struct {
	Screen
}

func newChatViewModel() *ChatViewModel {
	return &ChatViewModel{
		Screen{},
	}
}

func (c *ChatViewModel) Init() tea.Cmd {
	return nil
}

func (c *ChatViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c *ChatViewModel) View() string {
	return c.Screen.view("")
}
