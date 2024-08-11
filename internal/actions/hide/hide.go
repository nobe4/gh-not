/*
Package hide implements an [actions.Runner] that hides a notification.
It hides the notifications completely.
*/
package hide

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (_ *Runner) Run(n *notifications.Notification, w io.Writer) error {
	slog.Debug("marking notification as hidden", "notification", n.Id)

	n.Meta.Hidden = true

	fmt.Fprint(w, colors.Red("HIDE ")+n.String())

	return nil
}
