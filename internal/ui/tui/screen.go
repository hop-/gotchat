package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	screenWidth  = 80
	screenHeight = 24
)

type Screen struct {
}

func (m *Screen) update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		screenWidth = msg.Width - 4
		screenHeight = msg.Height - 2
	}
}

func (m *Screen) view(content string) string {
	return boarderStyle.Render(lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, content))
}
