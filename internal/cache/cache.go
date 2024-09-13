package cache

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type ExpiringReadWriter interface {
	Read(any) error
	Write(any) error
	Expired() bool
	RefreshedAt() time.Time
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

func (c *FileCache) Read(out any) error {
	content, err := os.ReadFile(c.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Debug("cache doesn't exist", "path", c.path)
			return nil
		}
		return err
	}

	return json.Unmarshal(content, out)
}

func (c *FileCache) Write(in any) error {
	marshalled, err := json.Marshal(in)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(c.path), 0755); err != nil {
		return err
	}

	return os.WriteFile(c.path, marshalled, 0644)
}

func (c *FileCache) Expired() bool {
	return time.Now().After(c.RefreshedAt().Add(c.ttl))
}

func (c *FileCache) RefreshedAt() time.Time {
	info, err := os.Stat(c.path)

	if err != nil {
		// Returning a valid time.Time that's the 0 epoch allows to not return
		// an error but still process the Expiration date logic correctly.
		slog.Warn("Could not read file info", "file", c.path, "error", err)
		return time.Unix(0, 0)
	}

	return info.ModTime()
}
