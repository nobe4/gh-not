package print

import (
	"fmt"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct{}

func (_ *Actor) Run(n notifications.Notification) (notifications.Notification, error) {
	if n.Meta.Hidden {
		return n, nil
	}

	fmt.Println(n.ToString())

	return n, nil
}
