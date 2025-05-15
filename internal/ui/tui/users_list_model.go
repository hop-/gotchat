package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type User struct {
	name, info string
}

func (i User) Title() string       { return i.name }
func (i User) Description() string { return i.info }
func (i User) FilterValue() string { return i.name }

type UsersListModel struct {
	// Screen component
	Screen
	// Focusable component
	FocusContainer

	// Component
	list *ItemList

	// Stack
	stack *Stack
}

func newUsersListModel() *UsersListModel {
	l := newItemList([]list.Item{
		User{"HoP", "Last login 5 minutes ago"},
		User{"Asd", "Last login 2 hours ago"},
		User{"KAKASH", "Last login 1 day ago"},
		User{"Yesimov", "Last login April 5"},
		User{"Khkhunj", "Last login June 2022"},
		User{"Jujupuluz", "last login December 2021"},
	})
	l.Title = "Users"

	newLoginButton := newButton("New Login")
	newLoginButton.SetActive(true)
	newLoginButton.OnAction(func() tea.Msg { return PushPageMsg{newLoginModel()} })

	exitButton := newButton("Exit")
	exitButton.SetActive(true)
	exitButton.OnAction(tea.Quit)

	return &UsersListModel{
		Screen{},
		FocusContainer{[]FocusableModel{l, newLoginButton, exitButton}, 0},
		l,
		newStack(Vertical, 3, l, newStack(Horizontal, 2, newLoginButton, exitButton)),
	}
}

func (m *UsersListModel) Init() tea.Cmd {
	return tea.Batch(m.FocusContainer.Init(), m.stack.Init())
}

func (m *UsersListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update screen size
	m.Screen.update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width/2, msg.Height/2)
	}

	_, cmd := m.FocusContainer.Update(msg)

	return m, cmd
}

func (m *UsersListModel) View() string {
	return m.Screen.view(m.stack.View())
}
