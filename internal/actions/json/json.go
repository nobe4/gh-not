/*
Package json implements an [actions.Runner] that prints a notification in JSON.
*/
package json

import (
	"fmt"
	"io"

	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (*Runner) Run(n *notifications.Notification, filters []string, w io.Writer) error {
	if n.Meta.Hidden {
		return nil
	}

	filter := ""
	if len(filters) > 0 {
		filter = filters[0]
	}

	result, err := jq.Run(filter, *n)
	if err != nil {
		return fmt.Errorf("failed to run json with filter '%s': %w", filter, err)
	}

	fmt.Fprint(w, result)

	return nil
}
