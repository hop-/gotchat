package components

import (
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	yourColor    = lipgloss.Color("205")
	unknownColor = lipgloss.Color("240")
	memberColors = []lipgloss.Color{
		lipgloss.Color("39"),
		lipgloss.Color("45"),
		lipgloss.Color("33"),
		lipgloss.Color("63"),
		lipgloss.Color("57"),
		lipgloss.Color("75"),
		lipgloss.Color("170"),
	}
)

type ChatMessage struct {
	Member string
	Text   string
	At     time.Time
}

type ChatHistory struct {
	viewport.Model
	memberStyles map[string]lipgloss.Style
	messages     []ChatMessage
	Title        string
	titleStyle   *lipgloss.Style
}

func NewChatHistory(you string, otherMembers ...string) *ChatHistory {
	vp := viewport.New(2, 2)
	vp.SetContent("Chat history will be displayed here.")

	memberStyles := make(map[string]lipgloss.Style, 1+len(otherMembers))
	memberStyles[you] = lipgloss.NewStyle().Foreground(yourColor).Bold(true).Underline(true)
	for i, name := range otherMembers {
		memberStyles[name] = lipgloss.NewStyle().Foreground(memberColors[i]).Bold(true).Italic(true).Underline(true)
	}

	return &ChatHistory{
		vp,
		memberStyles,
		[]ChatMessage{},
		"Chat History",
		&focusedTitleStyle,
	}
}

func (ch *ChatHistory) Init() tea.Cmd {
	ch.Model.SetContent(ch.renderMessages())
	ch.Model.GotoBottom()

	return ch.Model.Init()
}

func (ch *ChatHistory) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	vp, cmd := ch.Model.Update(msg)
	ch.Model = vp

	return ch, cmd
}

func (ch *ChatHistory) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, titleBarStyle.Render(ch.titleStyle.Render(ch.Title)), ch.Model.View())
}

func (ch *ChatHistory) Focus() tea.Cmd {
	ch.titleStyle = &focusedTitleStyle

	return nil
}

func (ch *ChatHistory) Blur() tea.Cmd {
	ch.titleStyle = &blurredTitleStyle

	return nil
}

func (ch *ChatHistory) Focused() bool {
	return true
}

func (ch *ChatHistory) IsActive() bool {
	return true
}

func (ch *ChatHistory) SetActive(active bool) {
	// Do nothing
}

func (ch *ChatHistory) SetHeight(height int) {
	ch.Model.Height = height - 2
}

func (ch *ChatHistory) SetWidth(width int) {
	ch.Model.Width = width
}

func (ch *ChatHistory) SetSize(width, height int) {
	ch.SetWidth(width)
	ch.SetHeight(height)
}

func (ch *ChatHistory) renderMessages() string {
	components := make([]string, 0, len(ch.messages))
	for _, message := range ch.messages {
		style, ok := ch.memberStyles[message.Member]
		if !ok {
			style = lipgloss.NewStyle().Foreground(unknownColor).Italic(true)
		}

		components = append(components, style.Render(message.Member)+":\n\t"+message.Text)
	}

	return lipgloss.JoinVertical(lipgloss.Left, components...)
}

func (ch *ChatHistory) SetMessages(messages []ChatMessage) {
	ch.messages = messages
	ch.SetContent(ch.renderMessages())
	ch.GotoBottom()
}

func (ch *ChatHistory) AddMessage(message ChatMessage) {
	ch.messages = append(ch.messages, message)
	ch.SetContent(ch.renderMessages())
	ch.GotoBottom()
}

func (ch *ChatHistory) UpdateMessage(index int, message ChatMessage) {
	if index < 0 || index >= len(ch.messages) {
		return
	}
	ch.messages[index] = message
	ch.SetContent(ch.renderMessages())
}
