package actors

import (
	"fmt"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Actor interface {
	Run(notifications.Notification) (notifications.Notification, error)
}

type DebugActor struct{}

func (_ *DebugActor) Run(n notifications.Notification) (notifications.Notification, error) {
	fmt.Printf("DEBUG Run %s\n", n.ToString())
	return n, nil
}

type PrintActor struct{}

func (_ *PrintActor) Run(n notifications.Notification) (notifications.Notification, error) {
	if n.Meta.Hidden {
		return n, nil
	}

	fmt.Println(n.ToString())
	return n, nil
}

type HideActor struct{}

func (_ *HideActor) Run(n notifications.Notification) (notifications.Notification, error) {
	// fmt.Printf("HIDDING %s\n", n.ToString())
	n.Meta.Hidden = !n.Meta.Hidden
	return n, nil
}
