package assign

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nobe4/gh-not/internal/api/mock"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("not assignees", func(t *testing.T) {
		t.Parallel()

		w := &bytes.Buffer{}
		api := &mock.Mock{}
		client := gh.Client{API: api}
		runner := Runner{Client: &client}

		n := &notifications.Notification{URL: "http://example.com"}

		if err := runner.Run(n, []string{}, w); err == nil {
			t.Fatal("expected error", err)
		}

		if err := api.Done(); err != nil {
			t.Fatal("unexpected error", err)
		}
	})

	t.Run("not an issue or pull", func(t *testing.T) {
		t.Parallel()

		w := &bytes.Buffer{}
		api := &mock.Mock{}
		client := gh.Client{API: api}
		runner := Runner{Client: &client}
		n := &notifications.Notification{URL: "http://example.com"}

		if err := runner.Run(n, []string{"user"}, w); err != nil {
			t.Fatal("unexpected error", err)
		}

		if err := api.Done(); err != nil {
			t.Fatal("unexpected error", err)
		}
	})

	t.Run("return an API failure", func(t *testing.T) {
		t.Parallel()

		w := &bytes.Buffer{}
		api := &mock.Mock{}
		client := gh.Client{API: api}
		runner := Runner{Client: &client}
		expectedError := errors.New("expected error")

		api.Calls = append(api.Calls, mock.Call{
			Verb:  "POST",
			URL:   "https://api.github.com/repos/owner/repo/issues/123/assignees",
			Data:  `{"assignees":["user"]}`,
			Error: expectedError,
		})
		n := &notifications.Notification{
			Subject: notifications.Subject{
				URL: "https://api.github.com/repos/owner/repo/pulls/123",
			},
		}

		if err := runner.Run(n, []string{"user"}, w); !errors.Is(err, expectedError) {
			t.Fatalf("expected %#v but got %#v", expectedError, err)
		}

		if err := api.Done(); err != nil {
			t.Fatal("unexpected error", err)
		}
	})

	for _, c := range []string{"issues", "pulls"} {
		t.Run("works for "+c, func(t *testing.T) {
			t.Parallel()

			w := &bytes.Buffer{}
			api := &mock.Mock{}
			client := gh.Client{API: api}
			runner := Runner{Client: &client}

			api.Calls = append(api.Calls, mock.Call{
				Verb:     "POST",
				URL:      "https://api.github.com/repos/owner/repo/issues/123/assignees",
				Data:     `{"assignees":["user"]}`,
				Response: &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(""))},
			})
			n := &notifications.Notification{
				Subject: notifications.Subject{
					URL: "https://api.github.com/repos/owner/repo/" + c + "/123",
				},
			}

			if err := runner.Run(n, []string{"user"}, w); err != nil {
				t.Fatal("unexpected error", err)
			}

			if err := api.Done(); err != nil {
				t.Fatal("unexpected error", err)
			}
		})
	}
}

func TestIsIssueOrPull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		url   string
		want  string
		match bool
	}{
		{
			url:   "http://example.com",
			match: false,
		},
		{
			url:   "https://github.com",
			match: false,
		},
		{
			url:   "https://api.github.com",
			match: false,
		},
		{
			url:   "https://api.github.com/repos/owner/repo",
			match: false,
		},
		{
			url:   "https://api.github.com/repos/owner/repo/pulls",
			match: false,
		},
		{
			url:   "https://api.github.com/repos/owner/repo/issues",
			match: false,
		},
		{
			url:   "https://api.github.com/repos/owner/repo/pulls/123",
			want:  "https://api.github.com/repos/owner/repo/issues/123",
			match: true,
		},
		{
			url:   "https://api.github.com/repos/owner/repo/issues/123",
			want:  "https://api.github.com/repos/owner/repo/issues/123",
			match: true,
		},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			t.Parallel()

			got, match := issueURL(test.url)
			if match != test.match {
				t.Errorf("want %v but got %v", test.match, match)
			}

			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}
}
