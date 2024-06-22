package gh

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/nobe4/gh-not/internal/api/mock"
	"github.com/nobe4/gh-not/internal/notifications"
)

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "http 502",
			err:  &ghapi.HTTPError{StatusCode: 502},
			want: true,
		},
		{
			name: "http 504",
			err:  &ghapi.HTTPError{StatusCode: 504},
			want: true,
		},
		{
			name: "io.EOF",
			err:  io.EOF,
			want: true,
		},
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "other error",
			err:  errors.New("error"),
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := isRetryable(test.err)
			if got != test.want {
				t.Errorf("expected %v, got %v", test.want, got)
			}
		})
	}
}

func notificationsEqual(a, b []*notifications.Notification) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a != nil && b != nil {
			if a[i].Id != b[i].Id {
				return false
			}
		}
	}

	return true
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		expected []*notifications.Notification
		next     string
		fails    bool
	}{
		{
			name:     "empty body",
			response: &http.Response{Body: io.NopCloser(strings.NewReader(`[]`))},
		},
		{
			name:     "invalid body",
			response: &http.Response{Body: io.NopCloser(strings.NewReader(`]`))},
			fails:    true,
		},
		{
			name:     "single notification",
			response: &http.Response{Body: io.NopCloser(strings.NewReader(`[{"id": "1"}]`))},
			expected: []*notifications.Notification{{Id: "1"}},
		},
		{
			name:     "multiple notifications",
			response: &http.Response{Body: io.NopCloser(strings.NewReader(`[{"id": "1"},{"id": "2"}]`))},
			expected: []*notifications.Notification{{Id: "1"}, {Id: "2"}},
		},
		{
			name:     "next page",
			response: &http.Response{Body: io.NopCloser(strings.NewReader(`[{"id": "1"},{"id": "2"}]`)), Header: http.Header{"Link": []string{`<https://next.page>; rel="next"`}}},
			expected: []*notifications.Notification{{Id: "1"}, {Id: "2"}},
			next:     "https://next.page",
		},
		{
			name:     "next page with no notification",
			response: &http.Response{Body: io.NopCloser(strings.NewReader(`[]`)), Header: http.Header{"Link": []string{`<https://next.page>; rel="next"`}}},
			expected: []*notifications.Notification{},
			next:     "https://next.page",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			notifications, next, err := parse(test.response)

			if test.fails && err == nil {
				t.Errorf("expected test to fails")
			} else if !test.fails && err != nil {
				t.Errorf("expected test to pass, got %v", err)
			}

			if !notificationsEqual(test.expected, notifications) {
				t.Errorf("expected %#v, got %#v", test.expected, notifications)
			}

			if next != test.next {
				t.Errorf("expected %s, got %s", test.next, next)
			}
		})
	}
}

func TestNextPageLink(t *testing.T) {
	tests := []struct {
		name     string
		header   http.Header
		expected string
	}{
		{
			name:     "empty header",
			header:   http.Header{},
			expected: "",
		},
		{
			name:     "no link",
			header:   http.Header{"Link": []string{}},
			expected: "",
		},
		{
			name:     "no next link",
			header:   http.Header{"Link": []string{`<https://prev.page>; rel="prev"`}},
			expected: "",
		},
		{
			name:     "next link",
			header:   http.Header{"Link": []string{`<https://next.page>; rel="next"`}},
			expected: "https://next.page",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := nextPageLink(&test.header)
			if got != test.expected {
				t.Errorf("expected %s, got %s", test.expected, got)
			}
		})
	}
}

func NewMockClient(t *testing.T, r []mock.Response) *Client {
	api, err := mock.New(r)
	if err != nil {
		t.Fatal(err)
	}
	return &Client{API: api}
}

func TestRequest(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		expectedError := errors.New("error")
		client := NewMockClient(t, []mock.Response{{Error: expectedError}})

		_, _, err := client.request("GET", "/notifications", nil)
		if err == nil {
			t.Errorf("expected test to fails")
		}

		if err != expectedError {
			t.Errorf("expected %v, got %v", expectedError, err)
		}
	})

	t.Run("calls parse", func(t *testing.T) {
		response := &http.Response{
			Body:   io.NopCloser(strings.NewReader(`[{"id":"0"}]`)),
			Header: http.Header{"Link": []string{`<https://next.page>; rel="next"`}},
		}
		client := NewMockClient(t, []mock.Response{{Response: response}})

		notifications, next, err := client.request("GET", "/notifications", nil)
		if err != nil {
			t.Errorf("expected test to pass")
		}

		if next != "https://next.page" {
			t.Errorf("expected https://next.page, got %s", next)
		}

		if len(notifications) != 1 {
			t.Errorf("expected 1 notification, got %d", len(notifications))
		}

		if notifications[0].Id != "0" {
			t.Errorf("expected notification id 0, got %s", notifications[0].Id)
		}
	})
}

func TestRetry2(t *testing.T) {
	sampleNotifications := []*notifications.Notification{{Id: "0"}}
	sampleError := errors.New("error")
	sampleResponse := func() *http.Response {
		return &http.Response{Body: io.NopCloser(strings.NewReader(`[{"id":"0"}]`))}
	}

	tests := []struct {
		name          string
		responses     []mock.Response
		maxRetry      int
		notifications []*notifications.Notification
		error         error
	}{
		{
			name:      "no retry, fails with an error",
			responses: []mock.Response{{Error: sampleError}},
			error:     sampleError,
		},
		{
			name:          "no retry, succeeds",
			responses:     []mock.Response{{Response: sampleResponse()}},
			notifications: sampleNotifications,
		},
		{
			name: "retry, fails with an error",
			responses: []mock.Response{
				{Error: &ghapi.HTTPError{StatusCode: 502}},
				{Error: sampleError},
			},
			error:    sampleError,
			maxRetry: 1,
		},
		{
			name: "retry, fails with too many retries",
			responses: []mock.Response{
				{Error: &ghapi.HTTPError{StatusCode: 502}},
				{Error: &ghapi.HTTPError{StatusCode: 502}},
				{Error: &ghapi.HTTPError{StatusCode: 502}},
			},
			error:    RetryError{"GET", "/notifications"},
			maxRetry: 2,
		},
		{
			name: "retry, succeeds",
			responses: []mock.Response{
				{Error: &ghapi.HTTPError{StatusCode: 502}},
				{Error: &ghapi.HTTPError{StatusCode: 502}},
				{Response: sampleResponse()},
			},
			notifications: sampleNotifications,
			maxRetry:      2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewMockClient(t, test.responses)
			client.maxRetry = test.maxRetry

			notifications, _, err := client.retry("GET", "/notifications", nil)

			if !errors.Is(err, test.error) {
				t.Errorf("want %#v, got %#v", test.error, err)
			}

			if !notificationsEqual(notifications, test.notifications) {
				t.Errorf("want %#v, got %#v", test.notifications, notifications)
			}
		})
	}
}
