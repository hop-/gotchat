package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	// Screen component
	Screen
	// Focusable container
	FocusContainer

	// Components
	usernameInput *TextInput
	passwordInput *TextInput
	submitButton  *Button
	backButton    *Button

	// Stack component
	stack *Stack
}

func newLoginModel() *LoginModel {
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

	submitButton := newButton("Submit")
	submitButton.SetActive(false)
	submitButton.OnAction(tea.Quit)

	backButton := newButton("Back")
	backButton.SetActive(true)
	backButton.OnAction(func() tea.Msg { return PopPageMsg{} })

	return &LoginModel{
		Screen{},
		FocusContainer{[]FocusableModel{usernameInput, passwordInput, submitButton, backButton}, 0},
		usernameInput,
		passwordInput,
		submitButton,
		backButton,
		newStack(Vertical, 3, usernameInput, passwordInput, newStack(Horizontal, 2, submitButton, backButton)),
	}
}

func (m *LoginModel) Init() tea.Cmd {
	m.updateActiveStates()

	return tea.Batch(m.FocusContainer.Init(), m.stack.Init())
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update screen size
	m.Screen.update(msg)
	m.updateActiveStates()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	_, cmd := m.FocusContainer.Update(msg)

	return m, cmd
}

func (m *LoginModel) updateActiveStates() {
	if m.usernameInput.Value() != "" {
		m.passwordInput.SetActive(true)
		if m.passwordInput.Value() != "" {
			m.submitButton.SetActive(true)
		} else {
			m.submitButton.SetActive(false)
		}
	} else {
		m.passwordInput.SetActive(false)
		m.submitButton.SetActive(false)
	}
}

func (m *LoginModel) View() string {
	return m.Screen.view(m.stack.View())
}
