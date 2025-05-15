package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Button struct {
	text     string
	focus    bool
	isActive bool

	focusedButton  string
	blurredButton  string
	inactiveButton string
	onActionCmd    tea.Cmd
}

func newButton(text string) *Button {
	buttonText := fmt.Sprintf("[ %s ]", text)

	return &Button{
		text,
		true,
		false,
		focusedStyle.Render(buttonText),
		fmt.Sprintf("[ %s ]", blurredStyle.Render(text)),
		blurredStyle.Render(buttonText),
		nil,
	}
}

func (b *Button) Init() tea.Cmd {
	// No initialization needed for button
	return nil
}

func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !b.focus {
		return b, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return b, b.onActionCmd
		}
	}

	return b, nil
}

func (b *Button) View() string {
	if b.focus {
		return b.focusedButton
	} else if !b.isActive {
		return b.inactiveButton
	}

	return b.blurredButton
}

func (b *Button) Focus() tea.Cmd {
	b.focus = true

	return nil
}

func (b *Button) Blur() tea.Cmd {
	b.focus = false
	return nil
}

func (b *Button) Focused() bool {
	return b.focus
}

func (b *Button) SetActive(active bool) tea.Cmd {
	b.isActive = active

	return nil
}

func (b *Button) IsActive() bool {
	return b.isActive
}

func (b *Button) OnAction(cmd tea.Cmd) {
	b.onActionCmd = cmd
}
