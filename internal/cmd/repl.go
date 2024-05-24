package cmd

import (
	"fmt"
	"log/slog"
	"strings"

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

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

type model struct {
	cursor      int
	choices     notifications.Notifications
	renderCache []string
	selected    map[int]bool
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("notification list")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
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

	out += "press ? for help"

	return out
}
