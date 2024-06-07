package print

import (
	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct{}

func (_ *Actor) Run(n *notifications.Notification) (string, error) {
	if n.Meta.Hidden {
		return "", nil
	}

	return n.ToString(), nil
}
