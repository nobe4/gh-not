package cache

import (
	"errors"
	"os"
	"time"
)

type ExpiringReadWriter interface {
	Read() ([]byte, error)
	Write([]byte) error
	Expired() (bool, error)
}

type FileCache struct {
	path string
	ttl  time.Duration
}

func NewFileCache(ttl time.Duration, path string) *FileCache {
	return &FileCache{
		path: path,
		ttl:  ttl,
	}
}

func (c *FileCache) Read() ([]byte, error) {
	content, err := os.ReadFile(c.path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil
		}
		return nil, err
	}

	return content, nil
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

func (c *FileCache) Write(content []byte) error {
	return os.WriteFile(c.path, content, 0644)
}
