package normal2

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	noStyle       = lipgloss.NewStyle()
	quitTextStyle = noStyle.MarginBottom(1)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	n, ok := listItem.(item)
	if !ok {
		return
	}

	selected := " "
	if n.selected {
		selected = "x"
	}
	cursor := " "

	str := n.notification.String()
	if index == m.Index() {
		cursor = ">"
		str = strings.ReplaceAll(str, " ", "⋅")
	}

	fmt.Fprintf(w, "%s%s%s", selected, cursor, str)
}

func (m *model) initView() {
	m.list.SetShowStatusBar(false)
	m.list.SetShowTitle(false)
	m.list.SetShowFilter(false)
	m.list.SetShowHelp(false)
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
	m.command.Placeholder = "filter"

	suggestions := []string{}
	for k := range m.actors {
		suggestions = append(suggestions, k)
	}

	m.command.SetSuggestions(suggestions)
	m.command.ShowSuggestions = true

	m.help.Styles = m.list.Help.Styles
}

func (m model) View() string {
	if m.help.ShowAll {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"default",
			m.help.View(m.keymap),
			"\nlist and filter",
			m.list.Help.View(m.list),
		)
	}

	paginationLine := m.list.Paginator.View() + " "
	if m.command.Focused() {
		paginationLine = m.command.View()
	} else {
		if m.list.FilterState() == list.Filtering {
			paginationLine += m.list.FilterInput.View()
		} else {
			paginationLine += m.help.Styles.ShortDesc.Render("? to toggle help")
		}
	}

	listView := m.list.View()

	content := noStyle.Height(m.list.Height() - 1).Render(listView)
	sections := []string{
		content,
		paginationLine,
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
