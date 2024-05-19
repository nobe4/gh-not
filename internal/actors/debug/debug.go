package debug

import (
	"fmt"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor struct{}

func (_ *Actor) Run(n notifications.Notification) (notifications.Notification, error) {
	fmt.Printf("DEBUG Run %s\n", n.ToString())
	return n, nil
}
