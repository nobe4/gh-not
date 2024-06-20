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
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
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

var Defaults = map[string]any{
	"cache.ttl_in_hours": 1,
	"cache.path":         path.Join(StateDir(), "cache.json"),

	"endpoint.all":       true,
	"endpoint.max_retry": 10,
	"endpoint.max_page":  5,

	"rules": []Rule{},

	"keymap.normal.cursor up":       []string{"up", "k"},
	"keymap.normal.cursor down":     []string{"down", "j"},
	"keymap.normal.next page":       []string{"right", "l"},
	"keymap.normal.previous page":   []string{"left", "h"},
	"keymap.normal.toggle selected": []string{" "},
	"keymap.normal.select all":      []string{"a"},
	"keymap.normal.select none":     []string{"A"},
	"keymap.normal.open in browser": []string{"o"},
	"keymap.normal.filter mode":     []string{"/"},
	"keymap.normal.command mode":    []string{":"},
	"keymap.normal.toggle help":     []string{"?"},
	"keymap.normal.quit":            []string{"q", "esc", "ctrl+c"},

	"keymap.filter.confirm": []string{"enter"},
	"keymap.filter.cancel":  []string{"esc", "ctrl+c"},

	"keymap.command.confirm": []string{"enter"},
	"keymap.command.cancel":  []string{"esc", "ctrl+c"},
}

func Default(path string) *viper.Viper {
	v := viper.New()

	for key, value := range Defaults {
		v.SetDefault(key, value)
	}

	if path == "" {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(ConfigDir())
	} else {
		v.SetConfigFile(path)
	}

	return v
}

func New(path string) (*Config, error) {
	v := Default(path)

	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) ||
			errors.Is(err, fs.ErrNotExist) {
			slog.Warn("Config file not found, using default")
		} else {
			slog.Error("Failed to read config file", "err", err)
			return nil, err
		}
	}

	config := &Config{}
	if err := v.Unmarshal(&config); err != nil {
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
