package manager

import (
	"log/slog"
	"sync"

	"github.com/nobe4/gh-not/internal/notifications"
)

const enrichWorkers = 10

func (m *Manager) Enrich(ns notifications.Notifications) notifications.Notifications {
	var wg sync.WaitGroup
	sem := make(chan struct{}, enrichWorkers)

	for _, n := range ns {
		if !m.shouldEnrich(n) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			m.enrichNotification(n)
		}()
	}

	wg.Wait()

	return ns
}

func (m *Manager) shouldEnrich(notification *notifications.Notification) bool {
	if notification == nil {
		return false
	}

	if m.ForceStrategy.Has(ForceEnrich) {
		return true
	}

	if notification.Meta.Done {
		return false
	}

	// Skip notifications that already carry enriched data. Sync preserves
	// these fields when UpdatedAt hasn't moved, so the values are still
	// fresh and we can avoid the API call.
	if notification.Author.Login != "" {
		return false
	}

	return true
}

func (m *Manager) enrichNotification(notification *notifications.Notification) {
	if err := m.client.Enrich(notification); err != nil {
		// Enrichment of a single notification should not prevent the
		// enrichment to continue.
		// TODO: suggest to re-run the enrichment
		slog.Warn("failed to enrich notification", "notification", notification.ID, "error", err.Error())
	}
}
