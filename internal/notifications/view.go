package notifications

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/cli/go-gh/v2/pkg/text"

	"github.com/nobe4/gh-not/internal/colors"
)

//nolint:gochecknoglobals // This map is used a lot.
var prettyRead = map[bool]string{
	false: colors.Red("RD"),
	true:  colors.Green("UR"),
}

//nolint:gochecknoglobals // This map is used a lot.
var prettyTypes = map[string]string{
	"Issue":       colors.Blue("IS"),
	"PullRequest": colors.Cyan("PR"),
	"Discussion":  colors.Green("DS"),
}

//nolint:gochecknoglobals // This map is used a lot.
var prettyState = map[string]string{
	"open":   colors.Green("OP"),
	"closed": colors.Red("CL"),
	"merged": colors.Magenta("MG"),
}

func (n Notification) String() string {
	if n.rendered != "" {
		return n.rendered
	}

	return fmt.Sprintf("[%s] %s at %s: %s", n.ID, n.Author.Login, n.UpdatedAt, n.Subject.Title)
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

func (n Notifications) TagsMap() map[string]int {
	tags := map[string]int{}

	for _, n := range n {
		for _, t := range n.Meta.Tags {
			tags[t]++
		}
	}

	return tags
}

func (n Notification) Visible() bool {
	return !n.Meta.Done && !n.Meta.Hidden
}

// Render the notifications in a human readable format.
// If possible, render a table, otherwise render a simple string.
//
//revive:disable:cognitive-complexity // TODO: simplify.
func (n Notifications) Render() error {
	if len(n) == 0 {
		return nil
	}

	// Default to a simple string
	for _, n := range n {
		n.rendered = fmt.Sprintf(
			"%s %s %s %s by %s at %s: '%s'",
			n.prettyRead(),
			n.prettyType(),
			n.prettyState(),
			n.Repository.FullName,
			n.Author.Login,
			text.RelativeTimeAgo(time.Now(), n.UpdatedAt),
			n.Subject.Title)
	}

	// Try to render a table
	out := bytes.Buffer{}

	t := term.FromEnv()

	w, _, err := t.Size()
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
	}

	printer := tableprinter.New(&out, t.IsTerminalOutput(), w)

	for _, n := range n {
		printer.AddField(n.prettyRead())
		printer.AddField(n.prettyType())
		printer.AddField(n.prettyState())
		printer.AddField(n.Repository.FullName)
		printer.AddField(n.Author.Login)
		printer.AddField(n.Subject.Title)

		relativeTime := text.RelativeTimeAgo(time.Now(), n.UpdatedAt)
		if n.LatestCommentor.Login != "" {
			relativeTime += " by " + n.LatestCommentor.Login
		}

		printer.AddField(relativeTime)
		printer.EndRow()
	}

	if err := printer.Render(); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	tableLines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")

	slog.Info("Rendered notifiation in a table", "notification count", len(n), "table line count", len(tableLines))

	for i, l := range tableLines {
		n[i].rendered = l
	}

	return nil
}
