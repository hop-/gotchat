package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/services"
)

type SigninModel struct {
	// Frame component
	Frame
	// Focusable container
	*FocusContainer

	// Components
	username      *Label
	passwordInput *TextInput
	loginButton   *Button
	backButton    *Button

	// Stack component
	stack *Stack

	// Services
	userManager *services.UserManager
}

func newSigninModel(
	user *core.User,
	userManager *services.UserManager,
	chatManager *services.ChatManager,
) *SigninModel {
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
			userManager.UpdateUser(user)

			return SetNewPageMsg{newChatViewModel(user, userManager, chatManager)}
		}

		return ErrorMsg{Message: "Invalid password"}
	})

	backButton := newButton("Back")
	backButton.SetActive(true)
	backButton.OnAction(PopPage)

	return &SigninModel{
		Frame{},
		&FocusContainer{[]FocusableModel{passwordInput, loginButton, backButton}, 0},

		usernameLabel,
		passwordInput,
		loginButton,
		backButton,
		newStack(Vertical, 1, usernameLabel, passwordInput, newStack(Horizontal, 3, loginButton, backButton)),
		userManager,
	}
}

func (m *SigninModel) Init() tea.Cmd {
	m.updateActiveStates()

	return tea.Batch(m.FocusContainer.Init(), m.stack.Init())
}

func (m *SigninModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle updates on frame
	frameCmd := m.Frame.Update(msg)

	m.updateActiveStates()

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(frameCmd, cmd)
}

func (m *SigninModel) View() string {
	return m.Frame.View(m.stack.View())
}

func (m *SigninModel) updateActiveStates() {
	if m.passwordInput.Value() != "" {
		m.loginButton.SetActive(true)
	} else {
		m.loginButton.SetActive(false)
	}
}
