package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	frameWidth  = 80
	frameHeight = 24
)

type Frame struct {
	errors []string
}

func (m *Frame) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ErrorMsg:
		m.addError(msg.Message)
	case tea.WindowSizeMsg:
		frameWidth = msg.Width - 4
		frameHeight = msg.Height - 2
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "ctrl+d" || msg.String() == "ctrl+q" {
			return InternalQuit
		}
	}

	return nil
}

func (m *Frame) View(content string) string {
	// TODO: improve
	for _, e := range m.errors {
		content += "\n" + e
	}

	return boarderStyle.Render(lipgloss.Place(frameWidth, frameHeight, lipgloss.Center, lipgloss.Center, content))
}

func (m *Frame) Width() int {
	return frameWidth
}

func (m *Frame) Height() int {
	return frameHeight
}

func (m *Frame) addError(e string) {
	m.errors = append(m.errors, e)
}
