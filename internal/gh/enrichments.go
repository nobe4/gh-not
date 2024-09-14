package gh

import (
	"encoding/json"
	"io"
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

	slog.Debug("enriching", "id", n.Id, "url", n.Subject.URL)
	resp, err := c.API.Request(http.MethodGet, n.Subject.URL, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	extra := Extra{}
	err = json.Unmarshal(body, &extra)
	if err != nil {
		return err
	}

	n.Author = extra.User
	n.Subject.State = extra.State
	n.Subject.HtmlUrl = extra.HtmlUrl

	return nil
}
