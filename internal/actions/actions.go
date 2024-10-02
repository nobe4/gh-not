package actions

import (
	"io"

	"github.com/nobe4/gh-not/internal/actions/assign"
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

type ActionsMap map[string]Runner

func Map(client *gh.Client) ActionsMap {
	return map[string]Runner{
		"pass":   &pass.Runner{},
		"debug":  &debug.Runner{},
		"print":  &print.Runner{},
		"hide":   &hide.Runner{},
		"read":   &read.Runner{Client: client},
		"done":   &done.Runner{Client: client},
		"open":   &open.Runner{Client: client},
		"assign": &assign.Runner{Client: client},
	}
}

type Runner interface {
	Run(*notifications.Notification, []string, io.Writer) error
}
