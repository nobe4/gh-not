package notifications

import (
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type NotificationsManager interface {
	List() notifications.Notifications
	Get(string) (*notifications.Notification, error)
}

type CachedNotificationsManager struct {
	cache  cache.ExpiringReadWriter
	client gh.NotificationCaller

	loaded bool

	list notifications.Notifications
}

func (m *CachedNotificationsManager) load() error {
	expired, err := m.cache.Expired()
	if err != nil {
		return err
	}

	if expired {
		// TODO: pass list as parameter instead of return
		list, err := m.client.List()
		if err != nil {
			return err
		}
		m.list = list
	} else {
		if err := m.cache.Read(&m.list); err != nil {
			return err
		}
	}

	m.loaded = true
	return nil
}

func (m *CachedNotificationsManager) Get(id string) (notifications.Notification, error) {
	if !m.loaded {
		if err := m.load(); err != nil {
			return notifications.Notification{}, err
		}
	}

	for _, n := range m.list {
		if id == n.Id {
			return n, nil
		}
	}

	return nil, nil
}
