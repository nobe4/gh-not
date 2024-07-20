package gh

import (
	"log/slog"
	"net/http"

	"github.com/nobe4/gh-not/internal/notifications"
)

type Extra struct {
	User    notifications.User `json:"user"`
	State   string             `json:"state"`
	HtmlUrl string             `json:"html_url"`
}

func (c *Client) Enrich(n *notifications.Notification) error {
	if n == nil {
		return nil
	}

	extra := Extra{}
	if err := c.API.Do(http.MethodGet, n.Subject.URL, nil, &extra); err != nil {
		return err
	}

	slog.Debug("enriching", "notifications", n.Debug())

	n.Author = extra.User
	n.Subject.State = extra.State
	n.Subject.HtmlUrl = extra.HtmlUrl

	return nil
}
