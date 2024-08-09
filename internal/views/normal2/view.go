package normal2

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	noStyle = lipgloss.NewStyle()
)

func (m *model) initView() {
	m.list.SetShowStatusBar(false)
	m.list.SetShowTitle(false)
	m.list.SetShowFilter(false)
	m.list.SetShowHelp(false)
	m.list.Help.ShowAll = false
	m.list.SetDelegate(itemDelegate{})
	m.list.SetShowPagination(false)

	m.list.Paginator.Type = paginator.Arabic

	m.list.FilterInput.Prompt = "/"
	m.list.FilterInput.Cursor.Style = noStyle
	m.list.FilterInput.PromptStyle = noStyle
	m.list.FilterInput.Placeholder = "filter"

	m.list.SetShowFilter(false)

	m.list.Styles = list.Styles{
		PaginationStyle: noStyle,
		HelpStyle:       noStyle,
		StatusBar:       noStyle,
	}

	m.command.Prompt = ":"
	m.command.Cursor.Style = noStyle
	m.command.PromptStyle = noStyle
	m.command.Placeholder = "command"

	suggestions := []string{}
	for k := range m.actors {
		suggestions = append(suggestions, k)
	}

	m.command.SetSuggestions(suggestions)
	m.command.ShowSuggestions = true
}

func (m *model) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	slog.Debug("resize", "width", msg.Width, "height", msg.Height)

	m.list.SetHeight(min(msg.Height, m.maxHeigth))
	m.list.SetWidth(msg.Width)

	m.result.Height = min(msg.Height, m.maxHeigth)
	m.result.Width = msg.Width

	if !m.ready {
		m.ready = true
	}

	return m, nil
}

type ResultUpdateMsg struct {
	content string
}

func (m model) renderResult(err error) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("renderResult")

		lines := []string{}

		if err != nil {
			lines = append(lines, err.Error())
		} else {
			lines = append(lines, m.resultStrings...)
			if len(m.processQueue) > 0 {
				lines = append(lines, fmt.Sprintf("%d more ...", len(m.processQueue)))
			}
		}

		lines = append(lines, fmt.Sprintf("press %s to continue ...", m.list.KeyMap.Quit.Keys()))

		return ResultUpdateMsg{lipgloss.JoinVertical(lipgloss.Top, lines...)}
	}
}

func (m model) viewFullHelp() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		"default", m.list.Help.FullHelpView(m.keymap.FullHelp()),
		"\nlist and filter", m.list.Help.FullHelpView(m.list.FullHelp()),
		"\nresults", m.list.Help.FullHelpView(ViewportKeymap{m.result.KeyMap}.FullHelp()),
	)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.showHelp {
		return m.viewFullHelp()
	}

	content := ""
	statusLine := ""

	if m.showResult {
		content = m.result.View()
		slog.Debug("showResult")

		if !m.result.AtTop() {
			statusLine = m.list.Help.Styles.ShortDesc.Render("↑ to scroll up ")
		}
		if !m.result.AtBottom() {
			statusLine += m.list.Help.Styles.ShortDesc.Render("↓ to scroll down ")
		}

		statusLine += m.list.Help.Styles.ShortDesc.Render("? to toggle help")
	} else {
		statusLine = m.list.Paginator.View() + " "

		if m.command.Focused() {
			statusLine = m.command.View()
		} else {
			if m.list.FilterState() == list.Filtering {
				statusLine += m.list.FilterInput.View()
			} else {
				statusLine += m.list.Help.Styles.ShortDesc.Render("? to toggle help")
			}
		}

		listView := m.list.View()

		content = noStyle.Height(m.list.Height() - 1).Render(listView)
	}

	sections := []string{
		content,
		statusLine,
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
