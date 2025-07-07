package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
	"github.com/hop-/gotchat/internal/ui/tui/commands"
)

type RootModel struct {
	pageStack []tea.Model
	emitter   core.EventEmitter
}

func newRootModel(initialPage tea.Model, emitter core.EventEmitter) *RootModel {
	return &RootModel{
		[]tea.Model{initialPage},
		emitter,
	}
}

func (m *RootModel) currentPage() tea.Model {
	if len(m.pageStack) == 0 {
		return nil
	}

	return m.pageStack[len(m.pageStack)-1]
}

func (m *RootModel) updateCurrnetPage(page tea.Model) {
	if len(m.pageStack) == 0 {
		return
	}

	m.pageStack[len(m.pageStack)-1] = page
}

func (m *RootModel) pushPage(page tea.Model) {
	m.pageStack = append(m.pageStack, page)
}

func (m *RootModel) popPage() {
	m.pageStack = m.pageStack[:len(m.pageStack)-1]
}

func (m *RootModel) Init() tea.Cmd {
	return m.currentPage().Init()
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case commands.SetNewPageMsg:
		// Reset the page stack with the new page
		m.pageStack = []tea.Model{msg.Page}

		return m, m.currentPage().Init()
	case commands.PushPageMsg:
		m.pushPage(msg.Page)

		return m, m.currentPage().Init()
	case commands.PopPageMsg:
		if len(m.pageStack) > 1 {
			m.popPage()
		} else {
			return m, commands.Shutdown
		}
	case commands.ShutdownMsg:
		// Setup shutdown screen
		m.pageStack = []tea.Model{newShutdownModel()}

		return m, m.currentPage().Init()
	case commands.InternalQuitMsg:
		m.emitter.Emit(core.QuitEvent{})
	}

	page, cmd := m.currentPage().Update(msg)

	m.updateCurrnetPage(page)

	return m, cmd
}

func (m *RootModel) View() string {
	return m.currentPage().View()
}
