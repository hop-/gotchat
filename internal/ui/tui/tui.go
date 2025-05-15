package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Tui struct {
	p *tea.Program
}

func New() *Tui {
	rootModel := newRootModel(newUsersListModel())
	//rootModel := newLoginModel()
	p := tea.NewProgram(rootModel)

	return &Tui{p}
}

func (ui *Tui) Init() error {
	// Not needed for Bubble Tea
	return nil
}

func (ui *Tui) Run() error {
	_, err := ui.p.Run()

	return err
}

func (ui *Tui) Close() error {
	// TODO: Send close signal to the program
	return nil
}
