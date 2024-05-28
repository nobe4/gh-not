package normal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
	"github.com/nobe4/gh-not/internal/views"
	"github.com/nobe4/gh-not/internal/views/command"
)

type Keymap struct {
	Up      key.Binding
	Down    key.Binding
	Toggle  key.Binding
	All     key.Binding
	Search  key.Binding
	Command key.Binding
	Help    key.Binding
	Quit    key.Binding
}

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Toggle, k.All},
		{k.Search, k.Command, k.Help, k.Quit},
	}
}

type filteredList []int

type SelectMsg struct {
	id       int
	selected bool
}

type Model struct {
	Mode views.Mode
	Keys Keymap
	help help.Model

	cursor         int
	choices        notifications.Notifications
	visibleChoices filteredList

	actors actors.ActorsMap

	renderCache []string
	selected    map[int]bool
	filter      textinput.Model
	command     tea.Model
	result      string
}

func New(client *gh.Client, notifications notifications.Notifications, renderCache string) Model {
	model := Model{
		Mode: views.NormalMode,
		Keys: Keymap{
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			Toggle: key.NewBinding(
				key.WithKeys(" ", "enter"),
				key.WithHelp("space/enter", "toggle selected"),
			),
			All: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "select all"),
			),
			Search: key.NewBinding(
				key.WithKeys("/"),
				key.WithHelp("/", "search mode"),
			),
			Command: key.NewBinding(
				key.WithKeys(":"),
				key.WithHelp(":", "command mode"),
			),
			Help: key.NewBinding(
				key.WithKeys("?"),
				key.WithHelp("?", "toggle help"),
			),
			Quit: key.NewBinding(
				key.WithKeys("q", "esc", "ctrl+c"),
				key.WithHelp("q/ESC/C-c", "quit"),
			),
		},
		help:        help.New(),
		cursor:      0,
		choices:     notifications,
		selected:    map[int]bool{},
		renderCache: strings.Split(renderCache, "\n"),
	}

	model.filter = textinput.New()
	model.filter.Prompt = "/"

	model.command = command.New(actors.Map(client), model.SelectedNotifications)

	return model
}

func (m Model) Init() tea.Cmd {
	return m.applyFilter()
}

func (m Model) SelectedNotifications(cb func(notifications.Notification)) {
	for i, selected := range m.selected {
		if selected {
			cb(m.choices[i])
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case filteredList:
		m.visibleChoices = msg

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
			case key.Matches(msg, m.Keys.Help):
				return m, views.ChangeMode(views.HelpMode)

			case key.Matches(msg, m.Keys.Search):
				m.filter.Focus()
				return m, views.ChangeMode(views.SearchMode)

			case key.Matches(msg, m.Keys.Command):
				return m, views.ChangeMode(views.CommandMode)

			case key.Matches(msg, m.Keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.Keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}

			case key.Matches(msg, m.Keys.Down):
				if m.cursor < len(m.visibleChoices)-1 {
					m.cursor++
				}

			case key.Matches(msg, m.Keys.Toggle):
				return m, m.toggleSelect()

			case key.Matches(msg, m.Keys.All):
				return m, m.selectAll()

			}

		case views.SearchMode:
			switch msg.String() {
			case "esc":
				m.filter.SetValue("")
				m.filter.Blur()
				return m, tea.Sequence(m.applyFilter(), views.ChangeMode(views.NormalMode))
			case "enter":
				m.filter.Blur()
				return m, tea.Sequence(m.applyFilter(), views.ChangeMode(views.NormalMode))
			default:
				m.filter, _ = m.filter.Update(msg)
			}
			return m, m.applyFilter()

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

func (m Model) applyFilter() tea.Cmd {
	return func() tea.Msg {
		m.cursor = 0
		f := m.filter.Value()

		visibleChoices := filteredList{}

		for i, line := range m.renderCache {
			if f == "" || strings.Contains(line, f) {
				visibleChoices = append(visibleChoices, i)
			}
		}

		return visibleChoices
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
func (m Model) selectAll() tea.Cmd {
	cmds := tea.BatchMsg{}

	for _, id := range m.visibleChoices {
		cmds = append(cmds,
			func() tea.Msg {
				return SelectMsg{id: id, selected: true}
			},
		)
	}

	return tea.Batch(cmds...)
}