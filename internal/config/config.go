package config

import (
	"log"
	"os"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Cache  Cache   `yaml:"cache"`
	Groups []Group `yaml:"groups"`
}

type Cache struct {
	TTLInHours int    `yaml:"ttl_in_hours"`
	Path       string `yaml:"path"`
}

type Group struct {
	Name    string   `yaml:"name"`
	Filters []string `yaml:"filters"`
	Action  string   `yaml:"action"`
}

func New(path string) (*Config, error) {
	config := &Config{
		// default values
		Cache: Cache{
			TTLInHours: 24 * 7,
			Path:       "/tmp/gh-not.cache.json",
		},
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Apply(n notifications.NotificationMap, actors map[string]actors.Actor) (notifications.NotificationMap, error) {
	err := error(nil)
	for _, group := range c.Groups {
		selectedNotifications := n.ToSlice()

		for _, filter := range group.Filters {
			selectedNotifications, err = jq.Filter(filter, selectedNotifications)
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
