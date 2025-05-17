package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	screenWidth  = 80
	screenHeight = 24
)

type Screen struct{}

func (m *Screen) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		screenWidth = msg.Width - 4
		screenHeight = msg.Height - 2
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return tea.Quit
		}
	}

	return nil
}

func (m *Screen) View(content string) string {
	return boarderStyle.Render(lipgloss.Place(screenWidth, screenHeight, lipgloss.Center, lipgloss.Center, content))
}

func (m *Screen) GetWidth() int {
	return screenWidth
}

func (m *Screen) GetHeight() int {
	return screenHeight
}
