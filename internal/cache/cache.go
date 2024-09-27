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

type CacheWrap struct {
	Data any `json:"data"`
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

	wrap := &CacheWrap{Data: out}

	err = json.Unmarshal(content, wrap)
	if err == nil {
		return nil
	}

	var jsonErr *json.UnmarshalTypeError
	if errors.As(err, &jsonErr) {
		return c.deprecatedRead(content, out)
	}

	return err
}

func (c *FileCache) deprecatedRead(content []byte, out any) error {
	slog.Warn("Cache is in an deprecated format. Attempting to read from the old format.")

	if err := json.Unmarshal(content, out); err != nil {
		return err
	}

	// TODO: here infer the other information in the CacheWrap
	return nil
}

func (c *FileCache) Write(in any) error {
	wrap := &CacheWrap{Data: in}

	marshalled, err := json.Marshal(wrap)
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
