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
	SetRefreshedAt(time.Time)
}

type FileCache struct {
	path string
	ttl  time.Duration
	wrap *CacheWrap
}

type CacheWrap struct {
	Data        any       `json:"data"`
	RefreshedAt time.Time `json:"refreshed_at"`
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

	c.wrap = &CacheWrap{Data: out}

	err = json.Unmarshal(content, c.wrap)
	if err == nil {
		return nil
	}

	var jsonErr *json.UnmarshalTypeError
	if errors.As(err, &jsonErr) {
		return c.deprecatedRead(content)
	}

	return err
}

func (c *FileCache) RefreshedAt() time.Time {
	return c.wrap.RefreshedAt
}

func (c *FileCache) SetRefreshedAt(t time.Time) {
	c.wrap.RefreshedAt = t
}

func (c *FileCache) Expired() bool {
	return time.Now().After(c.RefreshedAt().Add(c.ttl))
}

func (c *FileCache) deprecatedRead(content []byte) error {
	slog.Warn("Cache is in an deprecated format. Attempting to read from the old format.")

	if err := json.Unmarshal(content, c.wrap.Data); err != nil {
		return err
	}

	c.wrap.RefreshedAt = time.Unix(0, 0)

	return nil
}

func (c *FileCache) Write(in any) error {
	c.wrap.Data = in

	marshalled, err := json.Marshal(c.wrap)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(c.path), 0755); err != nil {
		return err
	}

	return os.WriteFile(c.path, marshalled, 0644)
}
