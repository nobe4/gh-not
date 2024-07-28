package normal2

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/notifications"
)

const listHeight = 10
const listDefaultWidth = 20

type item struct {
	*notifications.Notification
}

func (i item) FilterValue() string { return i.String() }

type model struct {
	list   list.Model
	choice *item
}

func (m model) Init() tea.Cmd {
	return nil
}

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

	switch keystroke := msg.String(); keystroke {
	case "?":
		m.list.SetShowHelp(!m.list.ShowHelp())

	case "enter":
		if i, ok := m.list.SelectedItem().(item); ok {
			m.choice = &i
		}
		return tea.Quit
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

func Init(n notifications.Notifications) {
	items := []list.Item{}

	for _, notification := range n {
		items = append(items, item{notification})
	}

	l := list.New(items, itemDelegate{}, listDefaultWidth, listHeight)
	m := model{list: l}
	m.initView()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
