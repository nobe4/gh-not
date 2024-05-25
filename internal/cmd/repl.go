package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
		choices:     notifications,
		selected:    map[int]bool{},
		renderCache: strings.Split(renderCache, "\n"),
	}

	model.filter = textinput.New()
	model.filter.Prompt = "/"

	model.command = textinput.New()
	model.command.Prompt = ":"

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

type Mode int64

const (
	Normal Mode = iota
	Search
	Command
)

type model struct {
	mode Mode

	cursor      int
	choices     notifications.Notifications
	renderCache []string
	selected    map[int]bool
	filter      textinput.Model
	command     textinput.Model
}

func (m model) Init() tea.Cmd {

	return tea.SetWindowTitle("notification list")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case Normal:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "?":
				panic("to implement")

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
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter", " ":
				v, ok := m.selected[m.cursor]
				if ok && v {
					m.selected[m.cursor] = false
				} else {
					m.selected[m.cursor] = true
				}
			}
		}
	case Search:
		switch msg := msg.(type) {
		case tea.KeyMsg:
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
		}
	case Command:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.mode = Normal
				m.command.Blur()
			case "enter":
				m.mode = Normal
				m.command.Blur()
			default:
				m.command, _ = m.command.Update(msg)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	out := ""

	for i := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if v, ok := m.selected[i]; ok && v {
			checked = "x"
		}

		out += fmt.Sprintf("%s%s%s\n", checked, cursor, m.renderCache[i])
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
