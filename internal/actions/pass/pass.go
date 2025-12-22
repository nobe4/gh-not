/*
Package pass implements an [actions.Runner] that does nothing.
*/
package pass

import (
	"fmt"
	"io"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (*Runner) Run(n *notifications.Notification, _ []string, w io.Writer) error {
	fmt.Fprint(w, colors.Blue("PASS ")+n.String())
	return nil
}
