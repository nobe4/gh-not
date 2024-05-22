package config

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Cache Cache  `yaml:"cache"`
	Rules []Rule `yaml:"rules"`
}

type Cache struct {
	TTLInHours int    `yaml:"ttl_in_hours"`
	Path       string `yaml:"path"`
}

type Rule struct {
	Name    string   `yaml:"name"`
	Filters []string `yaml:"filters"`
	Action  string   `yaml:"action"`
}

var (
	defaultCache = Cache{
		TTLInHours: 24 * 7,
		Path:       path.Join(StateDir(), "cache.json"),
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

	for _, rule := range c.Rules {
		slog.Debug("apply rule", "name", rule.Name)
		selectedNotifications := n.ToSlice()

		for _, filter := range rule.Filters {
			selectedNotifications, err = jq.Filter(filter, selectedNotifications)
			if err != nil {
				return nil, err
			}
		}

		for _, notification := range selectedNotifications {
			if actor, ok := actors[rule.Action]; ok == true {
				if noop {
					fmt.Printf("NOOP'ing action %s on notification %s\n", rule.Action, notification.ToString())
				} else {
					// Remove the notification temporarily from the list, it
					// will be added back after the actor runs.
					delete(n, notification.Id)

					notification, err = actor.Run(notification)
					if err != nil {
						slog.Error("action failed", "action", rule.Action, "err", err)
					}

					n[notification.Id] = notification
				}
			} else {
				slog.Error("unknown action", "action", rule.Action)
			}
		}
	}

	return n, nil
}
