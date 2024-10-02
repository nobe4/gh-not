package assign

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"github.com/nobe4/gh-not/internal/api/mock"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

func TestRun(t *testing.T) {
	w := &bytes.Buffer{}

	api := &mock.Mock{}
	client := gh.Client{API: api}

	runner := Runner{Client: &client}

	t.Run("not assignees", func(t *testing.T) {
		n := &notifications.Notification{URL: "http://example.com"}

		if err := runner.Run(n, []string{}, w); err == nil {
			t.Fatal("expected error", err)
		}

		if err := api.Done(); err != nil {
			t.Fatal("unexpected error", err)
		}
	})

	t.Run("not an issue or pull", func(t *testing.T) {
		n := &notifications.Notification{URL: "http://example.com"}

		if err := runner.Run(n, []string{"user"}, w); err != nil {
			t.Fatal("unexpected error", err)
		}

		if err := api.Done(); err != nil {
			t.Fatal("unexpected error", err)
		}
	})

	t.Run("return an API failure", func(t *testing.T) {
		expectedError := errors.New("expected error")

		api.Calls = append(api.Calls, mock.Call{
			Verb:  "POST",
			Url:   "https://api.github.com/repos/owner/repo/issues/123/assignees",
			Data:  `{"assignees":["user"]}`,
			Error: expectedError,
		})
		n := &notifications.Notification{
			Subject: notifications.Subject{
				URL: "https://api.github.com/repos/owner/repo/pulls/123",
			},
		}

		if err := runner.Run(n, []string{"user"}, w); err != expectedError {
			t.Fatalf("expected %#v but got %#v", expectedError, err)
		}

		if err := api.Done(); err != nil {
			t.Fatal("unexpected error", err)
		}
	})

	t.Run("works for an issue", func(t *testing.T) {
		api.Calls = append(api.Calls, mock.Call{
			Verb:     "POST",
			Url:      "https://api.github.com/repos/owner/repo/issues/123/assignees",
			Data:     `{"assignees":["user"]}`,
			Response: &http.Response{StatusCode: http.StatusCreated},
		})
		n := &notifications.Notification{
			Subject: notifications.Subject{
				URL: "https://api.github.com/repos/owner/repo/issues/123",
			},
		}

		if err := runner.Run(n, []string{"user"}, w); err != nil {
			t.Fatal("unexpected error", err)
		}

		if err := api.Done(); err != nil {
			t.Fatal("unexpected error", err)
		}
	})

	t.Run("works for a pull", func(t *testing.T) {
		api.Calls = append(api.Calls, mock.Call{
			Verb:     "POST",
			Url:      "https://api.github.com/repos/owner/repo/issues/123/assignees",
			Data:     `{"assignees":["user"]}`,
			Response: &http.Response{StatusCode: http.StatusCreated},
		})
		n := &notifications.Notification{
			Subject: notifications.Subject{
				URL: "https://api.github.com/repos/owner/repo/pulls/123",
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

func TestIsIssueOrPull(t *testing.T) {
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
