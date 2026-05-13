package manager

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nobe4/gh-not/internal/api/mock"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/notifications"
)

var errMockEnrich = errors.New("mock enrich failure")

func TestEnrich(t *testing.T) {
	t.Parallel()

	mockResponse := func(body string) *http.Response {
		t.Helper()

		return &http.Response{
			Body: io.NopCloser(strings.NewReader(body)),
		}
	}

	testManager := func(t *testing.T, calls []mock.Call) (*Manager, *mock.Mock) {
		t.Helper()

		m := New(&config.Data{
			Cache: config.Cache{Path: filepath.Join(t.TempDir(), "cache.json")},
		})
		caller := &mock.Mock{Calls: calls}
		m.SetCaller(caller)

		return m, caller
	}

	t.Run("continues after per-notification failure", func(t *testing.T) {
		t.Parallel()

		m, caller := testManager(t, []mock.Call{
			{URL: "bad-url", Error: errMockEnrich},
			{URL: "good-url", Response: mockResponse(`{"state":"open"}`)},
		})
		ns := notifications.Notifications{
			&notifications.Notification{ID: "bad", Subject: notifications.Subject{URL: "bad-url"}},
			&notifications.Notification{ID: "good", Subject: notifications.Subject{URL: "good-url"}},
		}

		m.Enrich(ns)

		if ns[0].Meta.Enriched {
			t.Fatal("expected failed notification to remain unenriched")
		}

		if !ns[1].Meta.Enriched {
			t.Fatal("expected later notification to be enriched")
		}

		if err := caller.Done(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("force enriches cached notification", func(t *testing.T) {
		t.Parallel()

		m, caller := testManager(t, []mock.Call{
			{URL: "subject-url", Response: mockResponse(`{"state":"closed"}`)},
		})
		m.ForceStrategy = ForceEnrich
		ns := notifications.Notifications{
			&notifications.Notification{
				ID:      "0",
				Subject: notifications.Subject{URL: "subject-url", State: "open"},
				Meta:    notifications.Meta{Enriched: true},
			},
		}

		m.Enrich(ns)

		if !ns[0].Meta.Enriched {
			t.Fatal("expected notification to remain enriched")
		}

		if ns[0].Subject.State != "closed" {
			t.Fatalf("expected state to be closed but got %q", ns[0].Subject.State)
		}

		if err := caller.Done(); err != nil {
			t.Fatal(err)
		}
	})
}
