package gh

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/nobe4/gh-not/internal/api/mock"
	"github.com/nobe4/gh-not/internal/notifications"
)

const (
	verb     = "GET"
	endpoint = "/notifications"
)

var (
	retriableError = &api.HTTPError{StatusCode: 502}
	sampleError    = errors.New("error")
	retryError     = RetryError{verb, endpoint}
)

func mockSubjectUrl(id int) string {
	return "https://subject.url/" + strconv.Itoa(id)
}

func mockNotification(id int) *notifications.Notification {
	return &notifications.Notification{
		Id: strconv.Itoa(id),
		Subject: notifications.Subject{
			URL: mockSubjectUrl(id),
		},
	}
}

func mockNotifications(ids []int) []*notifications.Notification {
	n := []*notifications.Notification{}
	for _, id := range ids {
		n = append(n, mockNotification(id))
	}
	return n
}

func mockNotificationsResponse(t *testing.T, ids []int, next bool) *http.Response {
	n := mockNotifications(ids)
	body, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}

	link := ""
	if next {
		link = `<https://next.page>; rel="next"`
	}

	return &http.Response{
		Body:   io.NopCloser(bytes.NewReader(body)),
		Header: http.Header{"Link": []string{link}},
	}
}

func notificationsEqual(a, b []*notifications.Notification) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] == nil && b[i] == nil {
			continue
		}

		if a[i] == nil || b[i] == nil {
			return false
		}

		if a[i].Id != b[i].Id {
			return false
		}
	}

	return true
}

func mockClient(c []mock.Call) (*Client, *mock.Mock) {
	mock := &mock.Mock{Calls: c}
	return &Client{
		API:      mock,
		endpoint: endpoint,
		maxRetry: 100,
		maxPage:  100,
	}, mock
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "http 404",
			err:  &api.HTTPError{StatusCode: 404},
			want: true,
		},
		{
			name: "http 502",
			err:  retriableError,
			want: true,
		},
		{
			name: "http 504",
			err:  &api.HTTPError{StatusCode: 504},
			want: true,
		},
		{
			name: "DecodeError",
			err:  decodeError,
			want: true,
		},
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "other error",
			err:  sampleError,
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := isRetryable(test.err)
			if got != test.want {
				t.Errorf("expected %#v, got %#v", test.want, got)
			}
		})
	}
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
			name:     "invalid body",
			response: &http.Response{Body: io.NopCloser(strings.NewReader(`{"a":1]`))},
			fails:    true,
		},
		{
			name:     "single notification",
			response: mockNotificationsResponse(t, []int{0}, false),
			expected: mockNotifications([]int{0}),
		},
		{
			name:     "multiple notifications",
			response: mockNotificationsResponse(t, []int{0, 1}, false),
			expected: mockNotifications([]int{0, 1}),
		},
		{
			name:     "next page",
			response: mockNotificationsResponse(t, []int{0, 1}, true),
			expected: mockNotifications([]int{0, 1}),
			next:     "https://next.page",
		},
		{
			name:     "next page with no notification",
			response: mockNotificationsResponse(t, []int{}, true),
			next:     "https://next.page",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			notifications, next, err := parse(test.response)

			if test.fails && err == nil {
				t.Errorf("expected test to fails")
			} else if !test.fails && err != nil {
				t.Errorf("expected test to pass, got %#v", err)
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
			name:   "empty header",
			header: http.Header{},
		},
		{
			name:   "no link",
			header: http.Header{"Link": []string{}},
		},
		{
			name:   "no next link",
			header: http.Header{"Link": []string{`<https://prev.page>; rel="prev"`}},
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

func TestRequest(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		expectedError := errors.New("error")
		client, api := mockClient([]mock.Call{{Error: expectedError}})

		_, _, err := client.request(verb, endpoint, nil)
		if err == nil {
			t.Errorf("expected test to fails")
		}

		if err != expectedError {
			t.Errorf("expected %#v, got %#v", expectedError, err)
		}

		if err := api.Done(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("calls parse", func(t *testing.T) {
		response := &http.Response{
			Body:   io.NopCloser(strings.NewReader(`[{"id":"0"}]`)),
			Header: http.Header{"Link": []string{`<https://next.page>; rel="next"`}},
		}
		client, api := mockClient([]mock.Call{{Response: response}})

		notifications, next, err := client.request(verb, endpoint, nil)
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

		if err := api.Done(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestRetry(t *testing.T) {
	tests := []struct {
		name          string
		calls         []mock.Call
		maxRetry      int
		notifications []*notifications.Notification
		error         error
	}{
		{
			name:  "no retry, fails with an error",
			calls: []mock.Call{{Error: sampleError}},
			error: sampleError,
		},
		{
			name: "no retry, succeeds",
			calls: []mock.Call{
				{Response: mockNotificationsResponse(t, []int{0}, false)},
			},
			notifications: mockNotifications([]int{0}),
		},
		{
			name: "retry, fails with an error",
			calls: []mock.Call{
				{Error: retriableError},
				{Error: sampleError},
			},
			error:    sampleError,
			maxRetry: 1,
		},
		{
			name: "retry, fails with too many retries",
			calls: []mock.Call{
				{Error: retriableError},
				{Error: retriableError},
				{Error: retriableError},
			},
			error:    retryError,
			maxRetry: 2,
		},
		{
			name: "retry, succeeds",
			calls: []mock.Call{
				{Error: retriableError},
				{Error: retriableError},
				{Response: mockNotificationsResponse(t, []int{0}, false)},
			},
			notifications: mockNotifications([]int{0}),
			maxRetry:      2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, api := mockClient(test.calls)
			client.maxRetry = test.maxRetry

			notifications, _, err := client.retry(verb, endpoint, nil)

			if !errors.Is(err, test.error) {
				t.Errorf("want %#v, got %#v", test.error, err)
			}

			if !notificationsEqual(notifications, test.notifications) {
				t.Errorf("want %#v, got %#v", test.notifications, notifications)
			}
			if err := api.Done(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPaginate(t *testing.T) {
	tests := []struct {
		name          string
		calls         []mock.Call
		maxRetry      int
		maxPage       int
		notifications []*notifications.Notification
		error         error
	}{
		{
			name: "zero page",
		},
		{
			name:    "one page, fails with an error",
			maxPage: 1,
			calls:   []mock.Call{{Error: sampleError}},
			error:   sampleError,
		},
		{
			name:     "one page, retries and fails with an error",
			maxRetry: 1,
			maxPage:  1,
			calls: []mock.Call{
				{Error: retriableError},
				{Error: sampleError},
			},
			error: sampleError,
		},
		{
			name:     "one page, retries to many times and fails",
			maxRetry: 1,
			maxPage:  1,
			calls: []mock.Call{
				{Error: retriableError},
				{Error: retriableError},
			},
			error: retryError,
		},
		{
			name:    "one page, succeeds",
			maxPage: 1,
			calls: []mock.Call{
				{Response: mockNotificationsResponse(t, []int{0, 1}, false)},
			},
			notifications: mockNotifications([]int{0, 1}),
		},
		{
			name:     "one page, retries and succeeds",
			maxRetry: 1,
			maxPage:  1,
			calls: []mock.Call{
				{Error: retriableError},
				{Response: mockNotificationsResponse(t, []int{0, 1}, false)},
			},
			notifications: mockNotifications([]int{0, 1}),
		},
		{
			name:    "two pages available, fetch only one",
			maxPage: 1,
			calls: []mock.Call{
				{Response: mockNotificationsResponse(t, []int{0}, true)},
			},
			notifications: mockNotifications([]int{0}),
		},
		{
			name:    "two pages, fails with an error on the second page",
			maxPage: 2,
			calls: []mock.Call{
				{Response: mockNotificationsResponse(t, []int{0}, true)},
				{Error: sampleError},
			},
			error: sampleError,
		},
		{
			name:    "two pages, succeeds",
			maxPage: 2,
			calls: []mock.Call{
				{Response: mockNotificationsResponse(t, []int{0}, true)},
				{Response: mockNotificationsResponse(t, []int{1}, true)},
			},
			notifications: mockNotifications([]int{0, 1}),
		},
		{
			name:    "three pages, but only two are requested",
			maxPage: 2,
			calls: []mock.Call{
				{Response: mockNotificationsResponse(t, []int{0}, true)},
				{Response: mockNotificationsResponse(t, []int{1}, true)},
			},
			notifications: mockNotifications([]int{0, 1}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, api := mockClient(test.calls)
			client.maxRetry = test.maxRetry
			client.maxPage = test.maxPage

			notifications, err := client.paginate()

			if !errors.Is(err, test.error) {
				t.Errorf("want %#v, got %#v", test.error, err)
			}

			if !notificationsEqual(notifications, test.notifications) {
				t.Errorf("want %#v, got %#v", test.notifications, notifications)
			}
			if err := api.Done(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestEnrich(t *testing.T) {
	tests := []struct {
		name         string
		calls        []mock.Call
		notification *notifications.Notification
		assertError  func(*testing.T, error)
	}{

		{
			name: "no notification",
		},
		{
			name: "one notification",
			calls: []mock.Call{
				{
					Endpoint: mockSubjectUrl(0),
					Response: &http.Response{
						Body: io.NopCloser(strings.NewReader(`{"author":{"login":"author"},"subject":{"title":"subject"}}`)),
					},
				},
			},
			notification: mockNotification(0),
		},
		{
			name:         "fail to enrich",
			calls:        []mock.Call{{Error: sampleError}},
			notification: mockNotification(0),
			assertError: func(t *testing.T, err error) {
				if err != sampleError {
					t.Errorf("expected to fail with %#v but got %#v", sampleError, err)
				}
			},
		},
		{
			name: "fail to unmarshal",
			calls: []mock.Call{
				{
					Endpoint: mockSubjectUrl(0),
					Response: &http.Response{
						Body: io.NopCloser(strings.NewReader(`not json`)),
					},
				},
			},
			notification: mockNotification(0),
			assertError: func(t *testing.T, err error) {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}

				expected := &json.SyntaxError{}
				if !errors.As(err, &expected) {
					t.Errorf("expected to fail with %#v but got %#v", expected, err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, api := mockClient(test.calls)

			err := client.Enrich(test.notification)

			// TODO: make this test check for the author/subject
			if test.assertError != nil {
				test.assertError(t, err)
			}
			if err := api.Done(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
