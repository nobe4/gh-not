package actors

import (
	"fmt"

	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type ActorsMap map[string]Actor

func Map(client *gh.Client) ActorsMap {
	return map[string]Actor{
		"pass":  &PassActor{},
		"debug": &DebugActor{},
		"print": &PrintActor{},
		"hide":  &HideActor{},
		"read":  &ReadActor{client: client},
		"done":  &DoneActor{client: client},
	}
}

type Actor interface {
	Run(notifications.Notification) (notifications.Notification, error)
}

type PassActor struct{}

func (_ *PassActor) Run(n notifications.Notification) (notifications.Notification, error) {
	return n, nil
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
	n.Meta.Hidden = !n.Meta.Hidden
	return n, nil
}

// Actor that marks a notification as read.
// Ref: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#mark-a-thread-as-read
type ReadActor struct {
	client *gh.Client
}

func (a *ReadActor) Run(n notifications.Notification) (notifications.Notification, error) {
	err := a.client.API.Patch(n.URL, nil, nil)

	// go-gh currently fails to handle HTTP-205 correctly, however it's possible
	// to catch this case.
	// ref: https://github.com/cli/go-gh/issues/161
	if err != nil && err.Error() != "unexpected end of JSON input" {
		return n, err
	}

	n.Unread = false

	return n, nil
}

// Actor that marks a notification as done.
// Ref: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#mark-a-thread-as-done
type DoneActor struct {
	client *gh.Client
}

func (a *DoneActor) Run(n notifications.Notification) (notifications.Notification, error) {
	emptyNotification := notifications.Notification{}

	err := a.client.API.Delete(n.URL, nil)
	if err != nil {
		return emptyNotification, err
	}

	return emptyNotification, nil
}
