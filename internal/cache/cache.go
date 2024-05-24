package cache

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type ExpiringReadWriter interface {
	Read(any) error
	Write(any) error
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

func (c *FileCache) Read(out any) error {
	content, err := os.ReadFile(c.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
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
