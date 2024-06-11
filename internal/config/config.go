package config

import (
	"log/slog"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Cache Cache  `yaml:"cache"`
	Rules []Rule `yaml:"rules"`
}

type Cache struct {
	TTLInHours int    `yaml:"ttl_in_hours"`
	Path       string `yaml:"path"`
}

var (
	defaultCache = Cache{
		TTLInHours: 24 * 7,
		Path:       path.Join(StateDir(), "cache.json"),
	}
)

func New(path string) (*Config, error) {
	config := &Config{Cache: defaultCache}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("config file not found, using default configuration", "path", path)
			return config, nil
		}

		return nil, err
	}

	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}
