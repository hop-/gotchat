package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
)

type User struct {
	id, name, lastLogin string
}

func (i User) Title() string       { return i.name }
func (i User) Description() string { return i.lastLogin }
func (i User) FilterValue() string { return i.name }

type UsersListModel struct {
	// Frame component
	Frame
	// Focusable component
	*FocusContainer

	// Component
	list *ItemList

	// Stack
	stack *Stack

	// Repos
	userRepo core.Repository[core.User]
}

func newUsersListModel(
	userRepo core.Repository[core.User],
	channelRepo core.Repository[core.Channel],
	attendanceRepo core.Repository[core.Attendance],
) *UsersListModel {
	l := newItemList([]list.Item{})
	l.Title = "Users"
	l.OnSelect(func(item list.Item) tea.Cmd {
		if userItem, ok := item.(User); ok {
			user, err := userRepo.GetOneBy("unique_id", userItem.id)
			if err != nil {
				return Error(err.Error())
			}

			return func() tea.Msg { return PushPageMsg{newSigninModel(user, userRepo, channelRepo, attendanceRepo)} }
		}

		return nil
	})

	newLoginButton := newButton("New Login")
	newLoginButton.SetActive(true)
	newLoginButton.OnAction(func() tea.Msg { return PushPageMsg{newSignupModel(userRepo, channelRepo, attendanceRepo)} })

	exitButton := newButton("Exit")
	exitButton.SetActive(true)
	exitButton.OnAction(InternalQuit)

	return &UsersListModel{
		Frame{},
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
	// Handle updates on frame
	frameCmd := m.Frame.Update(msg)

	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(m.Frame.Width()/2, m.Frame.Height()/2)
	}

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(frameCmd, cmd)
}

func (m *UsersListModel) View() string {
	return m.Frame.View(m.stack.View())
}

func (m *UsersListModel) getUsers() []list.Item {
	users, err := m.userRepo.GetAll()
	if err != nil {
		m.Frame.addError(err.Error())
		return nil
	}

	items := make([]list.Item, len(users))
	for i, user := range users {
		items[i] = User{
			id:        user.UniqueId,
			name:      user.Name,
			lastLogin: FormatLastLogin(user.LastLogin),
		}
	}

	return items
}
