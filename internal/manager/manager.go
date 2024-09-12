package manager

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nobe4/gh-not/internal/actions"
	"github.com/nobe4/gh-not/internal/api"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Manager struct {
	Notifications notifications.Notifications
	Cache         cache.ExpiringReadWriter
	config        *config.Data
	client        *gh.Client
	Actions       actions.ActionsMap

	RefreshStrategy RefreshStrategy
	ForceStrategy   ForceStrategy
}

func New(config *config.Data) *Manager {
	m := &Manager{}

	m.config = config
	m.Cache = cache.NewFileCache(m.config.Cache.TTLInHours, m.config.Cache.Path)

	return m
}

func (m *Manager) SetCaller(caller api.Requestor) {
	m.client = gh.NewClient(caller, m.Cache, m.config.Endpoint)
	m.Actions = actions.Map(m.client)
}

func (m *Manager) Load() error {
	if err := m.Cache.Read(&m.Notifications); err != nil {
		slog.Warn("cannot read the cache", "error", err)
	}

	slog.Info("Loaded notifications", "count", len(m.Notifications))

	return nil
}

func (m *Manager) Refresh() error {
	if m.RefreshStrategy.ShouldRefresh(m.Cache.Expired()) {
		return m.refreshNotifications()
	}

	slog.Info("Refreshed notifications", "count", len(m.Notifications))

	return nil
}

func (m *Manager) refreshNotifications() error {
	if m.client == nil {
		return fmt.Errorf("manager has no client, cannot refresh notifications")
	}

	fmt.Printf("Refreshing notifications...\n")

	remoteNotifications, err := m.client.Notifications()
	if err != nil {
		return err
	}

	m.Notifications = notifications.Sync(m.Notifications, remoteNotifications)
	m.Notifications = m.Notifications.Uniq()
	m.Notifications, err = m.Enrich(m.Notifications)

	return err
}

func (m *Manager) Save() error {
	return m.Cache.Write(m.Notifications.Compact())
}

func (m *Manager) Enrich(ns notifications.Notifications) (notifications.Notifications, error) {
	for i, n := range ns {
		if n.Meta.Done && !m.ForceStrategy.Has(ForceEnrich) {
			continue
		}

		if err := m.client.Enrich(n); err != nil {
			return nil, err
		}

		ns[i] = n
	}

	return ns, nil
}

func (m *Manager) Apply() error {
	for _, rule := range m.config.Rules {
		runner, ok := m.Actions[rule.Action]
		if !ok {
			slog.Error("unknown action", "action", rule.Action)
			continue
		}

		selectedNotifications, err := rule.Filter(m.Notifications)
		if err != nil {
			return err
		}

		slog.Debug("apply rule", "name", rule.Name, "count", len(selectedNotifications))

		for _, notification := range selectedNotifications {
			if notification.Meta.Done && !m.ForceStrategy.Has(ForceApply) {
				slog.Debug("skipping done notification", "id", notification.Id)
				continue
			}

			if m.ForceStrategy.Has(ForceNoop) {
				fmt.Printf("NOOP'ing action %s on notification %s\n", rule.Action, notification.String())
				continue
			}

			if err := runner.Run(notification, os.Stdout); err != nil {
				slog.Error("action failed", "action", rule.Action, "err", err)
			}
			fmt.Fprintln(os.Stdout, "")
		}
	}

	m.Notifications = m.Notifications.Compact()

	return nil
}
