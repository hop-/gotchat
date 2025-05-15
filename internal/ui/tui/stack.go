package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Orientation int

const (
	Vertical Orientation = iota
	Horizontal
)

type Stack struct {
	components  []tea.Model
	orientation Orientation
	delimiter   string
}

func newStack(orientation Orientation, gap int, components ...tea.Model) *Stack {
	delimiter := " "
	if orientation == Vertical {
		delimiter = "\n"
	}
	gapDelimiter := strings.Repeat(delimiter, gap)

	return &Stack{
		components,
		orientation,
		gapDelimiter,
	}
}

func (s *Stack) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(s.components))
	for i := 0; i <= len(s.components)-1; i++ {
		cmds[i] = s.components[i].Init()
	}

	return tea.Batch(cmds...)
}

func (s *Stack) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s *Stack) View() string {
	components := make([]string, len(s.components))
	for i, component := range s.components {
		components[i] = component.View()
	}

	return strings.Join(components, s.delimiter)
}
