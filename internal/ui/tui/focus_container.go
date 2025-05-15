package tui

import tea "github.com/charmbracelet/bubbletea"

type FocusContainer struct {
	items      []FocusableModel
	focusIndex int
}

func (fc *FocusContainer) Init() tea.Cmd {
	_, cmd := fc.ChangeFocus(0, false)

	return cmd
}

func (fc *FocusContainer) Update(msg tea.Msg) (*FocusContainer, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			_, cmd := fc.ChangeFocus(1, true)

			return fc, cmd
		case "enter":
			// TODO: Handle enter key for different components
		case "shift+tab":
			_, cmd := fc.ChangeFocus(-1, true)

			return fc, cmd
		}
	}

	cmd := fc.updateFocusedComponent(msg)

	return fc, cmd
}

func (fc *FocusContainer) ChangeFocus(val int, loop bool) (bool, tea.Cmd) {
	status := false

	if val != 0 {
		focusIndex := fc.focusIndex
		var findNextFocusIndex func()

		findNextFocusIndex = func() {
			focusIndex += val

			if focusIndex < 0 {
				if !loop {
					return
				}

				focusIndex = len(fc.items) - 1
			} else if focusIndex >= len(fc.items) {
				if !loop {
					return
				}

				focusIndex = 0
			}

			switch fc.items[focusIndex].(type) {
			case Activatable:
				if fc.items[focusIndex].(Activatable).IsActive() {
					status = true
					fc.focusIndex = focusIndex
				}
			default:
				status = true
				fc.focusIndex = focusIndex
			}

			if fc.focusIndex == focusIndex {
				status = false

				return
			}

			findNextFocusIndex()
		}

		findNextFocusIndex()
	}

	cmds := make([]tea.Cmd, len(fc.items))
	for i := 0; i <= len(fc.items)-1; i++ {
		if i == fc.focusIndex {
			cmds[i] = fc.items[i].Focus()
			continue
		}

		cmds[i] = fc.items[i].Blur()
	}

	return status, tea.Batch(cmds...)
}

func (fc *FocusContainer) updateFocusedComponent(msg tea.Msg) tea.Cmd {
	// TODO: update the model with the returned model
	_, cmd := fc.items[fc.focusIndex].Update(msg)
	// TODO: Use lipgloss.JoinHorizontal and lipgloss.JoinVertical to join the components

	return cmd
}
