/*
Package cache provides a simple file-based cache implementation that fulfills
the RefreshReadWriter interface.

It writes and reads the cache to a file in JSON format, along with the last
refresh time.
*/
package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type RefreshReadWriter interface {
	Read(d any) error
	Write(d any) error
	Refresh(t time.Time)
	RefreshedAt() time.Time
}

type FileCache struct {
	path string
	wrap *CacheWrap
}

type CacheWrap struct {
	Data        any       `json:"data"`
	RefreshedAt time.Time `json:"refreshed_at"`
}

func NewFileCache(path string) *FileCache {
	return &FileCache{
		path: path,
		wrap: &CacheWrap{RefreshedAt: time.Unix(0, 0)},
	}
}

func (c *FileCache) Read(out any) error {
	slog.Debug("Reading cache", "path", c.path)

	content, err := os.ReadFile(c.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Debug("cache doesn't exist", "path", c.path)
			return nil
		}

		return fmt.Errorf("failed to read cache: %w", err)
	}

	c.wrap.Data = out

	err = json.Unmarshal(content, c.wrap)
	if err == nil {
		return nil
	}

	var jsonErr *json.UnmarshalTypeError
	if errors.As(err, &jsonErr) {
		return c.deprecatedRead(content)
	}

	return fmt.Errorf("failed to unmarshal cache: %w", err)
}

func (c *FileCache) Refresh(t time.Time) {
	c.wrap.RefreshedAt = t
}

func (c *FileCache) RefreshedAt() time.Time {
	return c.wrap.RefreshedAt
}

func (c *FileCache) deprecatedRead(content []byte) error {
	slog.Warn("Cache is in an format deprecated in v0.5.0. Attempting to read from the old format.")

	if err := json.Unmarshal(content, c.wrap.Data); err != nil {
		return fmt.Errorf("failed to unmarshal deprecated cache: %w", err)
	}

	c.wrap.RefreshedAt = time.Unix(0, 0)

	return nil
}

func (c *FileCache) Write(in any) error {
	c.wrap.Data = in

	marshaled, err := json.Marshal(c.wrap)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.MkdirAll(
		filepath.Dir(c.path),
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IXUSR,
	); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	if err := os.WriteFile(
		c.path,
		marshaled,
		syscall.S_IRUSR|syscall.S_IWUSR,
	); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}
