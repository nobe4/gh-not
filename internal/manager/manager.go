package manager

import (
	"fmt"
	"log/slog"
	"os"

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
	config        *config.Data
	client        *gh.Client
	Actors        actors.ActorsMap
	refresh       RefreshStrategy
}

func New(config *config.Data) *Manager {
	m := &Manager{}

	m.config = config
	m.cache = cache.NewFileCache(m.config.Cache.TTLInHours, m.config.Cache.Path)

	return m
}

func (m *Manager) WithCaller(caller api.Caller) *Manager {
	m.client = gh.NewClient(caller, m.cache, m.config.Endpoint)
	m.Actors = actors.Map(m.client)

	return m
}

func (m *Manager) WithRefresh(refresh RefreshStrategy) *Manager {
	m.refresh = refresh

	return m
}

func (m *Manager) Load() error {
	m.Notifications = notifications.Notifications{}

	cachedNotifications, expired, err := m.loadCache()
	if err != nil {
		slog.Warn("cannot read the cache: %#v\n", err)
	} else if cachedNotifications != nil {
		m.Notifications = cachedNotifications
	}

	if m.shouldRefresh(expired) {
		return m.refreshNotifications()
	}

	return nil
}

func (m *Manager) shouldRefresh(expired bool) bool {
	if !expired && m.refresh == ForceRefresh {
		slog.Info("forcing a refresh")
		return true
	}

	if expired && m.refresh == PreventRefresh {
		slog.Info("preventing a refresh")
		return false
	}

	slog.Debug("refresh", "refresh", expired)
	return expired
}

func (m *Manager) refreshNotifications() error {
	fmt.Printf("Refreshing the cache...\n")

	remoteNotifications, err := m.client.Notifications()
	if err != nil {
		return err
	}

	m.Notifications = notifications.Sync(m.Notifications, remoteNotifications)

	if err := m.cache.Write(m.Notifications); err != nil {
		slog.Error("Error while writing the cache: %#v", err)
	}

	m.Notifications = m.Notifications.Uniq()

	m.Notifications, err = m.client.Enrich(m.Notifications)

	return err
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

		for _, notification := range m.Notifications.FilterFromIds(selectedIds) {
			// TODO: add --force flag to ignore this
			if notification.Meta.Done {
				continue
			}

			if noop {
				fmt.Printf("NOOP'ing action %s on notification %s\n", rule.Action, notification.String())
				continue
			}

			if err := actor.Run(notification, os.Stdout); err != nil {
				slog.Error("action failed", "action", rule.Action, "err", err)
			}
			fmt.Fprintln(os.Stdout, "")
		}
	}

	m.Notifications = m.Notifications.Compact()

	return nil
}
