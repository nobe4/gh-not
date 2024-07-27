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

func (n Notification) String() string {
	return fmt.Sprintf("%s %s %s %s by %s at %s: '%s'", n.prettyRead(), n.prettyType(), n.prettyState(), n.Repository.FullName, n.Author.Login, text.RelativeTimeAgo(time.Now(), n.UpdatedAt), n.Subject.Title)
}

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
		out += n.String() + "\n"
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

func (n Notifications) Table() (string, error) {
	out := bytes.Buffer{}

	t := term.FromEnv()
	w, _, err := t.Size()
	if err != nil {
		return "", err
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
		return "", err
	}

	return strings.TrimRight(out.String(), "\n"), nil
}
