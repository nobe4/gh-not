/*
Package pass implements an [actions.Runner] that does nothing.
*/
package pass

import (
	"io"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (*Runner) Run(_ *notifications.Notification, _ []string, _ io.Writer) error {
	return nil
}
