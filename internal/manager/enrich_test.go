package manager

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nobe4/gh-not/internal/api/mock"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

func TestEnrichWorkers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *config.Data
		want   int
	}{
		{"nil config", nil, 1},
		{"zero workers", &config.Data{}, 1},
		{"negative workers", &config.Data{Enrichment: config.Enrichment{Workers: -1}}, 1},
		{"set workers", &config.Data{Enrichment: config.Enrichment{Workers: 5}}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := &Manager{config: tt.config}
			if got := m.enrichWorkers(); got != tt.want {
				t.Fatalf("want %d, got %d", tt.want, got)
			}
		})
	}
}

func TestShouldEnrich(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		n     *notifications.Notification
		force bool
		want  bool
	}{
		{"nil notification", nil, false, false},
		{"default", &notifications.Notification{}, false, true},
		{"already enriched", &notifications.Notification{Meta: notifications.Meta{Enriched: true}}, false, false},
		{"done", &notifications.Notification{Meta: notifications.Meta{Done: true}}, false, false},
		{"done and enriched", &notifications.Notification{Meta: notifications.Meta{Done: true, Enriched: true}}, false, false},
		{"force on enriched", &notifications.Notification{Meta: notifications.Meta{Enriched: true}}, true, true},
		{"force on done", &notifications.Notification{Meta: notifications.Meta{Done: true}}, true, true},
		{"force on done and enriched", &notifications.Notification{Meta: notifications.Meta{Done: true, Enriched: true}}, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := &Manager{}
			if tt.force {
				m.ForceStrategy = ForceEnrich
			}

			if got := m.shouldEnrich(tt.n); got != tt.want {
				t.Fatalf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestEnrich(t *testing.T) {
	t.Parallel()

	mockCall := func(t *testing.T, url string, body any) mock.Call {
		t.Helper()

		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal mock body: %v", err)
		}

		return mock.Call{
			URL:      url,
			Response: &http.Response{Body: io.NopCloser(strings.NewReader(string(b)))},
		}
	}

	subjectCall := func(id string) mock.Call {
		return mockCall(t, "https://subject.url/"+id, struct {
			User    notifications.User `json:"user"`
			HTMLURL string             `json:"html_url"`
		}{
			User:    notifications.User{Login: "author", Type: "User"},
			HTMLURL: "https://html.url/",
		})
	}

	commentCall := func(id string) mock.Call {
		return mockCall(t, "https://latest.comment.url/"+id, struct {
			User notifications.User `json:"user"`
		}{
			User: notifications.User{Login: "commentor", Type: "User"},
		})
	}

	testNotification := func(id string) *notifications.Notification {
		return &notifications.Notification{
			ID: id,
			Subject: notifications.Subject{
				URL:              "https://subject.url/" + id,
				LatestCommentURL: "https://latest.comment.url/" + id,
			},
		}
	}

	tests := []struct {
		name  string
		calls []mock.Call
		ns    notifications.Notifications
		want  []bool
	}{
		{
			name: "enriches all notifications",
			calls: []mock.Call{
				subjectCall("1"),
				commentCall("1"),
				subjectCall("2"),
				commentCall("2"),
			},
			ns:   notifications.Notifications{testNotification("1"), testNotification("2")},
			want: []bool{true, true},
		},
		{
			name: "continues after failure",
			calls: []mock.Call{
				subjectCall("2"),
				commentCall("2"),
			},
			ns: func() notifications.Notifications {
				failed := testNotification("1")
				failed.Subject.URL = "https://unexpected.url/1"

				return notifications.Notifications{failed, testNotification("2")}
			}(),
			want: []bool{false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			requestor := &mock.Mock{Calls: tt.calls}

			m := &Manager{
				client: gh.NewClient(requestor, nil, gh.Endpoint{}),
				config: &config.Data{Enrichment: config.Enrichment{Workers: 1}},
			}

			m.Enrich(tt.ns)

			for i, n := range tt.ns {
				if n.Meta.Enriched != tt.want[i] {
					t.Fatalf("notification %s: want enriched=%v, got %v", n.ID, tt.want[i], n.Meta.Enriched)
				}
			}

			if err := requestor.Done(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
