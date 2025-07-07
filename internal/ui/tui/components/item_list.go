package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListOnAction func(item list.Item) tea.Cmd

type ItemList struct {
	list.Model

	focus bool

	onSelectAction ListOnAction
}

func NewItemList(items []list.Item) *ItemList {
	return &ItemList{
		Model: list.New(items, list.NewDefaultDelegate(), 0, 0),
		focus: true,
	}
}

func (il *ItemList) SetItems(items []list.Item) {
	il.Model.SetItems(items)
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
	cmds := []tea.Cmd{cmd}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			cmds = append(cmds, il.onSelectAction(il.SelectedItem()))
		}
	}

	return il, tea.Batch(cmds...)
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

func (il *ItemList) SetSize(width, height int) {
	il.Model.SetSize(width, height)

	// Workaround for list width issue - https://github.com/charmbracelet/bubbles/issues/744
	il.Model.Styles.HelpStyle = il.Model.Styles.HelpStyle.Width(width - 2)
}

func (il *ItemList) OnSelect(action ListOnAction) {
	il.onSelectAction = action
}
