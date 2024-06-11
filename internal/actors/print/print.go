package print

import (
	"io"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct{}

func (_ *Actor) Run(n *notifications.Notification, w io.Writer) error {
	if !n.Meta.Hidden {
		w.Write([]byte(n.ToString()))
	}

	return nil
}
