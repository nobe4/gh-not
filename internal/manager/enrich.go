package manager

import (
	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/nobe4/gh-not/internal/notifications"
)

func (m *Manager) Enrich(ns notifications.Notifications) {
	g := new(errgroup.Group)
	g.SetLimit(m.enrichWorkers())

	for _, n := range ns {
		if !m.shouldEnrich(n) {
			continue
		}

		g.Go(func() error {
			if err := m.client.Enrich(n); err != nil {
				slog.Warn("failed to enrich notification", "notification", n.ID, "error", err.Error())
			}

			return nil
		})
	}

	//nolint:errcheck // We don't do anything with the errgroup's final error.
	g.Wait()
}

func (m *Manager) enrichWorkers() int {
	if m.config == nil || m.config.Enrichment.Workers < 1 {
		return 1
	}

	return m.config.Enrichment.Workers
}

func (m *Manager) shouldEnrich(notification *notifications.Notification) bool {
	if notification == nil {
		return false
	}

	if m.ForceStrategy.Has(ForceEnrich) {
		return true
	}

	if notification.Meta.Enriched {
		return false
	}

	if notification.Meta.Done {
		return false
	}

	return true
}
