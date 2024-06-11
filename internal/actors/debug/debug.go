package debug

import (
	"fmt"
	"io"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct{}

func (_ *Actor) Run(n *notifications.Notification, w io.Writer) error {
	fmt.Fprint(w, colors.Yellow("DEBUG"+n.ToString()))

	return nil
}
