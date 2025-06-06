/*
Package read implements an [actions.Runner] that marks a notification as read.

It updates Unread and marks the notification's thread as read on GitHub.
Ref: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#mark-a-thread-as-read
*/
package read

import (
	"fmt"
	"io"
	"net/http"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct {
	Client *gh.Client
}

func (a *Runner) Run(n *notifications.Notification, _ []string, w io.Writer) error {
	r, err := a.Client.API.Request(http.MethodPatch, n.URL, nil)

	// go-gh currently fails to handle HTTP-205 correctly, however it's possible
	// to catch this case.
	// ref: https://github.com/cli/go-gh/issues/161
	if err != nil && err.Error() != "unexpected end of JSON input" {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	defer r.Body.Close()

	n.Unread = false

	fmt.Fprint(w, colors.Yellow("READ ")+n.String())

	return nil
}
