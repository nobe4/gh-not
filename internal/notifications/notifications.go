package notifications

import (
	"encoding/json"
	"fmt"
)

type NotificationMap map[string]Notification

type Notifications []Notification

type Notification struct {
	// API fields
	Title      string     `json:"title"`
	Id         string     `json:"id"`
	Unread     bool       `json:"unread"`
	Repository Repository `json:"repository"`

	// gh-not fields
	Meta Meta `json:"meta"`
}

type Meta struct {
	Hidden bool `json:"hidden"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`

	Owner Owner `json:"owner"`
}

type Owner struct {
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
	return fmt.Sprintf("%s %+v", n.Id, n)
}

func (n NotificationMap) ToString() string {
	out := ""
	for _, n := range n {
		out += n.ToString() + "\n"
	}
	return out
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
