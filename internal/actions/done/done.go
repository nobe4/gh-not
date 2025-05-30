/*
Package done implements an [actions.Runner] that marks a notification as done.

It updates Meta.Done and marks the notification's thread as done on GitHub.
The notification will be hidden until the thread is updated.
Ref: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#mark-a-thread-as-done
*/
package done

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct {
	Client *gh.Client
}

func (a *Runner) Run(n *notifications.Notification, _ []string, w io.Writer) error {
	slog.Debug("marking notification as done", "notification", n)

	n.Meta.Done = true

	r, err := a.Client.API.Request(http.MethodDelete, n.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to mark notification as done: %w", err)
	}
	defer r.Body.Close()

	fmt.Fprint(w, colors.Red("DONE ")+n.String())

	return nil
}
