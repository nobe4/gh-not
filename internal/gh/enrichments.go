package gh

import (
	"encoding/json"
	"fmt"
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
	MergedBy       notifications.User   `json:"merged_by"`
}

func (c *Client) Enrich(n *notifications.Notification) error {
	if n == nil {
		return nil
	}

	threadExtra, err := c.getThreadExtra(n)
	if err != nil {
		return err
	}

	lastCommentor, err := c.getLastCommentor(n)
	if err != nil {
		return err
	}

	n.Author = threadExtra.User
	n.Subject.State = threadExtra.State
	n.Subject.HTMLURL = threadExtra.HTMLURL
	n.Assignees = threadExtra.Assignees
	n.Reviewers = threadExtra.Reviewers
	n.ReviewersTeams = threadExtra.ReviewersTeams
	n.MergedBy = threadExtra.MergedBy
	n.LatestCommentor = lastCommentor
	n.Meta.Enriched = true

	return nil
}

func (c *Client) getJSON(url string, v any) error {
	resp, err := c.API.Request(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to get %s: %w", url, err)
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v) //nolint:wrapcheck // This is wrapped by the caller
}

func (c *Client) getThreadExtra(n *notifications.Notification) (ThreadExtra, error) {
	if n.Subject.URL == "" {
		return ThreadExtra{}, nil
	}

	slog.Debug("getting the thread extra", "id", n.ID, "url", n.Subject.URL)

	extra := ThreadExtra{}
	if err := c.getJSON(n.Subject.URL, &extra); err != nil {
		return ThreadExtra{}, fmt.Errorf("failed to get thread extra: %w", err)
	}

	return extra, nil
}

func (c *Client) getLastCommentor(n *notifications.Notification) (notifications.User, error) {
	if n.Subject.LatestCommentURL == "" {
		return notifications.User{}, nil
	}

	slog.Debug("getting the last commentor", "id", n.ID, "url", n.Subject.LatestCommentURL)

	comment := struct {
		User notifications.User `json:"user"`
	}{}
	if err := c.getJSON(n.Subject.LatestCommentURL, &comment); err != nil {
		return notifications.User{}, fmt.Errorf("failed to get last commentor: %w", err)
	}

	return comment.User, nil
}
