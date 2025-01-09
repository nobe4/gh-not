package repl

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
		return m.handleResult(msg)

	case m.command.Focused():
		return m.handleCommand(msg)

	case m.list.FilterState() == list.Filtering:
		return m.handleFiltering(msg)

	default:
		return m.handleBrowsing(msg)
	}
}

type CleanListMsg struct{}

func (CleanListMsg) apply(m model) (tea.Model, tea.Cmd) {
	items := []list.Item{}

	for _, e := range m.list.Items() {
		i, ok := e.(item)

		if !ok {
			continue
		}

		if !i.notification.Visible() {
			continue
		}

		items = append(items, e)
	}

	return m, tea.Sequence(
		// SetItems is needed here because the list might have less items now.
		m.list.SetItems(items),
		m.setIndexes(),
	)
}

func (m model) handleCommand(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keymap.CommandAccept):
		return m.acceptCommand()

	case key.Matches(msg, m.keymap.CommandCancel):
		return m.cancelCommand()
	}

	var cmd tea.Cmd
	m.command, cmd = m.command.Update(msg)

	return m, cmd
}

func (m model) handleResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.list.KeyMap.ShowFullHelp):
		m.showHelp = !m.showHelp
		slog.Debug("toggle help", "showAll", m.showHelp)

	case key.Matches(msg, m.list.KeyMap.Quit):
		m.showResult = false

		return m, nil
	}

	var cmd tea.Cmd
	m.result, cmd = m.result.Update(msg)

	return m, cmd
}

func (m model) handleBrowsing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	slog.Debug("browsing", "key", msg.String())

	switch {
	case key.Matches(msg, m.list.KeyMap.ShowFullHelp):
		m.showHelp = !m.showHelp
		slog.Debug("toggle help", "showAll", m.showHelp)

	case key.Matches(msg, m.keymap.Toggle):
		if i, ok := m.list.SelectedItem().(item); ok {
			i.selected = !i.selected
			slog.Debug("toggle selected", "item", i.notification.Subject.Title, "selected", i.selected)

			return m, m.list.SetItem(i.index, i)
		}

	case key.Matches(msg, m.keymap.All):
		return m.selectAll(true)

	case key.Matches(msg, m.keymap.None):
		return m.selectAll(false)

	case key.Matches(msg, m.keymap.Open):
		current, ok := m.list.SelectedItem().(item)
		if ok {
			err := m.actions["open"].Run(current.notification, nil, os.Stderr)
			if err != nil {
				slog.Warn("error opening", "notification", current.notification, "error", err)
			}
		}

	case key.Matches(msg, m.keymap.CommandMode):
		slog.Debug("focus command")

		return m, m.command.Focus()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) setIndexes() tea.Cmd {
	cmds := []tea.Cmd{}

	for index, e := range m.list.Items() {
		if i, ok := e.(item); ok {
			i.index = index
			cmds = append(cmds, m.list.SetItem(index, i))
		}
	}

	return tea.Batch(cmds...)
}

func (m model) handleFiltering(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	slog.Debug("filtering", "key", msg.String())

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}
