/*
Package print implements an [actions.Runner] that prints a notification.
*/
package print

import (
	"fmt"
	"io"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (_ *Runner) Run(n *notifications.Notification, w io.Writer) error {
	if !n.Meta.Hidden {
		fmt.Fprint(w, n.String())
	}

	return nil
}
