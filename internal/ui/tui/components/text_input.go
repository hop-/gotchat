package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInput struct {
	textinput.Model
	label    string
	isActive bool
}

func NewTextInput(label string) *TextInput {
	ti := textinput.New()

	return &TextInput{
		Model:    ti,
		label:    label,
		isActive: true,
	}
}

func (ti *TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (ti *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	ti.Model, cmd = ti.Model.Update(msg)

	return ti, cmd
}

func (ti *TextInput) View() string {
	label := ""
	if ti.label != "" {
		if ti.Model.Focused() {
			label = focusedStyle.Render(ti.label)
		} else {
			label = blurredStyle.Render(ti.label)
		}
		label += "\n"
	}

	return fmt.Sprintf("%s%s", label, ti.Model.View())
}

func (ti *TextInput) Focus() tea.Cmd {
	ti.PromptStyle = focusedStyle
	ti.TextStyle = focusedStyle

	return ti.Model.Focus()
}

func (ti *TextInput) Blur() tea.Cmd {
	if ti.isActive {
		ti.PromptStyle = noStyle
	} else {
		ti.PromptStyle = blurredStyle
	}
	ti.TextStyle = blurredStyle

	ti.Model.Blur()

	return nil
}

func (ti *TextInput) SetActive(active bool) tea.Cmd {
	ti.isActive = active
	if !ti.Focused() {
		if active {
			ti.PromptStyle = noStyle
		} else {
			ti.PromptStyle = blurredStyle
		}
	}

	return nil
}

func (ti *TextInput) IsActive() bool {
	return ti.isActive
}
