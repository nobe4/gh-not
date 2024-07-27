package open

import (
	"io"
	"log/slog"

	"github.com/cli/go-gh/pkg/browser"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct {
	Client *gh.Client
}

func (a *Actor) Run(n *notifications.Notification, w io.Writer) error {
	slog.Debug("open notification in browser", "notification", n)

	browser := browser.New("", w, w)
	return browser.Browse(n.Subject.HtmlUrl)
}
