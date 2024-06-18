package gh

import (
	"log/slog"
	"net/http"

	"github.com/nobe4/gh-not/internal/notifications"
)

func (c *Client) enrichNotification(n *notifications.Notification) error {
	if n.Meta.Done {
		return nil
	}

	extra := struct {
		User    notifications.User `json:"user"`
		State   string             `json:"state"`
		HtmlUrl string             `json:"html_url"`
	}{}
	if err := c.API.Do(http.MethodGet, n.Subject.URL, nil, &extra); err != nil {
		return err
	}

	slog.Debug("enriching", "notifications", n.Debug())

	n.Author = extra.User
	n.Subject.State = extra.State
	n.Subject.HtmlUrl = extra.HtmlUrl

	return nil
}
