package gh

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/nobe4/gh-not/internal/notifications"
)

type ThreadExtra struct {
	User    notifications.User `json:"user"`
	State   string             `json:"state"`
	HtmlUrl string             `json:"html_url"`
}

func (c *Client) Enrich(n *notifications.Notification) error {
	if n == nil {
		return nil
	}

	threadExtra, err := c.getThreadExtra(n)
	if err != nil {
		return err
	}
	n.Author = threadExtra.User
	n.Subject.State = threadExtra.State
	n.Subject.HtmlUrl = threadExtra.HtmlUrl

	LastCommentor, err := c.getLastCommentor(n)
	if err != nil {
		return err
	}
	n.LatestCommentor = LastCommentor

	return nil
}

func (c *Client) getThreadExtra(n *notifications.Notification) (ThreadExtra, error) {
	if n.Subject.URL == "" {
		return ThreadExtra{}, nil
	}

	slog.Debug("getting the thread extra", "id", n.Id, "url", n.Subject.URL)
	resp, err := c.API.Request(http.MethodGet, n.Subject.URL, nil)
	if err != nil {
		return ThreadExtra{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ThreadExtra{}, err
	}

	extra := ThreadExtra{}
	err = json.Unmarshal(body, &extra)
	if err != nil {
		return ThreadExtra{}, err
	}

	return extra, nil
}

func (c *Client) getLastCommentor(n *notifications.Notification) (notifications.User, error) {
	if n.Subject.LatestCommentUrl == "" {
		return notifications.User{}, nil
	}

	slog.Debug("getting the last commentor", "id", n.Id, "url", n.Subject.LatestCommentUrl)
	resp, err := c.API.Request(http.MethodGet, n.Subject.LatestCommentUrl, nil)
	if err != nil {
		return notifications.User{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return notifications.User{}, err
	}

	comment := struct {
		User notifications.User `json:"user"`
	}{}
	err = json.Unmarshal(body, &comment)
	if err != nil {
		return notifications.User{}, err
	}

	return comment.User, nil
}
