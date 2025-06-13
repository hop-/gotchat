package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatInputMessageSentMsg struct {
	Content string
}

func ChatInputMessageSent(content string) tea.Cmd {
	return func() tea.Msg {
		return ChatInputMessageSentMsg{content}
	}
}

type ChatInput struct {
	textarea.Model
	isActive bool
}

func newChatInput() *ChatInput {
	textArea := textarea.New()
	textArea.Placeholder = "Send a message..."
	textArea.Prompt = "┃ "
	textArea.SetHeight(3)
	textArea.SetWidth(20)
	textArea.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textArea.ShowLineNumbers = false

	return &ChatInput{
		textArea,
		true,
	}
}

func (ci *ChatInput) Init() tea.Cmd {
	return nil
}

func (ci *ChatInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd
	ci.Model, cmd = ci.Model.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if ci.Model.Value() != "" {
				// Send the message
				cmds = append(cmds, ChatInputMessageSent(ci.Model.Value()))
				// Reset the input field
				ci.Model.Reset()
			}
		}
	}

	return ci, tea.Batch(cmds...)
}

func (ci *ChatInput) View() string {
	return ci.Model.View()
}

func (ci *ChatInput) Focus() tea.Cmd {
	return ci.Model.Focus()
}

func (ci *ChatInput) Blur() tea.Cmd {
	ci.Model.Blur()

	return nil
}

func (ci *ChatInput) SetActive(active bool) tea.Cmd {
	ci.isActive = active

	return nil
}

func (ci *ChatInput) IsActive() bool {
	return ci.isActive
}
