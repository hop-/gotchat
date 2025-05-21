package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hop-/gotchat/internal/core"
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
	case SetNewPageMsg:
		m.pageStack = []tea.Model{msg.Page}

		return m, m.currentPage().Init()
	case PushPageMsg:
		m.pushPage(msg.Page)

		return m, m.currentPage().Init()
	case PopPageMsg:
		if len(m.pageStack) > 1 {
			m.popPage()
		} else {
			return m, internalQuit
		}
	case InternalQuitMsg:
		m.emitter.Emit(core.QuitEvent{})
	}

	page, cmd := m.currentPage().Update(msg)

	m.updateCurrnetPage(page)

	return m, cmd
}

func (m *RootModel) View() string {
	return m.currentPage().View()
}
