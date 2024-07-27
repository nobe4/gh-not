package notifications

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/cli/go-gh/v2/pkg/text"
	"github.com/nobe4/gh-not/internal/colors"
)

var prettyRead = map[bool]string{
	false: colors.Red("RD"),
	true:  colors.Green("UR"),
}

var prettyTypes = map[string]string{
	"Issue":       colors.Blue("IS"),
	"PullRequest": colors.Cyan("PR"),
}

var prettyState = map[string]string{
	"open":   colors.Green("OP"),
	"closed": colors.Red("CL"),
	"merged": colors.Magenta("MG"),
}

func (n Notification) String() string {
	return n.rendered
}

func (n Notification) prettyRead() string {
	if p, ok := prettyRead[n.Unread]; ok {
		return p
	}

	return colors.Yellow("R?")
}

func (n Notification) prettyType() string {
	if p, ok := prettyTypes[n.Subject.Type]; ok {
		return p
	}

	return colors.Yellow("T?")
}

func (n Notification) prettyState() string {
	if p, ok := prettyState[n.Subject.State]; ok {
		return p
	}

	return colors.Yellow("S?")
}

func (n Notifications) String() string {
	out := ""
	for _, n := range n {
		out += fmt.Sprintf("%s\n", n)
	}
	return out
}

func (n Notifications) Visible() Notifications {
	visible := Notifications{}
	for _, n := range n {
		if n.Visible() {
			visible = append(visible, n)
		}
	}
	return visible
}

func (n Notification) Visible() bool {
	return !n.Meta.Done && !n.Meta.Hidden
}

// Render the notifications in a human readable format.
// If possible, render a table, otherwise render a simple string.
func (n Notifications) Render() error {
	// Default to a simple string
	for _, n := range n {
		n.rendered = fmt.Sprintf("%s %s %s %s by %s at %s: '%s'", n.prettyRead(), n.prettyType(), n.prettyState(), n.Repository.FullName, n.Author.Login, text.RelativeTimeAgo(time.Now(), n.UpdatedAt), n.Subject.Title)
	}

	// Try to render a table
	out := bytes.Buffer{}

	t := term.FromEnv()
	w, _, err := t.Size()
	if err != nil {
		return err
	}

	printer := tableprinter.New(&out, t.IsTerminalOutput(), w)

	for _, n := range n {
		printer.AddField(n.prettyRead())
		printer.AddField(n.prettyType())
		printer.AddField(n.prettyState())
		printer.AddField(n.Repository.FullName)
		printer.AddField(n.Author.Login)
		printer.AddField(n.Subject.Title)
		printer.AddField(text.RelativeTimeAgo(time.Now(), n.UpdatedAt))
		printer.EndRow()
	}

	if err := printer.Render(); err != nil {
		return err
	}

	for i, l := range strings.Split(strings.TrimRight(out.String(), "\n"), "\n") {
		n[i].rendered = l
	}

	return nil
}
