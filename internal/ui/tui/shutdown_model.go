package tui

import tea "github.com/charmbracelet/bubbletea"

type ShutdownModel struct {
	// Frame component
	Frame
}

func newShutdownModel() *ShutdownModel {
	return &ShutdownModel{
		Frame: Frame{},
	}
}

func (m *ShutdownModel) Init() tea.Cmd {
	return nil
}

func (m *ShutdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle updates on frame
	frameCmd := m.Frame.Update(msg)

	return m, tea.Batch(frameCmd, Shutdown)
}

func (m *ShutdownModel) View() string {
	return m.Frame.View("Bye...")
}
