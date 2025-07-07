package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/ui/tui/commands"
	"github.com/hop-/gotchat/internal/ui/tui/components"
)

type ShutdownModel struct {
	// Frame component
	components.Frame
}

func newShutdownModel() *ShutdownModel {
	return &ShutdownModel{
		Frame: components.Frame{},
	}
}

func (m *ShutdownModel) Init() tea.Cmd {
	return nil
}

func (m *ShutdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle updates on frame
	frameCmd := m.Frame.Update(msg)

	return m, tea.Batch(frameCmd, commands.InternalQuit)
}

func (m *ShutdownModel) View() string {
	return m.Frame.View("Bye...")
}
