// Handles storing notification types and provides some helpers.
//
// Reference: https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28

package notifications

import (
	"encoding/json"
	"fmt"
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
	Meta Meta `json:"meta"`
}

type Meta struct {
	Hidden bool `json:"hidden"`

	// TODO: Rename to `Done`
	ToDelete bool `json:"to_delete"`
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

// Sync merges the local and remote notifications.
//
// It applies the following rules:
// | remote \ local | Missing   | Open      | Done      | Hidden    |
// | ---            | ---       | ---       | ---       | ---       |
// | Exist          | (1)Insert | (2)Update | (2)Update | (2)Update |
// | Missing        | (3)Noop   | (3)Noop   | (4)Drop   | (2)Update |
//
// 1. Insert: Add the notification to the new.
// 2. Update: Update the local notification with the remote data.
// 3. Noop: Do nothing.
// 4. Drop: Remove the notification from the local list.
func Sync(local, remote Notifications) Notifications {
	remoteMap := remote.Map()
	localMap := local.Map()

	n := Notifications{}

	// Add any new notifications to the list
	for remoteId, remote := range remoteMap {
		if _, ok := localMap[remoteId]; !ok {
			// (1)Insert
			n = append(n, remote)
		}
	}

	for localId, local := range localMap {
		remote, remoteExist := remoteMap[localId]

		if remoteExist {
			// (2)Update
			remote.Meta = local.Meta
			n = append(n, remote)
		} else {
			if local.Meta.ToDelete {
				// (4)Drop
				continue
			}

			// (3)Noop
			n = append(n, local)
		}
	}

	n.Sort()

	return n
}

func (n Notifications) IDList() []string {
	ids := []string{}
	for _, n := range n {
		ids = append(ids, n.Id)
	}
	return ids
}

func (n Notifications) Compact() Notifications {
	return slices.DeleteFunc(n, func(n *Notification) bool {
		return n == nil || n.Meta.ToDelete
	})
}

func (n Notifications) Sort() {
	slices.SortFunc(n, func(a, b *Notification) int {
		return strings.Compare(a.Id, b.Id)
	})
}

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

func (n Notifications) Marshall() (interface{}, error) {
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
