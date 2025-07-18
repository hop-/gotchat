package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/services"
	"github.com/hop-/gotchat/internal/ui/tui/commands"
	"github.com/hop-/gotchat/internal/ui/tui/components"
)

type SigninModel struct {
	// Frame component
	components.Frame
	// Focusable container
	*components.FocusContainer

	// Components
	username      *components.Label
	passwordInput *components.TextInput
	loginButton   *components.Button
	backButton    *components.Button

	// Stack component
	stack *components.Stack

	// Services
	userManager *services.UserManager
}

func newSigninModel(
	user *core.User,
	userManager *services.UserManager,
	chatManager *services.ChatManager,
) *SigninModel {
	usernameLabel := components.NewLabel(user.Name)

	passwordInput := components.NewTextInput("Password")
	passwordInput.Placeholder = "Enter your password"
	passwordInput.CharLimit = 256
	passwordInput.Width = 20
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.SetActive(true)

	loginButton := components.NewButton("Login")
	loginButton.SetActive(false)
	loginButton.OnAction(func() tea.Msg {
		user, err := userManager.LoginUser(user, passwordInput.Value())
		if err != nil {
			switch err {
			case services.ErrorInvalidInput, services.ErrorInvalidCredentials:
				return commands.ErrorMsg{Message: "Invalid credentials"}
			}

			return commands.ErrorMsg{Message: "An error occurred while logging in"}
		}

		return commands.SetNewPageMsg{Page: newChatViewModel(user, userManager, chatManager)}
	})

	backButton := components.NewButton("Back")
	backButton.SetActive(true)
	backButton.OnAction(commands.PopPage)

	return &SigninModel{
		components.Frame{},
		components.NewFocusContainer(passwordInput, loginButton, backButton),

		usernameLabel,
		passwordInput,
		loginButton,
		backButton,
		components.NewStack(
			components.Vertical, 1,
			usernameLabel, passwordInput, components.NewStack(
				components.Horizontal, 3,
				loginButton, backButton,
			),
		),
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
