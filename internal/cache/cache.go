package cache

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/nobe4/gh-not/internal/notifications"
)

type ExpiringReadWriter interface {
	Read() (notifications.Notifications, error)
	Write(notifications.Notifications) error
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

func (c *FileCache) Read() (notifications.Notifications, error) {
	content, err := os.ReadFile(c.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	notifications := notifications.Notifications{}
	if err := json.Unmarshal(content, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (c *FileCache) Write(n notifications.Notifications) error {
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
