package normal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
	"github.com/nobe4/gh-not/internal/views"
	"github.com/nobe4/gh-not/internal/views/command"
	"github.com/nobe4/gh-not/internal/views/filter"
)

type keymap struct {
	up      key.Binding
	down    key.Binding
	toggle  key.Binding
	all     key.Binding
	unall   key.Binding
	search  key.Binding
	command key.Binding
	help    key.Binding
	quit    key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.help}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down, k.toggle, k.all, k.unall},
		{k.search, k.command, k.help, k.quit},
	}
}

type SelectMsg struct {
	id       int
	selected bool
}

type Model struct {
	Mode views.Mode
	Keys keymap
	help help.Model

	cursor         int
	choices        notifications.Notifications
	visibleChoices []int

	actors actors.ActorsMap

	renderCache []string
	selected    map[int]bool

	filter  tea.Model
	command tea.Model

	result string
}

func New(client *gh.Client, notifications notifications.Notifications, renderCache string) Model {
	model := Model{
		Mode: views.NormalMode,
		Keys: keymap{
			up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			toggle: key.NewBinding(
				key.WithKeys(" ", "enter"),
				key.WithHelp("space/enter", "toggle selected"),
			),
			all: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "select all"),
			),
			unall: key.NewBinding(
				key.WithKeys("A"),
				key.WithHelp("A", "unselect all"),
			),
			search: key.NewBinding(
				key.WithKeys("/"),
				key.WithHelp("/", "search mode"),
			),
			command: key.NewBinding(
				key.WithKeys(":"),
				key.WithHelp(":", "command mode"),
			),
			help: key.NewBinding(
				key.WithKeys("?"),
				key.WithHelp("?", "toggle help"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "esc", "ctrl+c"),
				key.WithHelp("q/esc/ctrl+c", "quit"),
			),
		},
		help:        help.New(),
		cursor:      0,
		choices:     notifications,
		selected:    map[int]bool{},
		renderCache: strings.Split(renderCache, "\n"),
	}

	model.command = command.New(actors.Map(client), model.SelectedNotificationsFunc)
	model.filter = filter.New(model.VisibleLinesFunc)

	return model
}

func (m Model) Init() tea.Cmd {
	return m.filter.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case filter.FilterMsg:
		m.visibleChoices = msg.IntSlice()

	case views.ResultMsg:
		m.result = msg.ToString()
		return m, views.ChangeMode(views.ResultMode)

	case SelectMsg:
		m.selected[msg.id] = msg.selected

	case views.ChangeModeMsg:
		m.Mode = msg.Mode

	case tea.KeyMsg:
		switch m.Mode {
		case views.NormalMode:
			switch {
			case key.Matches(msg, m.Keys.help):
				return m, views.ChangeMode(views.HelpMode)

			case key.Matches(msg, m.Keys.search):
				return m, views.ChangeMode(views.SearchMode)

			case key.Matches(msg, m.Keys.command):
				return m, views.ChangeMode(views.CommandMode)

			case key.Matches(msg, m.Keys.quit):
				return m, tea.Quit

			case key.Matches(msg, m.Keys.up):
				if m.cursor > 0 {
					m.cursor--
				}

			case key.Matches(msg, m.Keys.down):
				if m.cursor < len(m.visibleChoices)-1 {
					m.cursor++
				}

			case key.Matches(msg, m.Keys.toggle):
				return m, m.toggleSelect()

			case key.Matches(msg, m.Keys.all):
				return m, m.selectAll(true)

			case key.Matches(msg, m.Keys.unall):
				return m, m.selectAll(false)

			}

		case views.SearchMode:
			m.filter, cmd = m.filter.Update(msg)
			cmds = append(cmds, cmd)

		case views.CommandMode:
			m.command, cmd = m.command.Update(msg)
			cmds = append(cmds, cmd)

		case views.ResultMode:
			return m, views.ChangeMode(views.NormalMode)

		case views.HelpMode:
			return m, views.ChangeMode(views.NormalMode)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Mode == views.ResultMode {
		return m.result
	}

	if m.Mode == views.HelpMode {
		m.help.ShowAll = true
		return m.help.View(m.Keys)
	} else {
		m.help.ShowAll = false
	}

	out := ""

	for i, id := range m.visibleChoices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if v, ok := m.selected[id]; ok && v {
			checked = "x"
		}

		out += fmt.Sprintf("%s%s%s\n", checked, cursor, m.renderCache[id])
	}

	switch m.Mode {
	case views.NormalMode:
		out += m.help.View(m.Keys)
	case views.SearchMode:
		out += m.filter.View()
	case views.CommandMode:
		out += m.command.View()
	}

	return out
}

func (m Model) SelectedNotificationsFunc(cb func(*notifications.Notification)) {
	for i, selected := range m.selected {
		if selected {
			cb(m.choices[i])
		}
	}
}

func (m Model) VisibleLinesFunc(cb func(string, int)) {
	for i, line := range m.renderCache {
		cb(line, i)
	}
}

func (m Model) toggleSelect() tea.Cmd {
	return func() tea.Msg {
		visibleLineId := m.visibleChoices[m.cursor]
		selected, ok := m.selected[visibleLineId]

		return SelectMsg{
			id:       visibleLineId,
			selected: !(selected && ok),
		}
	}
}
func (m Model) selectAll(selected bool) tea.Cmd {
	cmds := tea.BatchMsg{}

	for _, id := range m.visibleChoices {
		cmds = append(
			cmds,
			func() tea.Msg {
				return SelectMsg{id: id, selected: selected}
			},
		)
	}

	return tea.Batch(cmds...)
}
