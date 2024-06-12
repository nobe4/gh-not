package filter

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/views"
)

type Keymap struct {
	cancel  key.Binding
	confirm key.Binding
}

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.cancel}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.cancel, k.confirm},
	}
}

type FilterMsg []int

type Model struct {
	Keys         Keymap
	input        textinput.Model
	visibleLines func(func(string, int))
}

func New(visibleLines func(func(string, int)), keymap config.Keymap) Model {
	model := Model{
		Keys: Keymap{
			cancel:  keymap["filter"]["cancel"].Binding("cancel"),
			confirm: keymap["filter"]["confirm"].Binding("confirm"),
		},
		input:        textinput.New(),
		visibleLines: visibleLines,
	}

	model.input.Prompt = keymap["normal"]["filter"][0]
	model.input.Placeholder = "filter"

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

		case key.Matches(msg, m.Keys.cancel):
			m.input.SetValue("")
			m.input.Blur()
			cmds = append(cmds, views.ChangeMode(views.NormalMode))

		case key.Matches(msg, m.Keys.confirm):
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
