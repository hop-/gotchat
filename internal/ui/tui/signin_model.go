package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
)

type SigninModel struct {
	// Screen component
	Screen
	// Focusable container
	*FocusContainer

	// Components
	username      *Label
	passwordInput *TextInput
	loginButton   *Button
	backButton    *Button

	// Stack component
	stack *Stack

	// Repos
	userRepo core.Repository[core.User]

	// User ID
	user *core.User
}

func newSigninModel(userId string, userRepo core.Repository[core.User], channelRepo core.Repository[core.Channel]) *SigninModel {
	user, err := userRepo.GetOneBy("unique_id", userId)
	if err != nil {
		// TODO: Handle error
		panic(err)
	}

	usernameLabel := newLabel(user.Name)

	passwordInput := newTextInput("Password")
	passwordInput.Placeholder = "Enter your password"
	passwordInput.CharLimit = 256
	passwordInput.Width = 20
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.SetActive(true)

	loginButton := newButton("Login")
	loginButton.SetActive(false)
	loginButton.OnAction(func() tea.Msg {
		if core.CheckPasswordHash(passwordInput.Value(), user.Password) {
			user.LastLogin = time.Now()
			userRepo.Update(user)

			return SetNewPageMsg{newChatViewModel(channelRepo)}
		}

		return ErrorMsg{Message: "Invalid password"}
	})

	backButton := newButton("Back")
	backButton.SetActive(true)
	backButton.OnAction(func() tea.Msg { return PopPageMsg{} })

	return &SigninModel{
		Screen{},
		&FocusContainer{[]FocusableModel{passwordInput, loginButton, backButton}, 0},

		usernameLabel,
		passwordInput,
		loginButton,
		backButton,
		newStack(Vertical, 1, usernameLabel, passwordInput, newStack(Horizontal, 3, loginButton, backButton)),
		userRepo,
		user,
	}
}

func (m *SigninModel) Init() tea.Cmd {
	m.updateActiveStates()

	return tea.Batch(m.FocusContainer.Init(), m.stack.Init())
}

func (m *SigninModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle screen updates
	screenCmd := m.Screen.Update(msg)

	m.updateActiveStates()

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(screenCmd, cmd)
}

func (m *SigninModel) updateActiveStates() {
	if m.passwordInput.Value() != "" {
		m.loginButton.SetActive(true)
	} else {
		m.loginButton.SetActive(false)
	}
}

func (m *SigninModel) View() string {
	return m.Screen.View(m.stack.View())
}
