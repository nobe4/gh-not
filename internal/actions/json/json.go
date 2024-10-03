/*
Package json implements an [actions.Runner] that prints a notification in JSON.
*/
package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (_ *Runner) Run(n *notifications.Notification, _ []string, w io.Writer) error {
	if n.Meta.Hidden {
		return nil
	}

	marshalled, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(marshalled))

	return nil
}
