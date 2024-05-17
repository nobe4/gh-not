package gh

import (
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Client struct {
	*api.RESTClient
}

func NewClient() (*Client, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}

	return &Client{client}, err
}

func (c *Client) Notifications() ([]notifications.Notification, error) {
	var allNotifications []notifications.Notification

	if err := c.Get("notifications", &allNotifications); err != nil {
		return nil, err
	}

	return allNotifications, nil
}
