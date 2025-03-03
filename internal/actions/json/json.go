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

func (*Runner) Run(n *notifications.Notification, _ []string, w io.Writer) error {
	if n.Meta.Hidden {
		return nil
	}

	marshaled, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	fmt.Fprint(w, string(marshaled))

	return nil
}
