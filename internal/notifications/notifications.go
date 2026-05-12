// Packaeg notifications handles storing notification types and provides some
// helpers.
//
// Reference: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28

//revive:disable:max-public-structs // TODO: split this into multiple files.
package notifications

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
)

type (
	Notifications   []*Notification
	NotificationMap map[string]*Notification
)

type Notification struct {
	// Standard API fields
	ID         string     `json:"id"`
	Unread     bool       `json:"unread"`
	Reason     string     `json:"reason"`
	UpdatedAt  time.Time  `json:"updated_at"`
	URL        string     `json:"url"`
	Repository Repository `json:"repository"`
	Subject    Subject    `json:"subject"`

	// Enriched API fields
	Author          User   `json:"author"`
	LatestCommentor User   `json:"latest_commentor"`
	Assignees       []User `json:"assignees"`
	Reviewers       []User `json:"requested_reviewers"`
	ReviewersTeams  []Team `json:"requested_teams"`
	MergedBy        User   `json:"merged_by"`

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

	// Enriched marks notifications whose enriched API fields are cached.
	Enriched bool `json:"enriched"`

	// Tags is a list of tags that can be used to filter notifications.
	// They can be added/removed with the `tag` action.
	Tags []string `json:"tags"`
}

type Subject struct {
	// Standard API fields
	Title            string `json:"title"`
	URL              string `json:"url"`
	Type             string `json:"type"`
	LatestCommentURL string `json:"latest_comment_url"`

	// Enriched API fields
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
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

type Team struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
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

//nolint:cyclop // TODO: add sub struct equality.
//revive:disable:cyclomatic
func (n *Notification) Equal(other *Notification) bool {
	return n.ID == other.ID &&
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
		n.Subject.HTMLURL == other.Subject.HTMLURL &&
		n.Author.Login == other.Author.Login &&
		n.Author.Type == other.Author.Type &&
		n.LatestCommentor.Login == other.LatestCommentor.Login &&
		n.LatestCommentor.Type == other.LatestCommentor.Type &&
		n.Meta.Hidden == other.Meta.Hidden &&
		n.Meta.Done == other.Meta.Done &&
		n.Meta.RemoteExists == other.Meta.RemoteExists &&
		n.Meta.Enriched == other.Meta.Enriched
}

// MergeUpdatedNotification copies cached metadata from n onto remote. If remote
// is newer, it resets Done/Enriched and drops stale enrichment; otherwise it
// preserves cached enrichment. It mutates and returns remote.
func (n *Notification) MergeUpdatedNotification(remote *Notification) *Notification {
	meta := n.Meta
	meta.RemoteExists = true

	if remote.UpdatedAt.After(n.UpdatedAt) {
		meta.Done = false
		meta.Enriched = false

		remote.clearEnrichment()
	} else if meta.Enriched {
		remote.copyEnrichmentFrom(n)
	}

	remote.Meta = meta

	return remote
}

// BackfillEnriched sets Meta.Enriched on cached notifications saved before the
// field existed. Without this, every notification would re-enrich on the first
// refresh after upgrade.
func (n Notifications) BackfillEnriched() {
	for _, notification := range n {
		if notification != nil {
			notification.BackfillEnriched()
		}
	}
}

// BackfillEnriched sets Meta.Enriched=true when n carries enrichment-only
// subject fields. No-op if the flag is already set or no signal is present.
func (n *Notification) BackfillEnriched() {
	if n.Meta.Enriched {
		return
	}

	if n.Subject.State == "" && n.Subject.HTMLURL == "" {
		return
	}

	n.Meta.Enriched = true
}

func (n *Notification) Marshal() ([]byte, error) {
	marshaled, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal notifications: %w", err)
	}

	return marshaled, nil
}

func (n *Notification) Interface() (any, error) {
	marshaled, err := n.Marshal()
	if err != nil {
		return nil, err
	}

	var i any
	if err := json.Unmarshal(marshaled, &i); err != nil {
		return nil, fmt.Errorf("cannot convert to interface: %w", err)
	}

	return i, nil
}

func (n Notifications) Debug() string {
	out := make([]string, 0, len(n))
	for _, n := range n {
		out = append(out, n.Debug())
	}

	return strings.Join(out, "\n")
}

func (n *Notification) Debug() string {
	return fmt.Sprintf("%#v", *n)
}

func (n *Notification) copyEnrichmentFrom(other *Notification) {
	n.Author = other.Author
	n.LatestCommentor = other.LatestCommentor
	n.Assignees = other.Assignees
	n.Reviewers = other.Reviewers
	n.ReviewersTeams = other.ReviewersTeams
	n.MergedBy = other.MergedBy
	n.Subject.State = other.Subject.State
	n.Subject.HTMLURL = other.Subject.HTMLURL
}

func (n *Notification) clearEnrichment() {
	n.copyEnrichmentFrom(&Notification{})
}

func (n Notifications) Map() NotificationMap {
	m := NotificationMap{}
	for _, n := range n {
		m[n.ID] = n
	}

	return m
}

func (m NotificationMap) List() Notifications {
	l := make(Notifications, 0, len(m))
	for _, n := range m {
		l = append(l, n)
	}

	return l
}

func (n Notifications) IDList() []string {
	ids := make([]string, 0, len(n))
	for _, n := range n {
		ids = append(ids, n.ID)
	}

	return ids
}

// Compact remove all nil notifications.
// TODO: in-place update.
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

// Uniq remove all duplicated notifications.
// TODO: in-place update.
func (n Notifications) Uniq() Notifications {
	seenIDs := map[string]bool{}

	return slices.DeleteFunc(n, func(n *Notification) bool {
		if _, ok := seenIDs[n.ID]; ok {
			return true
		}

		seenIDs[n.ID] = true

		return false
	})
}

func (n Notifications) FilterFromIDs(ids []string) Notifications {
	newList := Notifications{}

	for _, id := range ids {
		for _, n := range n {
			if n.ID == id {
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

func (n Notifications) Interface() (any, error) {
	marshaled, err := n.Marshal()
	if err != nil {
		return nil, err
	}

	var i any
	if err := json.Unmarshal(marshaled, &i); err != nil {
		return nil, fmt.Errorf("cannot convert to interface: %w", err)
	}

	return i, nil
}
