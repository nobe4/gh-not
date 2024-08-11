package actions

import (
	"io"

	"github.com/nobe4/gh-not/internal/actions/debug"
	"github.com/nobe4/gh-not/internal/actions/done"
	"github.com/nobe4/gh-not/internal/actions/hide"
	"github.com/nobe4/gh-not/internal/actions/open"
	"github.com/nobe4/gh-not/internal/actions/pass"
	"github.com/nobe4/gh-not/internal/actions/print"
	"github.com/nobe4/gh-not/internal/actions/read"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type ActionsMap map[string]Actor

func Map(client *gh.Client) ActionsMap {
	return map[string]Actor{
		"pass":  &pass.Actor{},
		"debug": &debug.Actor{},
		"print": &print.Actor{},
		"hide":  &hide.Actor{},
		"read":  &read.Actor{Client: client},
		"done":  &done.Actor{Client: client},
		"open":  &open.Actor{Client: client},
	}
}

type Actor interface {
	Run(*notifications.Notification, io.Writer) error
}
