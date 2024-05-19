package done

import (
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

// Actor that marks a notification as done.
// Ref: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#mark-a-thread-as-done
type Actor struct {
	Client *gh.Client
}

func (a *Actor) Run(n notifications.Notification) (notifications.Notification, error) {
	emptyNotification := notifications.Notification{}

	err := a.Client.API.Delete(n.URL, nil)
	if err != nil {
		return emptyNotification, err
	}

	return emptyNotification, nil
}
