package normal2

import (
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

	m.help.Styles = m.list.Help.Styles
}

func (m *model) handleResize(msg tea.WindowSizeMsg) {
	slog.Debug("resize", "width", msg.Width, "height", msg.Height)

	m.list.SetHeight(min(msg.Height, m.maxHeigth))
	m.list.SetWidth(msg.Width)

	m.result.Height = min(msg.Height, m.maxHeigth)
	m.result.Width = msg.Width

	if !m.ready {
		m.ready = true
	}
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.help.ShowAll {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"default", m.help.View(m.keymap),
			"list and filter", m.list.Help.View(m.list),
			"\ncommand", m.help.View(ViewportKeymap{m.result.KeyMap}),
		)
	}

	content := ""
	statusLine := ""

	if m.showResult {
		content = m.result.View()

		if !m.result.AtTop() {
			statusLine = m.help.Styles.ShortDesc.Render("↑ to scroll up ")
		}
		if !m.result.AtBottom() {
			statusLine += m.help.Styles.ShortDesc.Render("↓ to scroll down ")
		}

		statusLine += m.help.Styles.ShortDesc.Render("? to toggle help")
	} else {
		statusLine = m.list.Paginator.View() + " "

		if m.command.Focused() {
			statusLine = m.command.View()
		} else {
			if m.list.FilterState() == list.Filtering {
				statusLine += m.list.FilterInput.View()
			} else {
				statusLine += m.help.Styles.ShortDesc.Render("? to toggle help")
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
