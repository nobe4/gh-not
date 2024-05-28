package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/notifications"
	"github.com/spf13/cobra"
)

var (
	replCmd = &cobra.Command{
		Use:     "repl",
		Aliases: []string{"r"},
		Short:   "Launch a REPL with notifications",
		RunE:    runRepl,
	}
)

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Toggle  key.Binding
	All     key.Binding
	Search  key.Binding
	Command key.Binding
	Help    key.Binding
	Quit    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Toggle, k.All},
		{k.Search, k.Command, k.Help, k.Quit},
	}
}

type Mode int64

const (
	Normal Mode = iota
	Search
	Command
	Result
	Help
)

type filteredList []int

type selection struct {
	id       int
	selected bool
}

type model struct {
	mode Mode
	keys keyMap
	help help.Model

	cursor         int
	choices        notifications.Notifications
	visibleChoices filteredList

	actors actors.ActorsMap

	renderCache []string
	selected    map[int]bool
	filter      textinput.Model
	command     textinput.Model
	result      string
}

func init() {
	rootCmd.AddCommand(replCmd)
}

func runRepl(cmd *cobra.Command, args []string) error {
	notifications, err := client.Notifications()
	if err != nil {
		slog.Error("Failed to list notifications", "err", err)
		return err
	}

	renderCache, err := notifications.ToTable()
	if err != nil {
		return err
	}

	model := model{
		mode: Normal,
		keys: keyMap{
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			Toggle: key.NewBinding(
				key.WithKeys("space", "enter"),
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
		actors:      actors.Map(client),
		choices:     notifications,
		selected:    map[int]bool{},
		renderCache: strings.Split(renderCache, "\n"),
	}

	model.filter = textinput.New()
	model.filter.Prompt = "/"

	model.command = textinput.New()
	model.command.Prompt = ":"

	suggestions := make([]string, 0, len(model.actors))
	for k := range model.actors {
		suggestions = append(suggestions, k)
	}
	model.command.SetSuggestions(suggestions)
	model.command.ShowSuggestions = true

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (m model) Init() tea.Cmd {
	return m.applyFilter()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case filteredList:
		m.visibleChoices = msg

	case result:
		m.result = msg.ToString()

	case selection:
		m.selected[msg.id] = msg.selected

	case tea.KeyMsg:
		switch m.mode {
		case Normal:
			switch {
			case key.Matches(msg, m.keys.Help):
				m.mode = Help
				m.help.ShowAll = !m.help.ShowAll

			case key.Matches(msg, m.keys.Search):
				m.mode = Search
				m.filter.Focus()

			case key.Matches(msg, m.keys.Command):
				m.mode = Command
				m.command.Focus()

			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}

			case key.Matches(msg, m.keys.Down):
				if m.cursor < len(m.visibleChoices)-1 {
					m.cursor++
				}

			case key.Matches(msg, m.keys.Toggle):
				return m, m.toggleSelect()

			case key.Matches(msg, m.keys.All):
				return m, m.selectAll()

			}

		case Search:
			switch msg.String() {
			case "esc":
				m.mode = Normal
				m.filter.SetValue("")
				m.filter.Blur()
			case "enter":
				m.mode = Normal
				m.filter.Blur()
			default:
				m.filter, _ = m.filter.Update(msg)
			}
			return m, m.applyFilter()

		case Command:
			switch msg.String() {
			case "esc":
				m.mode = Normal
				m.command.SetValue("")
				m.command.Blur()
			case "enter":
				command := m.command.Value()
				m.mode = Result
				m.command.SetValue("")
				m.command.Blur()
				return m, m.runCommand(command)
			default:
				m.command, _ = m.command.Update(msg)
			}

		case Result:
			m.mode = Normal

		case Help:
			m.mode = Normal
		}
	}

	return m, cmd
}

func (m model) View() string {
	if m.mode == Result {
		return m.result
	}

	if m.mode == Help {
		m.help.ShowAll = true
		return m.help.View(m.keys)
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

	switch m.mode {
	case Normal:
		out += m.help.View(m.keys)
	case Search:
		out += m.filter.View()
	case Command:
		out += m.command.View()
	}

	return out
}

type result struct {
	out string
	err error
}

func (r result) ToString() string {
	out := ""

	if r.err != nil {
		out = colors.Red(r.err.Error())
	} else {
		out = r.out
	}

	return out + "\npress any key to continue"
}

func (m model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		actor, ok := m.actors[command]
		if !ok {
			return result{
				out: "",
				err: fmt.Errorf("unknown command: %s", command),
			}
		}

		hasSelected := false
		out := ""
		for i, selected := range m.selected {
			if selected {
				hasSelected = true
				n, outn, err := actor.Run(m.choices[i])
				if err != nil {
					return result{err: err}
				}

				m.choices[i] = n
				out += outn + "\n"
			}
		}

		if !hasSelected {
			return result{err: fmt.Errorf("no notification selected")}
		}

		return result{out: out}
	}
}

func (m model) applyFilter() tea.Cmd {
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

func (m model) toggleSelect() tea.Cmd {
	return func() tea.Msg {
		visibleLineId := m.visibleChoices[m.cursor]
		selected, ok := m.selected[visibleLineId]

		return selection{
			id:       visibleLineId,
			selected: !(selected && ok),
		}
	}
}
func (m model) selectAll() tea.Cmd {
	cmds := tea.BatchMsg{}

	for _, id := range m.visibleChoices {
		cmds = append(cmds,
			func() tea.Msg {
				return selection{id: id, selected: true}
			},
		)
	}

	return tea.Batch(cmds...)
}
