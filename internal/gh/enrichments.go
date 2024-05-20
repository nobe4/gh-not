package gh

import (
	"log/slog"
	"net/http"

	"github.com/nobe4/gh-not/internal/notifications"
)

func (c *Client) enrichNotification(n notifications.Notification) (notifications.Notification, error) {
	var err error

	if n, err = c.addAuthor(n); err != nil {
		return n, err
	}

	return n, nil
}

func (c *Client) addAuthor(n notifications.Notification) (notifications.Notification, error) {
	slog.Debug("adding author to notification", "notifications", n)
	author := struct {
		User notifications.User `json:"user"`
	}{}
	if err := c.API.Do(http.MethodGet, n.Subject.URL, nil, &author); err != nil {
		return n, err
	}

	n.Author = author.User

	return n, nil
}
