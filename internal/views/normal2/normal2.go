package normal2

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/notifications"
)

type model struct {
	keymap Keymap
	actors actors.ActorsMap

	help    help.Model
	list    list.Model
	command textinput.Model
	result  viewport.Model

	ready     bool
	maxHeigth int
}

func Init(n notifications.Notifications, actors actors.ActorsMap, keymap config.Keymap, view config.View) error {
	items := []list.Item{}

	for _, notification := range n {
		items = append(items, item{notification: notification})
	}

	m := model{
		list:      list.New(items, itemDelegate{}, 0, 0),
		command:   textinput.New(),
		actors:    actors,
		result:    viewport.New(0, 0),
		maxHeigth: view.Height,
	}

	m.list.SetItems(items)
	m.initView()
	m.initKeymap(keymap)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}

	return nil
}

type item struct {
	notification *notifications.Notification
	selected     bool
}

func (i item) FilterValue() string { return i.notification.String() }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.handleResize(msg)
		return m, nil

	case tea.KeyMsg:
		if m.command.Focused() {
			return m, m.handleCommand(msg)
		}

		if m.list.FilterState() == list.Filtering {
			return m, m.handleFiltering(msg)
		}

		return m, m.handleBrowsing(msg)
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m *model) handleCommand(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keymap.CommandAccept):
		panic(m.command.Value())

	case key.Matches(msg, m.keymap.CommandCancel):
		m.command.SetValue("")
		m.command.Blur()
		return nil
	}

	m.command, cmd = m.command.Update(msg)
	return cmd
}

func (m *model) handleBrowsing(msg tea.KeyMsg) tea.Cmd {
	slog.Debug("browsing", "key", msg.String())

	switch {
	case key.Matches(msg, m.list.KeyMap.ShowFullHelp):
		m.help.ShowAll = !m.help.ShowAll

	case key.Matches(msg, m.keymap.Toggle):
		if i, ok := m.list.SelectedItem().(item); ok {
			i.selected = !i.selected
			return m.list.SetItem(m.list.GlobalIndex(), i)
		}

	case key.Matches(msg, m.keymap.Test):
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
