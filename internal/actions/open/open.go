/*
Package open implements an [actions.Runner] that opens a notification in the browser.
*/
package open

import (
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/cli/go-gh/pkg/browser"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

var errNoURL = errors.New("no URL to open")

type Runner struct {
	Client *gh.Client
}

func (*Runner) Run(n *notifications.Notification, _ []string, w io.Writer) error {
	slog.Debug("open notification in browser", "notification", n)

	b := browser.New("", w, w)

	if n.Subject.HTMLURL == "" {
		return errNoURL
	}

	fmt.Fprint(w, colors.Blue("OPEN ")+n.Subject.URL+" ")

	if err := b.Browse(n.Subject.HTMLURL); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	return nil
}
