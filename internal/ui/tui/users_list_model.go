package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/services"
	"github.com/hop-/gotchat/internal/ui/tui/commands"
	"github.com/hop-/gotchat/internal/ui/tui/components"
)

type User struct {
	id, name, lastLogin string
}

func (i User) Title() string       { return i.name }
func (i User) Description() string { return i.lastLogin }
func (i User) FilterValue() string { return i.name }

type UsersListModel struct {
	// Frame component
	components.Frame
	// Focusable component
	*components.FocusContainer

	// Component
	list *components.ItemList

	// Stack
	stack *components.Stack

	// Services
	userManager *services.UserManager
}

func newUsersListModel(
	userManager *services.UserManager,
	chatManager *services.ChatManager,
) *UsersListModel {
	l := components.NewItemList([]list.Item{})
	l.Title = "Users"
	l.OnSelect(func(item list.Item) tea.Cmd {
		if userItem, ok := item.(User); ok {
			user, err := userManager.GetUserByUniqueId(userItem.id)
			if err != nil {
				return commands.Error(err.Error())
			}

			return commands.PushPage(newSigninModel(user, userManager, chatManager))
		}

		return nil
	})

	newLoginButton := components.NewButton("New Login")
	newLoginButton.SetActive(true)
	newLoginButton.OnAction(commands.PushPage(newSignupModel(userManager, chatManager)))

	exitButton := components.NewButton("Exit")
	exitButton.SetActive(true)
	exitButton.OnAction(commands.Shutdown)

	return &UsersListModel{
		components.Frame{},
		components.NewFocusContainer([]components.FocusableModel{l, newLoginButton, exitButton}),
		l,
		components.NewStack(
			components.Vertical, 2,
			l, components.NewStack(
				components.Horizontal, 3,
				newLoginButton, exitButton,
			),
		),
		userManager,
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
	users, err := m.userManager.GetAllUsers()
	if err != nil {
		m.Frame.AddError(err.Error())
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
