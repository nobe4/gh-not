package pass

import (
	"io"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (_ *Runner) Run(n *notifications.Notification, _ io.Writer) error {
	return nil
}
