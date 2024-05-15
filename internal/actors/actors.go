package actors

import (
	"fmt"

	"github.com/nobe4/ghnot/internal/notifications"
)

type Actor interface {
	Run(notifications.Notification) error
}

type DebugActor struct{}

func (_ *DebugActor) Run(n notifications.Notification) error {
	fmt.Printf("DEBUG Run %#v\n", n)
	return nil
}

type PrintActor struct{}

func (_ *PrintActor) Run(n notifications.Notification) error {
	fmt.Printf("%#v\n", n)
	return nil
}
