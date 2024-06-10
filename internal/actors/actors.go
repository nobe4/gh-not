package actors

import (
	"github.com/nobe4/gh-not/internal/actors/debug"
	"github.com/nobe4/gh-not/internal/actors/done"
	"github.com/nobe4/gh-not/internal/actors/hide"
	"github.com/nobe4/gh-not/internal/actors/pass"
	"github.com/nobe4/gh-not/internal/actors/print"
	"github.com/nobe4/gh-not/internal/actors/read"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type ActorsMap map[string]Actor

func Map(client *gh.Client) ActorsMap {
	return map[string]Actor{
		"pass":  &pass.Actor{},
		"debug": &debug.Actor{},
		"print": &print.Actor{},
		"hide":  &hide.Actor{},
		"read":  &read.Actor{Client: client},
		"done":  &done.Actor{Client: client},
	}
}

type Actor interface {
	Run(*notifications.Notification) (string, error)
}
