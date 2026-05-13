package manager

import (
	"log/slog"
	"sync"

	"github.com/nobe4/gh-not/internal/notifications"
)

func (m *Manager) Enrich(ns notifications.Notifications) {
	workers := m.enrichWorkers()

	var wg sync.WaitGroup

	sem := make(chan struct{}, workers)

	for _, n := range ns {
		if !m.shouldEnrich(n) {
			continue
		}

		sem <- struct{}{}

		wg.Go(func() {
			defer func() { <-sem }()

			if err := m.client.Enrich(n); err != nil {
				// Enrichment of a single notification should not prevent the
				// enrichment to continue.
				// TODO: suggest to re-run the enrichment
				slog.Warn("failed to enrich notification", "notification", n.ID, "error", err.Error())
			}
		})
	}

	wg.Wait()
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
