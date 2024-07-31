package normal2

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/notifications"
)

type model struct {
	list   list.Model
	keymap Keymap
	help   help.Model
}

func Init(n notifications.Notifications, keymap config.Keymap, view config.View) error {
	items := []list.Item{}

	for _, notification := range n {
		items = append(items, item{notification: notification})
	}

	l := list.New(items, itemDelegate{}, 0, view.Height)

	m := model{list: l}

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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			return m, m.handleFiltering(msg)
		} else {
			return m, m.handleBrowsing(msg)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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
		for _, i := range m.list.Items() {
			if i, ok := i.(item); ok && i.selected {
				slog.Debug("selected", "notification", i.notification.String())
			}
		}
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
