package integration

import (
	"encoding/json"
	"os"
	"testing"

	apiMock "github.com/nobe4/gh-not/internal/api/mock"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/manager"
	"github.com/nobe4/gh-not/internal/notifications"
)

type MockManagerConfig struct {
	ConfigPath      string
	CallsPath       string
	ForceStrategy   manager.ForceStrategy
	RefreshStrategy manager.RefreshStrategy
}

func MockManager(t *testing.T, conf MockManagerConfig) *manager.Manager {
	c, err := config.New(conf.ConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	m := manager.New(c.Data)

	m.ForceStrategy = conf.ForceStrategy
	m.RefreshStrategy = conf.RefreshStrategy

	calls, err := apiMock.LoadCallsFromFile(conf.CallsPath)
	m.SetCaller(apiMock.New(calls))

	return m
}

func Compare(t *testing.T, expectedPath string, got notifications.Notifications) {
	expected := notifications.Notifications{}
	raw, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &expected); err != nil {
		t.Fatal(err)
	}

	if len(got) != len(expected) {
		t.Fatalf("expected %d notifications, got %d", len(expected), len(got))
	}

	for i := range got {
		g := got[i]
		e := expected[i]

		if g.Id != e.Id {
			t.Fatalf("Id mismatch: expected %s, got %s", e.Id, g.Id)
		}

		if g.Subject.State != e.Subject.State {
			t.Fatalf("Subject.State mismatch: expected %s, got %s", e.Subject.State, g.Subject.State)
		}
	}
}
