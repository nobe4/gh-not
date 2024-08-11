package tests

import (
	"testing"

	"github.com/nobe4/gh-not/internal/logger"
	"github.com/nobe4/gh-not/internal/manager"
	"github.com/nobe4/gh-not/test/integration"
)

func Test000(t *testing.T) {
	logger.Init(5)

	m := integration.MockManager(t, integration.MockManagerConfig{
		ConfigPath:      "./000_config.yaml",
		CallsPath:       "./000_calls.json",
		RefreshStrategy: manager.ForceRefresh,
	})

	// Test the manager
	if m == nil {
		t.Fatal("manager is nil")
	}

	if err := m.Load(); err != nil {
		t.Fatal(err)
	}

	if len(m.Notifications) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(m.Notifications))
	}

	if err := m.Refresh(); err != nil {
		t.Fatal(err)
	}

	integration.Compare(t, "./000_expected.json", m.Notifications)
}
