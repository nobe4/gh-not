package repl

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/notifications"
)

type item struct {
	// It's not possible to rely on bubbles/list's model to get the global index
	// of an item, so we have to manage it manually.
	index        int
	notification *notifications.Notification
	selected     bool
}

func (i item) FilterValue() string { return i.notification.String() }

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
