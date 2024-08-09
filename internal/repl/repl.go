package repl

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/notifications"
)

type model struct {
	keymap       Keymap
	actors       actors.ActorsMap
	currentActor actors.Actor

	showHelp bool
	list     list.Model
	command  textinput.Model

	result        viewport.Model
	resultStrings []string

	ready        bool
	showResult   bool
	processQueue []item
	maxHeigth    int
}

func Init(n notifications.Notifications, actors actors.ActorsMap, keymap config.Keymap, view config.View) error {
	items := []list.Item{}

	for _, notification := range n {
		items = append(items, item{notification: notification})
	}

	m := model{
		list:      list.New(items, itemDelegate{}, 0, 0),
		command:   textinput.New(),
		actors:    actors,
		result:    viewport.New(0, 0),
		maxHeigth: view.Height,
	}

	m.list.SetItems(items)
	m.initView()
	m.initKeymap(keymap)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}

	return nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("update", "msg", msg)
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return m.handleResize(msg)

	case ResultUpdateMsg:
		m.result.SetContent(msg.content)

	case ApplyCommandMsg:
		return msg.apply(m)

	case AppliedCommandMsg:
		return msg.apply(m)

	case CleanListMsg:
		return msg.apply(m)

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
