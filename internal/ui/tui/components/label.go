package components

import tea "github.com/charmbracelet/bubbletea"

type Label struct {
	text string
}

func NewLabel(text string) *Label {
	return &Label{text: text}
}

func (l *Label) Init() tea.Cmd {
	return nil
}

func (l *Label) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}

func (l *Label) View() string {
	return l.text
}
