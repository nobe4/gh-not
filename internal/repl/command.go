package repl

import (
	"fmt"
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/actions"
)

type Run struct {
	Runner actions.Runner
	Args   []string
}

func (m model) commandAndArgs(value, suggestion string) (string, []string) {
	if value == suggestion {
		return value, nil
	}

	values := strings.Split(value, " ")
	command, args := values[0], values[1:]

	// Manually get the new suggestion
	m.command.SetValue(command)
	// Workaround until https://github.com/charmbracelet/bubbles/pull/630 is
	// merged.
	m.command.SetSuggestions(m.command.AvailableSuggestions())
	command = m.command.CurrentSuggestion()

	return command, args
}

func (m model) acceptCommand() (tea.Model, tea.Cmd) {
	command, args := m.commandAndArgs(
		m.command.Value(),
		m.command.CurrentSuggestion(),
	)

	slog.Debug("acceptCommand", "command", command, "args", args)

	m.command.SetValue("")
	m.command.Blur()
	m.showResult = true

	return m, m.applyCommand(command, args)
}

func (m model) cancelCommand() (tea.Model, tea.Cmd) {
	slog.Debug("cancelCommand")

	m.command.SetValue("")
	m.command.Blur()

	return m, nil
}

type ApplyCommandMsg struct {
	Items   []item
	Command string
	Args    []string
}

func (m model) applyCommand(command string, args []string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("applyCommand", "command", command)

		selected := []item{}

		for _, i := range m.list.Items() {
			if n, ok := i.(item); ok && n.selected {
				selected = append(selected, n)
			}
		}

		return ApplyCommandMsg{Items: selected, Command: command, Args: args}
	}
}

func (msg ApplyCommandMsg) apply(m model) (tea.Model, tea.Cmd) {
	slog.Debug("apply command", "command", msg.Command)

	runner, ok := m.actions[msg.Command]
	if !ok {
		return m, m.renderResult(fmt.Errorf("Invalid command %s", msg.Command))
	}

	m.resultStrings = []string{}
	m.currentRun = Run{
		Runner: runner,
		Args:   msg.Args,
	}
	m.processQueue = msg.Items

	return m, tea.Sequence(m.renderResult(nil), m.applyNext())
}

type AppliedCommandMsg struct {
	Message string
}

func (msg AppliedCommandMsg) apply(m model) (tea.Model, tea.Cmd) {
	slog.Debug("applied command", "message", msg.Message)

	m.processQueue = m.processQueue[1:]

	m.resultStrings = append(m.resultStrings, msg.Message)

	return m, tea.Sequence(m.renderResult(nil), m.applyNext())
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
		if err := m.currentRun.Runner.Run(current.notification, m.currentRun.Args, out); err != nil {
			message = fmt.Sprintf("Error for '%s': %s", current.notification.Subject.Title, err.Error())
		} else {
			message = out.String()
		}

		return AppliedCommandMsg{Message: message}
	}
}
