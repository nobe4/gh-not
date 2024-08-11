package hide

import (
	"io"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct{}

func (_ *Actor) Run(n *notifications.Notification, _ io.Writer) error {
	n.Meta.Hidden = !n.Meta.Hidden

	n = &notifications.Notification{
		Id: n.Id,
		Meta: notifications.Meta{
			Hidden: n.Meta.Hidden,
		},
	}

	return nil
}
