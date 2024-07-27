package normal2

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nobe4/gh-not/internal/notifications"
)

const listHeight = 10
const listDefaultWidth = 20

var (
	paginationStyle   = lipgloss.NewStyle().PaddingLeft(0)
	helpStyle         = lipgloss.NewStyle().PaddingLeft(0).PaddingBottom(0)
	filterPromptStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#ff0000", Dark: "#00ff00"})
	filterCursorStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#00ffff", Dark: "#ffff00"})
	quitTextStyle     = lipgloss.NewStyle().MarginBottom(1)
)

type item struct {
	*notifications.Notification
}

func (i item) FilterValue() string { return i.String() }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	n, ok := listItem.(item)
	if !ok {
		return
	}

	cursor := " "
	str := n.String()

	if index == m.Index() {
		cursor = ">"
		str = strings.ReplaceAll(str, " ", "⋅")
	}

	fmt.Fprint(w, cursor+str)
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.Subject.Title)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}

	return m.list.View() + "\n" + m.list.Paginator.View() + " " + m.list.FilterInput.View()
}

func Init(n notifications.Notifications) {
	items := []list.Item{}

	for _, notification := range n {
		items = append(items, item{notification})
	}

	l := list.New(items, itemDelegate{}, listDefaultWidth, listHeight)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	l.Paginator.Type = paginator.Arabic

	l.FilterInput.Prompt = "/"
	l.FilterInput.Cursor.Style = lipgloss.NewStyle()
	l.FilterInput.PromptStyle = lipgloss.NewStyle()

	l.Styles = list.Styles{
		PaginationStyle: paginationStyle,
		HelpStyle:       helpStyle,
		StatusBar:       lipgloss.NewStyle(),
	}

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
