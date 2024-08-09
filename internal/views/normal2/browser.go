package normal2

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) selectAll(selected bool) (tea.Model, tea.Cmd) {
	items := []list.Item{}
	for _, e := range m.list.Items() {
		if i, ok := e.(item); ok {
			i.selected = selected
			items = append(items, i)
		}
	}
	return m, m.list.SetItems(items)
}
