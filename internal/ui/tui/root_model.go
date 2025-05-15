package tui

import tea "github.com/charmbracelet/bubbletea"

type RootModel struct {
	pageStack []tea.Model
}

func newRootModel(initialPage tea.Model) *RootModel {
	return &RootModel{
		pageStack: []tea.Model{initialPage},
	}
}

func (m *RootModel) currentPage() tea.Model {
	if len(m.pageStack) == 0 {
		return nil
	}

	return m.pageStack[len(m.pageStack)-1]
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
	case PushPageMsg:
		m.pushPage(msg.Page)

		return m, m.currentPage().Init()
	case PopPageMsg:
		if len(m.pageStack) > 1 {
			m.popPage()
		} else {
			return m, tea.Quit
		}
	}

	page, cmd := m.currentPage().Update(msg)

	m.pageStack[len(m.pageStack)-1] = page

	return m, cmd
}

func (m *RootModel) View() string {
	return m.currentPage().View()
}
