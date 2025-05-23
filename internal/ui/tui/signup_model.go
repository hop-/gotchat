package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
)

type SignupModel struct {
	// Screen component
	Screen
	// Focusable container
	*FocusContainer

	// Components
	usernameInput *TextInput
	passwordInput *TextInput
	loginButton   *Button
	backButton    *Button

	// Stack component
	stack *Stack

	// Repos
	userRepo core.Repository[core.User]
}

func newSignupModel(userRepo core.Repository[core.User]) *SignupModel {
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

	loginButton := newButton("Login")
	loginButton.SetActive(false)
	loginButton.OnAction(func() tea.Msg {
		err := userRepo.Create(core.NewUser(usernameInput.Value()))
		if err != nil {
			fmt.Println("Error creating user:", err)
			return nil
		}

		return SetNewPageMsg{newChatViewModel()}
	})

	backButton := newButton("Back")
	backButton.SetActive(true)
	backButton.OnAction(func() tea.Msg { return PopPageMsg{} })

	return &SignupModel{
		Screen{},
		&FocusContainer{[]FocusableModel{usernameInput, passwordInput, loginButton, backButton}, 0},
		usernameInput,
		passwordInput,
		loginButton,
		backButton,
		newStack(Vertical, 1, usernameInput, passwordInput, newStack(Horizontal, 3, loginButton, backButton)),
		userRepo,
	}
}

func (m *SignupModel) Init() tea.Cmd {
	m.updateActiveStates()

	return tea.Batch(m.FocusContainer.Init(), m.stack.Init())
}

func (m *SignupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle screen updates
	screenCmd := m.Screen.Update(msg)

	m.updateActiveStates()

	fc, cmd := m.FocusContainer.Update(msg)
	m.FocusContainer = fc

	return m, tea.Batch(screenCmd, cmd)
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
	return m.Screen.View(m.stack.View())
}
