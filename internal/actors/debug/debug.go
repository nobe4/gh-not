package debug

import (
	"io"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct{}

func (_ *Actor) Run(n *notifications.Notification, w io.Writer) error {
	w.Write([]byte("DEBUG" + n.ToString()))

	return nil
}
