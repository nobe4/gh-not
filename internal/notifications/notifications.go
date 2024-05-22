// Handles storing notification types and provides some helpers.
//
// Reference: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28

package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/nobe4/gh-not/internal/colors"
)

type NotificationMap map[string]Notification

type Notifications []Notification

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
	Meta Meta `json:"meta"`
}

type Meta struct {
	Hidden bool `json:"hidden"`
}

type Subject struct {
	// Standard API fields
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`

	// Enriched API fields
	State string `json:"state"`
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

func (n NotificationMap) Append(notifications []Notification) {
	for _, notification := range notifications {
		if _, ok := n[notification.Id]; !ok {
			n[notification.Id] = notification
		}
	}
}

func (n NotificationMap) Uniq() NotificationMap {
	unique := NotificationMap{}

	for _, notification := range n {
		if _, ok := unique[notification.Id]; !ok {
			unique[notification.Id] = notification
		}
	}

	return unique
}

func (n Notification) ToString() string {
	return fmt.Sprintf("%s %s %s by %s: '%s' ", n.prettyType(), n.prettyState(), n.Repository.FullName, n.Author.Login, n.Subject.Title)
}

var prettyTypes = map[string]string{
	"Issue":       colors.Blue("IS"),
	"PullRequest": colors.Cyan("PR"),
}

var prettyState = map[string]string{
	"open":   colors.Green("OP"),
	"closed": colors.Red("CL"),
	"merged": colors.Magenta("MG"),
}

func (n Notification) prettyType() string {
	if p, ok := prettyTypes[n.Subject.Type]; ok {
		return p
	}

	return "T?"
}

func (n Notification) prettyState() string {
	if p, ok := prettyState[n.Subject.State]; ok {
		return p
	}

	return "S?"
}

func (n NotificationMap) ToString() string {
	out := ""
	for _, n := range n {
		out += n.ToString() + "\n"
	}
	return out
}

func (n NotificationMap) ToTable() (string, error) {
	out := bytes.Buffer{}

	t := term.FromEnv()
	w, _, err := t.Size()
	if err != nil {
		return "", err
	}

	printer := tableprinter.New(&out, t.IsTerminalOutput(), w)

	for _, n := range n {
		printer.AddField(n.prettyType())
		printer.AddField(n.prettyState())
		printer.AddField(n.Repository.FullName)
		printer.AddField(n.Author.Login)
		printer.AddField(n.Subject.Title)
		printer.EndRow()
	}

	if err := printer.Render(); err != nil {
		return "", err
	}

	fmt.Fprintf(&out, "Found %d notifications", len(n))

	return out.String(), nil
}

func (n NotificationMap) ToSlice() Notifications {
	s := Notifications{}

	for _, n := range n {
		s = append(s, n)
	}

	return s
}

func (n Notifications) ToInterface() (interface{}, error) {
	marshalled, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal notifications: %w", err)
	}

	var i interface{}
	if err := json.Unmarshal(marshalled, &i); err != nil {
		return nil, fmt.Errorf("cannot unmarshal interface: %w", err)
	}

	return i, nil
}

func FromInterface(i interface{}) (Notifications, error) {
	marshalled, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("cannot marshall interface: %w", err)
	}

	notifications := Notifications{}
	if err := json.Unmarshal(marshalled, &notifications); err != nil {
		return nil, fmt.Errorf("cannot unmarshall into notification: %w", err)
	}

	return notifications, nil
}
