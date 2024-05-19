package cache

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/nobe4/gh-not/internal/notifications"
)

type ExpiringReadWriter interface {
	Read() (notifications.NotificationMap, error)
	Write(notifications.NotificationMap) error
	Expired() (bool, error)
}

type FileCache struct {
	path string
	ttl  time.Duration
}

func NewFileCache(ttlInHours int, path string) *FileCache {
	return &FileCache{
		path: path,
		ttl:  time.Duration(ttlInHours) * time.Hour,
	}
}

func (c *FileCache) Read() (notifications.NotificationMap, error) {
	content, err := os.ReadFile(c.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	notifications := notifications.NotificationMap{}
	if err := json.Unmarshal(content, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (c *FileCache) Write(n notifications.NotificationMap) error {
	// In case we have items without IDs, we can safely delete it.
	delete(n, "")

	marshalled, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, marshalled, 0644)
}

func (c *FileCache) Expired() (bool, error) {
	info, err := os.Stat(c.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}
		return false, err
	}

	expiration := info.ModTime().Add(c.ttl)
	expired := time.Now().After(expiration)

	return expired, nil
}
