package normal2

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case m.showResult:
		return m, m.handleResult(msg)
	case m.command.Focused():
		return m, m.handleCommand(msg)
	case m.list.FilterState() == list.Filtering:
		return m, m.handleFiltering(msg)
	default:
		return m, m.handleBrowsing(msg)
	}
}

type CleanListMsg struct{}

func (_ CleanListMsg) apply(m model) (tea.Model, tea.Cmd) {
	items := []list.Item{}
	for _, e := range m.list.Items() {
		if i, ok := e.(item); ok {
			if !i.notification.Meta.Done {
				items = append(items, e)
			}
		}
	}

	return m, m.list.SetItems(items)
}

func (m *model) handleCommand(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keymap.CommandAccept):
		// TODO: move into a method
		slog.Debug("blur command")
		command := m.command.Value()
		m.command.SetValue("")
		m.command.Blur()
		m.showResult = true
		return m.applyCommand(command)

	case key.Matches(msg, m.keymap.CommandCancel):
		// TODO: move into a method
		slog.Debug("blur command")
		m.command.SetValue("")
		m.command.Blur()
		return nil
	}

	m.command, cmd = m.command.Update(msg)
	return cmd
}

func (m *model) handleResult(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {

	case key.Matches(msg, m.list.KeyMap.ShowFullHelp):
		m.help.ShowAll = !m.help.ShowAll
		// TODO: why is this not showing the full help?
		slog.Debug("toggle help", "showAll", m.help.ShowAll)

	case key.Matches(msg, m.list.KeyMap.Quit):
		m.showResult = false
		return nil
	}

	m.result, cmd = m.result.Update(msg)
	return cmd
}

func (m *model) handleBrowsing(msg tea.KeyMsg) tea.Cmd {
	slog.Debug("browsing", "key", msg.String())

	switch {
	case key.Matches(msg, m.list.KeyMap.ShowFullHelp):
		m.help.ShowAll = !m.help.ShowAll
		slog.Debug("toggle help", "showAll", m.help.ShowAll)

	case key.Matches(msg, m.keymap.Toggle):
		if i, ok := m.list.SelectedItem().(item); ok {
			i.selected = !i.selected
			slog.Debug("toggle selected", "item", i.notification.Subject.Title, "selected", i.selected)

			return m.list.SetItem(m.list.GlobalIndex(), i)
		}

	case key.Matches(msg, m.keymap.All):
		// TODO: move into a method
		items := []list.Item{}
		for _, e := range m.list.Items() {
			if i, ok := e.(item); ok {
				i.selected = true
				items = append(items, i)
			}
		}
		m.list.SetItems(items)

	case key.Matches(msg, m.keymap.None):
		// TODO: move into a method
		items := []list.Item{}
		for _, e := range m.list.Items() {
			if i, ok := e.(item); ok {
				i.selected = false
				items = append(items, i)
			}
		}
		m.list.SetItems(items)

	case key.Matches(msg, m.keymap.Open):
		// TODO: move into a method
		current, ok := m.list.SelectedItem().(item)
		if ok {
			m.actors["open"].Run(current.notification, os.Stderr)
		}

	case key.Matches(msg, m.keymap.CommandMode):
		slog.Debug("focus command")
		m.command.Focus()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *model) handleFiltering(msg tea.KeyMsg) tea.Cmd {
	slog.Debug("filtering", "key", msg.String())
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd
}
