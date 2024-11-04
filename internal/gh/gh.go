// Interact with GitHub's api, wrapper around cli/go-gh client object.
package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/nobe4/gh-not/internal/api"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/notifications"
)

var (
	linkRE = regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)

	DefaultUrl = url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/notifications",
	}
)

type Client struct {
	API      api.Requestor
	cache    cache.RefreshReadWriter
	maxRetry int
	maxPage  int
	url      string
}

func NewClient(api api.Requestor, cache cache.RefreshReadWriter, config config.Endpoint) *Client {
	url := DefaultUrl

	query := url.Query()
	if config.All {
		query.Set("all", "true")
	}
	if config.PerPage > 0 && config.PerPage != 100 {
		query.Set("per_page", fmt.Sprintf("%d", config.PerPage))
	}
	url.RawQuery = query.Encode()

	return &Client{
		API:      api,
		cache:    cache,
		maxRetry: config.MaxRetry,
		maxPage:  config.MaxPage,
		url:      url.String(),
	}
}

// isRetryable returns true if the error is retryable.
// It is pretty permissive, as the /notifications endpoint is flaky.
// Unexpected status codes and decoding errors are considered retryable.
// See https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#list-notifications-for-the-authenticated-user--status-codes
func isRetryable(e error) bool {
	var httpError *ghapi.HTTPError
	if errors.As(e, &httpError) {
		switch httpError.StatusCode {
		case 404, 502, 504: // expected status code
			return true
		}
	}

	var urlError *url.Error
	if errors.As(e, &urlError) {
		return true
	}

	if errors.Is(e, decodeError) {
		return true
	}

	return false
}

func parse(r *http.Response) ([]*notifications.Notification, string, error) {
	n := []*notifications.Notification{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&n); err != nil {
		slog.Warn("error decoding response", "error", err)

		// Returning a generic error makes it retryable.
		// A valid body can always be decoded, even if it is empty.
		return nil, "", decodeError
	}
	defer r.Body.Close()

	return n, nextPageLink(&r.Header), nil
}

func nextPageLink(h *http.Header) string {
	for _, m := range linkRE.FindAllStringSubmatch(h.Get("Link"), -1) {
		if len(m) > 2 && m[2] == "next" {
			return m[1]
		}
	}
	return ""
}

func (c *Client) request(verb, endpoint string, body io.Reader) ([]*notifications.Notification, string, error) {
	slog.Debug("request", "verb", verb, "endpoint", endpoint)
	response, err := c.API.Request(verb, endpoint, body)
	if err != nil {
		return nil, "", err
	}

	return parse(response)
}

func (c *Client) retry(verb, endpoint string, body io.Reader) ([]*notifications.Notification, string, error) {
	for i := c.maxRetry; i >= 0; i-- {
		notifications, next, err := c.request(verb, endpoint, body)
		if err == nil {
			return notifications, next, nil
		}

		if isRetryable(err) {
			slog.Warn("endpoint failed with retryable error", "error", err, "endpoint", endpoint, "retry left", i)
			continue
		}

		return nil, "", err
	}

	return nil, "", RetryError{verb, endpoint}
}

// inspired by https://github.com/cli/go-gh/blob/25db6b99518c88e03f71dbe9e58397c4cfb62caf/example_gh_test.go#L96-L134
func (c *Client) paginate() (notifications.Notifications, error) {
	var list notifications.Notifications
	var pageList []*notifications.Notification
	var err error

	pageLeft := c.maxPage
	endpoint := c.url

	for endpoint != "" && pageLeft > 0 {
		slog.Info("API REST request", "endpoint", endpoint, "page_left", pageLeft)

		pageList, endpoint, err = c.retry(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}

		list = append(list, pageList...)

		pageLeft--
	}

	return list, nil
}

func (c *Client) Notifications() (notifications.Notifications, error) {
	return c.paginate()
}
