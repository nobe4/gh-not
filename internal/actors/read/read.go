package read

import (
	"net/http"

	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

// Actor that marks a notification as read.
// Ref: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#mark-a-thread-as-read
type Actor struct {
	Client *gh.Client
}

func (a *Actor) Run(n notifications.Notification) (notifications.Notification, error) {
	err := a.Client.API.Do(http.MethodPatch, n.URL, nil, nil)

	// go-gh currently fails to handle HTTP-205 correctly, however it's possible
	// to catch this case.
	// ref: https://github.com/cli/go-gh/issues/161
	if err != nil && err.Error() != "unexpected end of JSON input" {
		return n, err
	}

	n.Unread = false

	return n, nil
}
