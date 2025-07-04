package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boarderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Align(lipgloss.Center).
			Padding(0, 0)

	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "255", Dark: "187"})
	blurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "240"})
	noStyle           = lipgloss.NewStyle()
	focusedTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 1)
	blurredTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Background(lipgloss.Color("244")).
				Padding(0, 1)
	titleBarStyle = lipgloss.NewStyle().Padding(0, 0, 1, 2)
)

type Focusable interface {
	Focus() tea.Cmd
	Blur() tea.Cmd
	Focused() bool
}

type Activatable interface {
	SetActive(active bool) tea.Cmd
	IsActive() bool
}

type FocusContainer1 interface {
	ChangeFocus(val int, loop bool) (status bool, cmd tea.Cmd)
}

type FocusableModel interface {
	tea.Model
	Focusable
}

type FocusableActivatableModel interface {
	tea.Model
	Focusable
	Activatable
}

// Custm messages

type SetNewPageMsg struct {
	Page tea.Model
}

type PushPageMsg struct {
	Page tea.Model
}

type PopPageMsg struct{}

type ShutdownMsg struct{}

type InternalQuitMsg struct{}

type ErrorMsg struct {
	Message string
}

// Custom commands and command factories

func SetNewPage(page tea.Model) tea.Cmd {
	return func() tea.Msg {
		return SetNewPageMsg{page}
	}
}

func PushPage(page tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushPageMsg{page}
	}
}

func PopPage() tea.Msg {
	return PopPageMsg{}
}

func Shutdown() tea.Msg {
	return ShutdownMsg{}
}

func InternalQuit() tea.Msg {
	return InternalQuitMsg{}
}

func Error(msg string) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{msg}
	}
}
