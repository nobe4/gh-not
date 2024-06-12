/*
Package config provides a way to load the configuration from a file.
It also comes with a default configuration that can be used if no file is found.

See individual types for more information on the configuration.

Output the default configuration (free of rules) with `gh-not config --init`.

Example rules:

	rules:
	  - name: showcasing conditionals
	    action: debug
	    # Filters are run one after the other, like they are joined by 'and'.
	    # Having 'or' can be done via '(cond1) or (cond2) or ...'.
	    filters:
	      - .author.login == "dependabot[bot]"
	      - >
	        (.subject.title | contains("something unimportant")) or
	        (.subject.title | contains("something already done"))

	  - name: ignore ci failures for the current repo
	    action: done
	    filters:
	      - .repository.full_name == "nobe4/gh-not"
	      - .reason == "ci_activity"
*/
package config

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Cache    Cache    `yaml:"cache"`
	Endpoint Endpoint `yaml:"endpoint"`
	Keymap   Keymap   `yaml:"keymap"`
	Rules    []Rule   `yaml:"rules"`
}

type Endpoint struct {
	All      bool `yaml:"all"`
	MaxRetry int  `yaml:"max_retry"`
	MaxPage  int  `yaml:"max_page"`
}

type Cache struct {
	TTLInHours int    `yaml:"ttl_in_hours"`
	Path       string `yaml:"path"`
}

var (
	defaultCache = Cache{
		TTLInHours: 1,
		Path:       path.Join(StateDir(), "cache.json"),
	}

	defaultEndpoint = Endpoint{
		All:      true,
		MaxRetry: 10,
		MaxPage:  5,
	}
)

func Default() *Config {
	return &Config{
		Cache:    defaultCache,
		Endpoint: defaultEndpoint,
		Keymap:   defaultKeymap,
	}
}

func New(path string) (*Config, error) {
	config := Default()

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("config file not found, using default configuration", "path", path)
			return config, nil
		}

		return nil, err
	}

	slog.Warn("config file found", "path", path)

	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Save(path string) error {
	marshalled, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, marshalled, 0644)
}
