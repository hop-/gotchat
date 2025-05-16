package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SigninModel struct {
	// Screen component
	Screen
	// Focusable container
	FocusContainer

	// Components
	username      *Label
	passwordInput *TextInput
	loginButton   *Button
	backButton    *Button

	// Stack component
	stack *Stack
}

func newSigninModel(username string) *SigninModel {
	usernameLabel := newLabel(username)

	passwordInput := newTextInput("Password")
	passwordInput.Placeholder = "Enter your password"
	passwordInput.CharLimit = 256
	passwordInput.Width = 20
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.SetActive(true)

	loginButton := newButton("Login")
	loginButton.SetActive(false)
	loginButton.OnAction(func() tea.Msg { return SetNewPageMsg{newChatViewModel()} })

	backButton := newButton("Back")
	backButton.SetActive(true)
	backButton.OnAction(func() tea.Msg { return PopPageMsg{} })

	return &SigninModel{
		Screen{},
		FocusContainer{[]FocusableModel{passwordInput, loginButton, backButton}, 0},

		usernameLabel,
		passwordInput,
		loginButton,
		backButton,
		newStack(Vertical, 3, usernameLabel, passwordInput, newStack(Horizontal, 2, loginButton, backButton)),
	}
}

func (m *SigninModel) Init() tea.Cmd {
	m.updateActiveStates()

	return tea.Batch(m.FocusContainer.Init(), m.stack.Init())
}

func (m *SigninModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *SigninModel) updateActiveStates() {
	if m.passwordInput.Value() != "" {
		m.loginButton.SetActive(true)
	} else {
		m.loginButton.SetActive(false)
	}
}

func (m *SigninModel) View() string {
	return m.Screen.view(m.stack.View())
}
