// Interact with GitHub's api, wrapper around cli/go-gh client object.
package gh

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/nobe4/gh-not/internal/api"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/notifications"
)

const (
	DefaultEndpoint = "https://api.github.com/notifications"
)

var (
	linkRE = regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)
)

type Client struct {
	API      api.Caller
	cache    cache.ExpiringReadWriter
	maxRetry int
	maxPage  int
	endpoint string
}

func NewClient(api api.Caller, cache cache.ExpiringReadWriter, config config.Endpoint) *Client {
	endpoint := DefaultEndpoint
	if config.All {
		endpoint += "?all=true"
	}

	return &Client{
		API:      api,
		cache:    cache,
		maxRetry: config.MaxRetry,
		maxPage:  config.MaxPage,
		endpoint: endpoint,
	}
}

func isRetryable(e error) bool {
	var httpError *ghapi.HTTPError
	if errors.As(e, &httpError) {
		switch httpError.StatusCode {
		case 502, 504:
			return true
		}
	}

	if errors.Is(e, io.EOF) {
		return true
	}

	return false
}

func parse(r *http.Response) ([]*notifications.Notification, string, error) {
	n := []*notifications.Notification{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&n); err != nil {
		return nil, "", err
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
			slog.Warn("endpoint failed with retryable error", "endpoint", endpoint, "retry left", i)
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
	endpoint := c.endpoint

	for endpoint != "" && pageLeft >= 0 {
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
	list, err := c.paginate()
	if err != nil {
		return nil, err
	}

	for i, n := range list {
		if err := c.enrichNotification(n); err != nil {
			return nil, err
		}

		list[i] = n
	}

	return list, nil
}
