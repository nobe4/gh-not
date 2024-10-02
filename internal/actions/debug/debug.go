/*
Package debug implements an [actions.Runner] that prints a notification with a DEBUG prefix.
*/
package debug

import (
	"fmt"
	"io"
	"strings"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (_ *Runner) Run(n *notifications.Notification, args []string, w io.Writer) error {
	fmt.Fprint(w, colors.Yellow("DEBUG ")+n.String()+" "+strings.Join(args, ", "))

	return nil
}
