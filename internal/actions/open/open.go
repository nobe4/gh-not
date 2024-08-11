/*
Package open implements an [actions.Runner] that opens a notification in the browser.
*/
package open

import (
	"io"
	"log/slog"

	"github.com/cli/go-gh/pkg/browser"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct {
	Client *gh.Client
}

func (a *Runner) Run(n *notifications.Notification, w io.Writer) error {
	slog.Debug("open notification in browser", "notification", n)

	browser := browser.New("", w, w)
	return browser.Browse(n.Subject.HtmlUrl)
}
