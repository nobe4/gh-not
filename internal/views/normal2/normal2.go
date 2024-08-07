package normal2

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/notifications"
)

type model struct {
	keymap       Keymap
	actors       actors.ActorsMap
	currentActor actors.Actor

	help    help.Model
	list    list.Model
	command textinput.Model

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

type item struct {
	notification *notifications.Notification
	selected     bool
}

func (i item) FilterValue() string { return i.notification.String() }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.handleResize(msg)
		return m, nil

	case ApplyCommandMsg:
		slog.Debug("apply command", "command", msg.Command)
		actor, ok := m.actors[msg.Command]
		if !ok {
			m.result.SetContent(fmt.Sprintf("Invalid command %s\nPress %s to continue ...", msg.Command, m.list.KeyMap.Quit.Keys()))
			return m, nil
		}

		m.resultStrings = []string{}
		m.currentActor = actor
		m.processQueue = msg.Items

		// TODO: the rendering should be done only in one place
		if len(m.processQueue) == 0 {
			m.result.SetContent(fmt.Sprintf("done, press %s to continue ...", m.list.KeyMap.Quit.Keys()))
		} else {
			m.result.SetContent(fmt.Sprintf("%d more ...", len(m.processQueue)))
		}
		return m, m.applyNext()

	case AppliedCommandMsg:
		slog.Debug("applied command", "message", msg.Message)
		m.processQueue = m.processQueue[1:]

		m.resultStrings = append(m.resultStrings, msg.Message)

		content := lipgloss.JoinVertical(lipgloss.Top, m.resultStrings...)
		if len(m.processQueue) == 0 {
			content = lipgloss.JoinVertical(lipgloss.Top, content, fmt.Sprintf("done, press %s to continue ...", m.list.KeyMap.Quit.Keys()))
		} else {
			content = lipgloss.JoinVertical(lipgloss.Top, content, fmt.Sprintf("%d more ...", len(m.processQueue)))
		}

		m.result.SetContent(content)

		return m, m.applyNext()

	case CleanListMsg:
		items := []list.Item{}
		for _, e := range m.list.Items() {
			if i, ok := e.(item); ok {
				if !i.notification.Meta.Done {
					items = append(items, e)
				}
			}
		}
		m.list.SetItems(items)

	case tea.KeyMsg:
		if m.showResult {
			return m, m.handleResult(msg)
		}

		if m.command.Focused() {
			return m, m.handleCommand(msg)
		}

		if m.list.FilterState() == list.Filtering {
			return m, m.handleFiltering(msg)
		}

		return m, m.handleBrowsing(msg)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) handleCommand(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keymap.CommandAccept):
		slog.Debug("blur command")
		command := m.command.Value()
		m.command.SetValue("")
		m.command.Blur()
		m.showResult = true
		return m.applyCommand(command)

	case key.Matches(msg, m.keymap.CommandCancel):
		slog.Debug("blur command")
		m.command.SetValue("")
		m.command.Blur()
		return nil
	}

	m.command, cmd = m.command.Update(msg)
	return cmd
}

type ApplyCommandMsg struct {
	Items   []item
	Command string
}

func (m model) applyCommand(command string) tea.Cmd {
	return func() tea.Msg {
		selected := []item{}

		for _, i := range m.list.Items() {
			n, ok := i.(item)

			if !ok {
				continue
			}
			if n.selected {
				selected = append(selected, n)
			}
		}

		return ApplyCommandMsg{Items: selected, Command: command}
	}
}

type AppliedCommandMsg struct {
	Message string
}
type CleanListMsg struct {
	Message string
}

func (m model) applyNext() tea.Cmd {
	return func() tea.Msg {
		if len(m.processQueue) == 0 {
			slog.Debug("no more command to apply")
			return CleanListMsg{}
		}

		current, tail := m.processQueue[0], m.processQueue[1:]
		m.processQueue = tail

		slog.Debug("apply next", "notification", current.notification.String())

		message := ""
		out := &strings.Builder{}
		if err := m.currentActor.Run(current.notification, out); err != nil {
			message = fmt.Sprintf("Error for '%s': %s", current.notification.Subject.Title, err.Error())
		} else {
			message = out.String()
		}

		return AppliedCommandMsg{Message: message}
	}
}

func (m *model) handleResult(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch {

	case key.Matches(msg, m.list.KeyMap.ShowFullHelp):
		m.help.ShowAll = !m.help.ShowAll
		// TODO: why is this not showing the full help?
		slog.Debug("toggle help", "showAll", m.help.ShowAll)

	case key.Matches(msg, m.list.KeyMap.Quit):
		m.showResult = false
		return nil
	}

	m.result, cmd = m.result.Update(msg)
	return cmd
}

func (m *model) handleBrowsing(msg tea.KeyMsg) tea.Cmd {
	slog.Debug("browsing", "key", msg.String())

	switch {
	case key.Matches(msg, m.list.KeyMap.ShowFullHelp):
		m.help.ShowAll = !m.help.ShowAll
		slog.Debug("toggle help", "showAll", m.help.ShowAll)

	case key.Matches(msg, m.keymap.Toggle):
		if i, ok := m.list.SelectedItem().(item); ok {
			i.selected = !i.selected
			slog.Debug("toggle selected", "item", i.notification.Subject.Title, "selected", i.selected)

			return m.list.SetItem(m.list.GlobalIndex(), i)
		}

	case key.Matches(msg, m.keymap.All):
		items := []list.Item{}
		for _, e := range m.list.Items() {
			if i, ok := e.(item); ok {
				i.selected = true
				items = append(items, i)
			}
		}
		m.list.SetItems(items)

	case key.Matches(msg, m.keymap.None):
		items := []list.Item{}
		for _, e := range m.list.Items() {
			if i, ok := e.(item); ok {
				i.selected = false
				items = append(items, i)
			}
		}
		m.list.SetItems(items)

	case key.Matches(msg, m.keymap.Open):
		current, ok := m.list.SelectedItem().(item)
		if ok {
			m.actors["open"].Run(current.notification, os.Stderr)
		}

	case key.Matches(msg, m.keymap.CommandMode):
		slog.Debug("focus command")
		m.command.Focus()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *model) handleFiltering(msg tea.KeyMsg) tea.Cmd {
	slog.Debug("filtering", "key", msg.String())
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd
}
