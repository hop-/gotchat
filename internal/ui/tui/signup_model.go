package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/services"
)

type SignupModel struct {
	// Frame component
	Frame
	// Focusable container
	*FocusContainer

	// Components
	usernameInput *TextInput
	passwordInput *TextInput
	loginButton   *Button
	backButton    *Button

	// Stack component
	stack *Stack

	// Services
	userManager *services.UserManager
}

func newSignupModel(
	userManager *services.UserManager,
	chatManager *services.ChatManager,
) *SignupModel {
	usernameInput := newTextInput("Username")
	usernameInput.Placeholder = "Enter your nickname"
	usernameInput.CharLimit = 128
	usernameInput.Width = 20

	passwordInput := newTextInput("Password")
	passwordInput.Placeholder = "Enter your password"
	passwordInput.CharLimit = 256
	passwordInput.Width = 20
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'

	// TODO: add password confirmation input

	loginButton := newButton("Login")
	loginButton.SetActive(false)
	loginButton.OnAction(func() tea.Msg {
		passwordHash, err := core.HashPassword(passwordInput.Value())
		if err != nil {
			// TODO: handle error properly
			return nil
		}

		user := core.NewUser(usernameInput.Value(), passwordHash)
		err = userManager.CreateUser(user)
		if err != nil {
			// TODO: handle error properly
			return nil
		}

		return SetNewPageMsg{newChatViewModel(user, userManager, chatManager)}
	})

	backButton := newButton("Back")
	backButton.SetActive(true)
	backButton.OnAction(PopPage)

	return &SignupModel{
		Frame{},
		&FocusContainer{[]FocusableModel{usernameInput, passwordInput, loginButton, backButton}, 0},
		usernameInput,
		passwordInput,
		loginButton,
		backButton,
		newStack(Vertical, 1, usernameInput, passwordInput, newStack(Horizontal, 3, loginButton, backButton)),
		userManager,
	}
}

func (m *SignupModel) Init() tea.Cmd {
	m.updateActiveStates()

	return tea.Batch(m.FocusContainer.Init(), m.stack.Init())
}

func (m *SignupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle supdates on frame
	frameCmd := m.Frame.Update(msg)

	m.updateActiveStates()

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(frameCmd, cmd)
}

func (m *SignupModel) updateActiveStates() {
	if m.usernameInput.Value() != "" {
		m.passwordInput.SetActive(true)
		if m.passwordInput.Value() != "" {
			m.loginButton.SetActive(true)
		} else {
			m.loginButton.SetActive(false)
		}
	} else {
		m.passwordInput.SetActive(false)
		m.loginButton.SetActive(false)
	}
}

func (m *SignupModel) View() string {
	return m.Frame.View(m.stack.View())
}
