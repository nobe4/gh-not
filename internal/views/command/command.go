package command

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/notifications"
	"github.com/nobe4/gh-not/internal/views"
)

type keymap struct {
	quit    key.Binding
	confirm key.Binding
}

type Model struct {
	keys                  keymap
	input                 textinput.Model
	actors                actors.ActorsMap
	selectedNotifications func(func(*notifications.Notification))
}

func New(actors actors.ActorsMap, selectedNotifications func(func(*notifications.Notification))) Model {
	model := Model{
		keys: keymap{
			quit: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("ESC", "quit"),
			),
			confirm: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "confirm command"),
			),
		},
		input:                 textinput.New(),
		actors:                actors,
		selectedNotifications: selectedNotifications,
	}

	suggestions := make([]string, 0, len(actors))
	for k := range actors {
		suggestions = append(suggestions, k)
	}

	model.input.Prompt = ":"
	model.input.SetSuggestions(suggestions)
	model.input.ShowSuggestions = true

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.input.Focus()

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.quit):
			m.input.SetValue("")
			m.input.Blur()
			return m, views.ChangeMode(views.NormalMode)

		case key.Matches(msg, m.keys.confirm):
			input := m.input.Value()
			m.input.SetValue("")
			m.input.Blur()
			return m, tea.Sequence(m.runCommand(input), views.ChangeMode(views.NormalMode))

		default:
			m.input, _ = m.input.Update(msg)
		}
	}

	return m, cmd
}

func (m Model) View() string {
	return m.input.View()
}

func (m Model) runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		actor, ok := m.actors[command]
		if !ok {
			return views.ResultMsg{
				Err: fmt.Errorf("unknown command: %s", command),
			}
		}

		result := views.ResultMsg{}

		hasSelected := false

		m.selectedNotifications(func(n *notifications.Notification) {
			hasSelected = true
			buff := &strings.Builder{}

			err := actor.Run(n, buff)
			if err != nil {
				result.Err = err
			}
			fmt.Fprintln(buff, "")

			result.Out += buff.String()
		})

		if !hasSelected {
			return views.ResultMsg{Err: fmt.Errorf("no notification selected")}
		}

		return result
	}
}
