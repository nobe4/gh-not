package gh

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/nobe4/gh-not/internal/api/mock"
	"github.com/nobe4/gh-not/internal/notifications"
)

func TestEnrichLeavesNotificationUnchangedOnFailure(t *testing.T) {
	t.Parallel()

	threadExtraBody := `{
		"user": { "login": "new-author", "type": "User" },
		"state": "open",
		"html_url": "https://example.com/issues/1",
		"assignees": [{ "login": "assignee1", "type": "User" }],
		"requested_reviewers": [{ "login": "reviewer1", "type": "User" }],
		"requested_teams": [{ "name": "team1", "id": 42 }],
		"merged_by": { "login": "merger", "type": "User" }
	}`

	oldUser := func(login string) notifications.User {
		return notifications.User{Login: login, Type: "User"}
	}

	// seed builds a notification that already has enriched fields, as it could
	// after a prior successful Enrich pass.
	seed := func() *notifications.Notification {
		n := mockNotification(0)
		n.Author = oldUser("old-author")
		n.Subject.State = "closed"
		n.Subject.HTMLURL = "https://example.com/old"
		n.Assignees = []notifications.User{oldUser("old-assignee")}
		n.Reviewers = []notifications.User{oldUser("old-reviewer")}
		n.ReviewersTeams = []notifications.Team{{Name: "old-team", ID: 7}}
		n.MergedBy = oldUser("old-merger")
		n.LatestCommentor = oldUser("old-commentor")

		return n
	}

	tests := []struct {
		name  string
		calls []mock.Call
	}{
		{
			name:  "thread extra fails",
			calls: []mock.Call{{URL: mockSubjectURL(0), Error: errSample}},
		},
		{
			name: "last commentor fails",
			calls: []mock.Call{
				{URL: mockSubjectURL(0), Response: &http.Response{Body: io.NopCloser(strings.NewReader(threadExtraBody))}},
				{URL: mockLatestCommentURL(0), Error: errSample},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			want := seed()
			n := seed()
			client, m := mockClient(test.calls)

			if err := client.Enrich(n); !errors.Is(err, errSample) {
				t.Fatalf("expected %#v, got %#v", errSample, err)
			}

			if !reflect.DeepEqual(n, want) {
				t.Errorf("notification mutated despite error\nwant: %#v\ngot:  %#v", want, n)
			}

			if err := m.Done(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
