package pass

import (
	"io"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct{}

func (_ *Actor) Run(n *notifications.Notification, _ io.Writer) error {
	return nil
}
