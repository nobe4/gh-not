package gh

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/notifications"
)

const (
	path = "notifications"
)

type APICaller interface {
	Do(string, string, io.Reader, interface{}) error
}

type Client struct {
	API   APICaller
	cache cache.ExpiringReadWriter
}

func NewClient(api APICaller, cache cache.ExpiringReadWriter) *Client {
	return &Client{
		API:   api,
		cache: cache,
	}
}

func (c *Client) loadCache() (notifications.NotificationMap, bool, error) {
	expired, err := c.cache.Expired()
	if err != nil {
		return nil, false, err
	}

	n, err := c.cache.Read()
	if err != nil {
		return nil, expired, err
	}

	return n, expired, nil
}

func (c *Client) pullNotificationFromApi() ([]notifications.Notification, error) {
	var list []notifications.Notification

	slog.Debug("API REST request", "path", path)
	if err := c.API.Do(http.MethodGet, path, nil, &list); err != nil {
		return nil, err
	}

	for i, n := range list {
		enriched, err := c.enrichNotification(n)
		if err != nil {
			return nil, err
		}

		list[i] = enriched
	}

	return list, nil
}

func (c *Client) Notifications() (notifications.NotificationMap, error) {
	allNotifications := make(notifications.NotificationMap)

	cachedNotifications, expired, err := c.loadCache()
	if err != nil {
		fmt.Printf("Error while reading the cache: %#v\n", err)
	} else if cachedNotifications != nil {
		allNotifications = cachedNotifications
	}

	if expired {
		fmt.Printf("Cache expired, pulling from the API...\n")
		pulledNotifications, err := c.pullNotificationFromApi()
		if err != nil {
			return nil, err
		}

		allNotifications.Append(pulledNotifications)

		if err := c.cache.Write(allNotifications); err != nil {
			fmt.Printf("Error while writing the cache: %#v", err)
		}
	}

	return allNotifications, nil
}
