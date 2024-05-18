package gh

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Client struct {
	client *api.RESTClient
	cache  cache.ExpiringReadWriter
}

func NewClient(cache cache.ExpiringReadWriter) (*Client, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		cache:  cache,
	}, err
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
		var pulledNotifications []notifications.Notification
		if err := c.client.Get("notifications", &pulledNotifications); err != nil {
			return nil, err
		}

		allNotifications.Append(pulledNotifications)

		// This will favor the cached notifications as they are first into the
		// slice.
		// allNotifications = allNotifications.Uniq()

		if err := c.cache.Write(allNotifications); err != nil {
			fmt.Printf("Error while writing the cache: %#v", err)
		}
	}

	return allNotifications, nil
}
