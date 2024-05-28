package filter

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/views"
)

type keymap struct {
	quit    key.Binding
	confirm key.Binding
}

type FilterMsg []int

type Model struct {
	keys         keymap
	input        textinput.Model
	visibleLines func(func(string, int))
}

func New(visibleLines func(func(string, int))) Model {
	model := Model{
		keys: keymap{
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl-c"),
				key.WithHelp("esc/ctrl-c", "cancel"),
			),
			confirm: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "confirm"),
			),
		},
		input:        textinput.New(),
		visibleLines: visibleLines,
	}

	model.input.Prompt = "/"

	return model
}

func (m Model) Init() tea.Cmd {
	return m.filter("")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.input.Focus()

	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.quit):
			m.input.SetValue("")
			m.input.Blur()
			cmds = append(cmds, views.ChangeMode(views.NormalMode))

		case key.Matches(msg, m.keys.confirm):
			m.input.Blur()
			cmds = append(cmds, views.ChangeMode(views.NormalMode))

		default:
			m.input, _ = m.input.Update(msg)
		}
	}

	cmds = append(cmds, m.filter(m.input.Value()))

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.input.View()
}

func (m Model) filter(f string) tea.Cmd {
	return func() tea.Msg {
		visibleChoices := FilterMsg{}

		m.visibleLines(func(line string, id int) {
			if f == "" || strings.Contains(line, f) {
				visibleChoices = append(visibleChoices, id)
			}
		})

		return visibleChoices
	}
}

func (f FilterMsg) IntSlice() []int {
	i := []int{}
	for _, j := range f {
		i = append(i, j)
	}
	return i
}
