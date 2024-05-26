package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/actors"
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

type Mode int64

const (
	Normal Mode = iota
	Search
	Command
)

type filteredList []int

type selection struct {
	id       int
	selected bool
}

type model struct {
	mode Mode

	cursor         int
	choices        notifications.Notifications
	visibleChoices filteredList

	actors actors.ActorsMap

	renderCache []string
	selected    map[int]bool
	filter      textinput.Model
	command     textinput.Model
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

	case selection:
		m.selected[msg.id] = msg.selected

	case tea.KeyMsg:
		switch m.mode {
		case Normal:
			switch msg.String() {
			case "?":
				panic("to implement with the help bubble")

			case "/":
				m.mode = Search
				m.filter.Focus()

			case ":":
				m.mode = Command
				m.command.Focus()

			case "esc":
				return m, tea.Quit

			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.visibleChoices)-1 {
					m.cursor++
				}

			case "enter", " ":
				return m, m.toggleSelect()
			case "a":
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
				return m, m.runCommand(false)
			case "enter":
				return m, m.runCommand(true)
			default:
				m.command, cmd = m.command.Update(msg)
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
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
		out += "press ? for help\n"
	case Search:
		out += m.filter.View()
	case Command:
		out += m.command.View()
	}

	return out
}

func (m model) runCommand(apply bool) tea.Cmd {
	command := m.command.Value()

	m.mode = Normal
	m.command.SetValue("")
	m.command.Blur()

	return func() tea.Msg {

		if _, ok := m.actors[command]; ok {
			if apply {
				return tea.Quit()
			}
		} else {
			panic(command)
		}

		return nil
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
