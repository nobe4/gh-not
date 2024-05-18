package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Config struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Name    string   `json:"name"`
	Filters []string `json:"filters"`
	Action  string   `json:"action"`
}

func New(path string) (*Config, error) {
	config := &Config{}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Apply(n notifications.NotificationMap, actors map[string]actors.Actor) (notifications.NotificationMap, error) {
	err := error(nil)
	for _, group := range c.Groups {
		notificationList := n.ToSlice()
		selectedNotifications := notifications.Notifications{}

		for _, filter := range group.Filters {
			selectedNotifications, err = jq.Filter(filter, notificationList)
			if err != nil {
				return nil, err
			}
		}

		for _, notification := range selectedNotifications {
			if actor, ok := actors[group.Action]; ok == true {
				notification, err = actor.Run(notification)
				if err != nil {
					log.Fatalf("action '%s' failed: %v", group.Action, err)
				}

				n[notification.Id] = notification
			} else {
				log.Fatalf("unknown action '%s'", group.Action)
			}
		}
	}

	return n, nil
}
