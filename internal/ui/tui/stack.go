package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	draw        func(...string) string
}

func newStack(orientation Orientation, gap int, components ...tea.Model) *Stack {
	return newStackWithPosition(lipgloss.Center, orientation, gap, components...)
}

func newStackWithPosition(pos lipgloss.Position, orientation Orientation, gap int, components ...tea.Model) *Stack {
	draw := func(components ...string) string {
		return lipgloss.JoinHorizontal(pos, components...)
	}

	delimiter := " "
	if orientation == Vertical {
		delimiter = "\n"
		draw = func(components ...string) string {
			return lipgloss.JoinVertical(pos, components...)
		}
	}
	gapDelimiter := strings.Repeat(delimiter, gap)

	return &Stack{
		components,
		orientation,
		gapDelimiter,
		draw,
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
	components := make([]string, 0, len(s.components)*2-1)
	for i, component := range s.components {
		components = append(components, component.View())
		if i < len(s.components)-1 {
			components = append(components, s.delimiter)
		}
	}

	return s.draw(components...)
}
