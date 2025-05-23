package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
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
	*FocusContainer

	// Component
	list *ItemList

	// Stack
	stack *Stack

	// Repos
	userRepo core.Repository[core.User]
}

func newUsersListModel(userRepo core.Repository[core.User]) *UsersListModel {
	l := newItemList([]list.Item{
		User{"HoP", "Last login 5 minutes ago"},
		User{"Asd", "Last login 2 hours ago"},
		User{"KAKASH", "Last login 1 day ago"},
		User{"Yesimov", "Last login April 5"},
		User{"Khkhunj", "Last login June 2022"},
		User{"Jujupuluz", "last login December 2021"},
	})
	l.Title = "Users"
	l.OnSelect(func(item list.Item) tea.Cmd {
		if user, ok := item.(User); ok {
			return func() tea.Msg { return PushPageMsg{newSigninModel(user.Title())} }
		}

		return nil
	})

	newLoginButton := newButton("New Login")
	newLoginButton.SetActive(true)
	newLoginButton.OnAction(func() tea.Msg { return PushPageMsg{newSignupModel(userRepo)} })

	exitButton := newButton("Exit")
	exitButton.SetActive(true)
	exitButton.OnAction(internalQuit)

	return &UsersListModel{
		Screen{},
		&FocusContainer{[]FocusableModel{l, newLoginButton, exitButton}, 0},
		l,
		newStack(Vertical, 2, l, newStack(Horizontal, 3, newLoginButton, exitButton)),
		userRepo,
	}
}

func (m *UsersListModel) Init() tea.Cmd {
	users := m.getUsers()
	if len(users) == 0 {
		// TODO: Disable the list if there are no users
	}
	m.list.SetItems(m.getUsers())

	return tea.Batch(m.FocusContainer.Init(), m.stack.Init())
}

func (m *UsersListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle screen updates
	screenCmd := m.Screen.Update(msg)

	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(m.Screen.Width()/2, m.Screen.Height()/2)
	}

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(screenCmd, cmd)
}

func (m *UsersListModel) View() string {
	return m.Screen.View(m.stack.View())
}

func (m *UsersListModel) getUsers() []list.Item {
	users, err := m.userRepo.GetAll()
	if err != nil {
		return nil
	}

	items := make([]list.Item, len(users))
	for i, user := range users {
		items[i] = User{
			name: user.Name,
			info: fmt.Sprintf("Last login: %s", user.LastLogin.String()),
		}
	}

	return items
}
