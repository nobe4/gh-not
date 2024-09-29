package repl

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) selectAll(selected bool) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	for _, e := range m.list.VisibleItems() {
		if i, ok := e.(item); ok {
			i.selected = selected
			cmds = append(cmds, m.list.SetItem(i.index, i))
		}
	}

	return m, tea.Sequence(cmds...)
}
