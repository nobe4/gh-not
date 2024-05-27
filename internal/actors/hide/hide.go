package hide

import "github.com/nobe4/gh-not/internal/notifications"

type Actor struct{}

func (_ *Actor) Run(n notifications.Notification) (notifications.Notification, string, error) {
	n.Meta.Hidden = !n.Meta.Hidden

	return n, "", nil
}
