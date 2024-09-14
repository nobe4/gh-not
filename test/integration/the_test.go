package tests

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"

	apiMock "github.com/nobe4/gh-not/internal/api/mock"
	configPkg "github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/logger"
	"github.com/nobe4/gh-not/internal/manager"
	"github.com/nobe4/gh-not/internal/notifications"
)

type config struct {
	Id string
	// TODO: move those into config so it can be set by default as well as via
	// CLI
	ForceStrategy   manager.ForceStrategy
	RefreshStrategy manager.RefreshStrategy
}

func setup(t *testing.T, conf config) (*manager.Manager, *apiMock.Mock, notifications.Notifications) {
	logger.Init(5)
	slog.Info("---- Starting test ----", "test", t.Name())

	configPath := fmt.Sprintf("./%s/config.yaml", conf.Id)
	callsPath := fmt.Sprintf("./%s/calls.json", conf.Id)
	wantPath := fmt.Sprintf("./%s/want.json", conf.Id)
	cachePath := fmt.Sprintf("./%s/cache.json", conf.Id)

	c, err := configPkg.New(configPath)
	if err != nil {
		t.Fatal(err)
	}
	c.Data.Cache.Path = cachePath

	m := manager.New(c.Data)

	// TODO: move those into config so it can be set by default as well as via
	// CLI
	m.ForceStrategy = conf.ForceStrategy
	m.RefreshStrategy = conf.RefreshStrategy

	calls, err := apiMock.LoadCallsFromFile(callsPath)
	caller := &apiMock.Mock{Calls: calls}
	m.SetCaller(caller)

	if err = m.Load(); err != nil {
		t.Fatal(err)
	}

	for _, n := range m.Notifications {
		slog.Info("Loaded notification", "id", n.Id)
	}

	if err = m.Refresh(); err != nil {
		t.Fatal(err)
	}

	for _, n := range m.Notifications {
		slog.Info("Refresh notification", "id", n.Id)
	}

	want := notifications.Notifications{}

	raw, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}

	return m, caller, want
}

func TestIntegration(t *testing.T) {
	dirs, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		t.Run(dir.Name(), func(t *testing.T) {
			m, c, want := setup(t, config{
				Id:              dir.Name(),
				RefreshStrategy: manager.ForceRefresh,
			})

			got := m.Notifications

			if !want.Equal(got) {
				t.Fatalf("mismatch notifications\nwant:\n%s\n\ngot:\n%s", want.Debug(), got.Debug())
			}

			if err := c.Done(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
