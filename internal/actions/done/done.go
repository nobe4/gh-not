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

func (a *Runner) Run(n *notifications.Notification, w io.Writer) error {
	slog.Debug("marking notification as done", "notification", n)

	n.Meta.Done = true

	err := a.Client.API.Do(http.MethodDelete, n.URL, nil, nil)
	if err != nil {
		return err
	}

	fmt.Fprint(w, colors.Red("DONE ")+n.String())

	return nil
}
