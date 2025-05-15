package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemList struct {
	list.Model

	focus bool
}

func newItemList(items []list.Item) *ItemList {
	return &ItemList{
		Model: list.New(items, list.NewDefaultDelegate(), 0, 0),
		focus: true,
	}
}

func (il *ItemList) Init() tea.Cmd {
	return nil
}

func (il *ItemList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !il.focus {
		return il, nil
	}

	var cmd tea.Cmd
	il.Model, cmd = il.Model.Update(msg)

	return il, cmd
}

func (il *ItemList) View() string {
	return il.Model.View()
}

func (il *ItemList) Focus() tea.Cmd {
	il.focus = true

	// TODO: Handle focus for list
	return nil
}

func (il *ItemList) Blur() tea.Cmd {
	il.focus = false

	// TODO: Handle blur for list
	return nil
}

func (il *ItemList) Focused() bool {
	return il.focus
}
