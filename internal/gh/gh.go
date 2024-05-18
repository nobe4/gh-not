package gh

import (
	"encoding/json"
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

func (c *Client) loadCache() ([]notifications.Notification, bool, error) {
	expired, err := c.cache.Expired()
	if err != nil {
		return nil, false, err
	}

	content, err := c.cache.Read()
	if err != nil {
		return nil, expired, err
	}

	notifications := []notifications.Notification{}
	if err := json.Unmarshal(content, &notifications); err != nil {
		return nil, expired, err
	}

	return notifications, expired, nil
}

func (c *Client) writeCache(n []notifications.Notification) error {
	marshalled, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return c.cache.Write(marshalled)
}

func (c *Client) Notifications() ([]notifications.Notification, error) {
	allNotifications := []notifications.Notification{}

	cachedNotifications, expired, err := c.loadCache()
	if err != nil {
		fmt.Printf("Error while reading the cache: %#v\n", err)
	} else {
		allNotifications = append(allNotifications, cachedNotifications...)
	}

	if expired {
		fmt.Printf("Cache expired, pulling from the API...\n")
		var pulledNotifications []notifications.Notification
		if err := c.client.Get("notifications", &pulledNotifications); err != nil {
			return nil, err
		}

		allNotifications = append(allNotifications, pulledNotifications...)

		allNotifications = notifications.Uniq(allNotifications)

		if err := c.writeCache(allNotifications); err != nil {
			fmt.Printf("Error while writing the cache: %#v", err)
			cachedNotifications = []notifications.Notification{}
		}
	}

	return allNotifications, nil
}
