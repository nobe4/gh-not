package gh

import (
	"log/slog"
	"net/http"

	"github.com/nobe4/gh-not/internal/notifications"
)

func (c *Client) enrichNotification(n *notifications.Notification) error {
	extra := struct {
		User  notifications.User `json:"user"`
		State string             `json:"state"`
	}{}
	if err := c.API.Do(http.MethodGet, n.Subject.URL, nil, &extra); err != nil {
		return err
	}

	slog.Debug("adding author to notification", "notifications", n)
	n.Author = extra.User

	slog.Debug("adding state to notification's suject", "notifications", n)
	n.Subject.State = extra.State

	return nil
}
