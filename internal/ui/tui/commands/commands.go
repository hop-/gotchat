package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

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

type ConnectMsg struct {
	Host string
	Port string
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

func Connect(host string, port string) tea.Cmd {
	return func() tea.Msg {
		return ConnectMsg{host, port}
	}
}
