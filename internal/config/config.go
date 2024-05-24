package config

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/nobe4/gh-not/internal/actors"
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

func (c *Config) Apply(n notifications.Notifications, actors map[string]actors.Actor, noop bool) (notifications.Notifications, error) {
	indexIDMap := map[string]int{}
	for i, n := range n {
		indexIDMap[n.Id] = i
	}

	for _, rule := range c.Rules {
		slog.Debug("apply rule", "name", rule.Name)

		selectedIds, err := rule.filterIds(n)
		if err != nil {
			return nil, err
		}

		for _, id := range selectedIds {
			i := indexIDMap[id]
			notification := n[i]

			if actor, ok := actors[rule.Action]; ok {
				if noop {
					fmt.Printf("NOOP'ing action %s on notification %s\n", rule.Action, notification.ToString())
				} else {
					notification, err = actor.Run(notification)
					if err != nil {
						slog.Error("action failed", "action", rule.Action, "err", err)
					}
				}
			} else {
				slog.Error("unknown action", "action", rule.Action)
			}

			n[i] = notification
		}
	}

	return n.DeleteNil(), nil
}
