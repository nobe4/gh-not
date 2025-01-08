package gh

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/nobe4/gh-not/internal/notifications"
)

type ThreadExtra struct {
	User           notifications.User   `json:"user"`
	State          string               `json:"state"`
	HTMLURL        string               `json:"html_url"`
	Assignees      []notifications.User `json:"assignees"`
	Reviewers      []notifications.User `json:"requested_reviewers"`
	ReviewersTeams []notifications.Team `json:"requested_teams"`
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
	n.Subject.HTMLURL = threadExtra.HTMLURL
	n.Assignees = threadExtra.Assignees
	n.Reviewers = threadExtra.Reviewers
	n.ReviewersTeams = threadExtra.ReviewersTeams

	lastCommentor, err := c.getLastCommentor(n)
	if err != nil {
		return err
	}

	n.LatestCommentor = lastCommentor

	return nil
}

func (c *Client) getThreadExtra(n *notifications.Notification) (ThreadExtra, error) {
	if n.Subject.URL == "" {
		return ThreadExtra{}, nil
	}

	slog.Debug("getting the thread extra", "id", n.ID, "url", n.Subject.URL)

	resp, err := c.API.Request(http.MethodGet, n.Subject.URL, nil)
	if err != nil {
		return ThreadExtra{}, fmt.Errorf("failed to get thread extra: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ThreadExtra{}, fmt.Errorf("failed to read thread extra: %w", err)
	}

	extra := ThreadExtra{}

	err = json.Unmarshal(body, &extra)
	if err != nil {
		return ThreadExtra{}, fmt.Errorf("failed to unmarshal thread extra: %w", err)
	}

	return extra, nil
}

func (c *Client) getLastCommentor(n *notifications.Notification) (notifications.User, error) {
	if n.Subject.LatestCommentURL == "" {
		return notifications.User{}, nil
	}

	slog.Debug("getting the last commentor", "id", n.ID, "url", n.Subject.LatestCommentURL)

	resp, err := c.API.Request(http.MethodGet, n.Subject.LatestCommentURL, nil)
	if err != nil {
		return notifications.User{}, fmt.Errorf("failed to get last commentor: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return notifications.User{}, fmt.Errorf("failed to read last commentor: %w", err)
	}

	comment := struct {
		User notifications.User `json:"user"`
	}{}

	err = json.Unmarshal(body, &comment)
	if err != nil {
		return notifications.User{}, fmt.Errorf("failed to unmarshal last commentor: %w", err)
	}

	return comment.User, nil
}
