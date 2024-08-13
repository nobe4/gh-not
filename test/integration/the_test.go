package tests

import (
	"encoding/json"
	"fmt"
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

func setup(t *testing.T, conf config) (*manager.Manager, notifications.Notifications) {
	logger.Init(5)

	configPath := fmt.Sprintf("./%s/config.yaml", conf.Id)
	callsPath := fmt.Sprintf("./%s/calls.json", conf.Id)
	wantPath := fmt.Sprintf("./%s/want.json", conf.Id)

	c, err := configPkg.New(configPath)
	if err != nil {
		t.Fatal(err)
	}

	m := manager.New(c.Data)

	// TODO: move those into config so it can be set by default as well as via
	// CLI
	m.ForceStrategy = conf.ForceStrategy
	m.RefreshStrategy = conf.RefreshStrategy

	calls, err := apiMock.LoadCallsFromFile(callsPath)
	m.SetCaller(apiMock.New(calls))

	if err := m.Load(); err != nil {
		t.Fatal(err)
	}
	if err := m.Refresh(); err != nil {
		t.Fatal(err)
	}

	want := notifications.Notifications{}

	raw, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}

	return m, want
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
			m, want := setup(t, config{
				Id:              dir.Name(),
				RefreshStrategy: manager.ForceRefresh,
			})

			got := m.Notifications

			if !want.Equal(got) {
				t.Fatalf("mismatch notifications\nwant %s\ngot  %s", want.Debug(), got.Debug())
			}
		})
	}
}
