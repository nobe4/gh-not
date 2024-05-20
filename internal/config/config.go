package config

import (
	"log/slog"
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

var (
	defaultCache = Cache{
		TTLInHours: 24 * 7,
		Path:       "/tmp/gh-not.cache.json",
	}
)

func New(path string) (*Config, error) {
	config := &Config{Cache: defaultCache}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Apply(n notifications.NotificationMap, actors map[string]actors.Actor, noop bool) (notifications.NotificationMap, error) {
	err := error(nil)

	for _, group := range c.Groups {
		slog.Debug("apply group", "name", group.Name)
		selectedNotifications := n.ToSlice()

		for _, filter := range group.Filters {
			slog.Debug("apply filter", "filter", filter)
			selectedNotifications, err = jq.Filter(filter, selectedNotifications)
			if err != nil {
				return nil, err
			}
		}

		for _, notification := range selectedNotifications {
			if actor, ok := actors[group.Action]; ok == true {
				if noop {
					slog.Info("NOOP'ing", "action", group.Action, "notification", notification.ToString())
				} else {
					// Remove the notification temporarily from the list, it
					// will be added back after the actor runs.
					delete(n, notification.Id)

					notification, err = actor.Run(notification)
					if err != nil {
						slog.Error("action failed", "action", group.Action, "err", err)
					}

					slog.Debug("returned notification", "notification", notification.ToString())
					n[notification.Id] = notification
				}
			} else {
				slog.Error("unknown action", "action", group.Action)
			}
		}
	}

	return n, nil
}
