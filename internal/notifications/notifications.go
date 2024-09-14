// Handles storing notification types and provides some helpers.
//
// Reference: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28

package notifications

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
)

type Notifications []*Notification
type NotificationMap map[string]*Notification

type Notification struct {
	// Standard API fields
	Id         string     `json:"id"`
	Unread     bool       `json:"unread"`
	Reason     string     `json:"reason"`
	UpdatedAt  time.Time  `json:"updated_at"`
	URL        string     `json:"url"`
	Repository Repository `json:"repository"`
	Subject    Subject    `json:"subject"`

	// Enriched API fields
	Author User `json:"author"`

	// gh-not specific fields
	// Those fields are not part of the GitHub API and will persist between
	// syncs.
	Meta Meta `json:"meta"`

	// Rendered string for display, set by Notifications.Render
	rendered string
}

type Meta struct {
	// Hide the notification from the user
	Hidden bool `json:"hidden"`

	// Mark the notification as done, will be deleted as soon as it's missing
	// from the remote notification list.
	Done bool `json:"done"`

	// RemoteExists is used to track if the notification is still present in the
	// remote list (i.e. the API).
	RemoteExists bool `json:"remote_exists"`
}

type Subject struct {
	// Standard API fields
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`

	// Enriched API fields
	State   string `json:"state"`
	HtmlUrl string `json:"html_url"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`

	Owner User `json:"owner"`
}

type User struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

func (n Notifications) Equal(others Notifications) bool {
	if len(n) != len(others) {
		return false
	}

	for i, n := range n {
		if !n.Equal(others[i]) {
			slog.Info("notification not equal", "n", n.Debug(), "other", others[i].Debug())
			return false
		}
	}

	return true
}

func (n Notification) Equal(other *Notification) bool {
	return n.Id == other.Id &&
		n.Unread == other.Unread &&
		n.Reason == other.Reason &&
		n.UpdatedAt.Equal(other.UpdatedAt) &&
		n.URL == other.URL &&
		n.Repository.Name == other.Repository.Name &&
		n.Repository.FullName == other.Repository.FullName &&
		n.Repository.Private == other.Repository.Private &&
		n.Repository.Fork == other.Repository.Fork &&
		n.Repository.Owner.Login == other.Repository.Owner.Login &&
		n.Repository.Owner.Type == other.Repository.Owner.Type &&
		n.Subject.Title == other.Subject.Title &&
		n.Subject.URL == other.Subject.URL &&
		n.Subject.Type == other.Subject.Type &&
		n.Subject.State == other.Subject.State &&
		n.Subject.HtmlUrl == other.Subject.HtmlUrl &&
		n.Author.Login == other.Author.Login &&
		n.Author.Type == other.Author.Type &&
		n.Meta.Hidden == other.Meta.Hidden &&
		n.Meta.Done == other.Meta.Done &&
		n.Meta.RemoteExists == other.Meta.RemoteExists
}

func (n Notifications) Debug() string {
	out := []string{}
	for _, n := range n {
		out = append(out, n.Debug())
	}
	return strings.Join(out, "\n")
}

func (n Notification) Debug() string {
	return fmt.Sprintf("%#v", n)
}

func (n Notifications) Map() NotificationMap {
	m := NotificationMap{}
	for _, n := range n {
		m[n.Id] = n
	}
	return m
}

func (m NotificationMap) List() Notifications {
	l := Notifications{}
	for _, n := range m {
		l = append(l, n)
	}
	return l
}

func (n Notifications) IDList() []string {
	ids := []string{}
	for _, n := range n {
		ids = append(ids, n.Id)
	}
	return ids
}

// TODO: in-place update
func (n Notifications) Compact() Notifications {
	return slices.DeleteFunc(n, func(n *Notification) bool {
		return n == nil
	})
}

func (n Notifications) Sort() {
	slices.SortFunc(n, func(a, b *Notification) int {
		if a.UpdatedAt.Before(b.UpdatedAt) {
			return 1
		} else if a.UpdatedAt.After(b.UpdatedAt) {
			return -1
		}
		return 0
	})
}

// TODO: in-place update
func (n Notifications) Uniq() Notifications {
	seenIds := map[string]bool{}
	return slices.DeleteFunc(n, func(n *Notification) bool {
		if _, ok := seenIds[n.Id]; ok {
			return true
		}
		seenIds[n.Id] = true
		return false
	})
}

func (n Notifications) FilterFromIds(ids []string) Notifications {
	newList := Notifications{}

	for _, id := range ids {
		for _, n := range n {
			if n.Id == id {
				newList = append(newList, n)
			}
		}
	}

	return newList
}

func (n Notifications) Marshal() ([]byte, error) {
	marshaled, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal notifications: %w", err)
	}

	return marshaled, nil
}

func (n Notifications) Interface() (interface{}, error) {
	marshaled, err := n.Marshal()
	if err != nil {
		return nil, err
	}

	var i interface{}
	if err := json.Unmarshal(marshaled, &i); err != nil {
		return nil, fmt.Errorf("cannot convert to interface: %w", err)
	}

	return i, nil
}
