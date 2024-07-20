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
	info, err := os.Stat(c.path)

	if err != nil {
		// We could return the error here, but we will get it again when we try
		// to read/write the cache, so we can ignore it here.
		return true
	}

	expiration := info.ModTime().Add(c.ttl)
	expired := time.Now().After(expiration)

	return expired
}
