package manager

import (
	"fmt"
	"log/slog"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/api"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Manager struct {
	Notifications notifications.Notifications
	cache         cache.ExpiringReadWriter
	config        *config.Config
	client        *gh.Client
	Actors        actors.ActorsMap
}

func New(config *config.Config, caller api.Caller) *Manager {
	m := &Manager{}

	m.config = config
	m.cache = cache.NewFileCache(m.config.Cache.TTLInHours, m.config.Cache.Path)
	m.client = gh.NewClient(caller, m.cache)
	m.Actors = actors.Map(m.client)

	return m
}

func (m *Manager) Load(refresh, noRefresh bool) error {
	allNotifications := notifications.Notifications{}

	cachedNotifications, refreshCache, err := m.loadCache()
	if err != nil {
		slog.Warn("cannot read the cache: %#v\n", err)
	} else if cachedNotifications != nil {
		allNotifications = cachedNotifications
	}

	if !refreshCache && refresh {
		slog.Info("forcing a refresh")
		refresh = true
	}

	if refreshCache && noRefresh {
		slog.Info("preventing a refresh")
		refresh = false
	}

	if refresh {
		fmt.Printf("Refreshing the cache...\n")

		pulledNotifications, err := m.client.Notifications(refresh, noRefresh)
		if err != nil {
			return err
		}

		allNotifications = append(allNotifications, pulledNotifications...)

		if err := m.cache.Write(allNotifications); err != nil {
			slog.Error("Error while writing the cache: %#v", err)
		}
	}

	m.Notifications = allNotifications.Uniq()

	return nil
}

func (m *Manager) Save() error {
	return m.cache.Write(m.Notifications.Compact())
}

func (m *Manager) loadCache() (notifications.Notifications, bool, error) {
	expired, err := m.cache.Expired()
	if err != nil {
		return nil, false, err
	}

	n := notifications.Notifications{}
	if err := m.cache.Read(&n); err != nil {
		return nil, expired, err
	}

	return n, expired, nil
}

func (m *Manager) Apply(noop bool) error {
	for _, rule := range m.config.Rules {
		actor, ok := m.Actors[rule.Action]
		if !ok {
			slog.Error("unknown action", "action", rule.Action)
			continue
		}

		selectedIds, err := rule.FilterIds(m.Notifications)
		if err != nil {
			return err
		}

		slog.Debug("apply rule", "name", rule.Name, "count", len(selectedIds))

		var out string

		for _, notification := range m.Notifications.FilterFromIds(selectedIds) {
			if noop {
				fmt.Printf("NOOP'ing action %s on notification %s\n", rule.Action, notification.ToString())
				continue
			}

			out, err = actor.Run(notification)
			if err != nil {
				slog.Error("action failed", "action", rule.Action, "err", err)
			}
			if out != "" {
				fmt.Println(out)
			}
		}
	}

	m.Notifications = m.Notifications.Compact()

	return nil
}
