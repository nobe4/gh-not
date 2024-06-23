/*
Package config provides a way to load the configuration from a file.
It also comes with a default configuration that can be used if no file is found.

See individual types for more information on the configuration.

Output the default configuration (free of rules) with `gh-not config --init`.
*/
package config

import (
	"errors"
	"io/fs"
	"log/slog"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config holds the configuration data.
type Config struct {
	viper *viper.Viper
	Path  string
	Data  *Data
}

// Data holds the configuration data.
type Data struct {
	Cache    Cache    `yaml:"cache"`
	Endpoint Endpoint `yaml:"endpoint"`
	Keymap   Keymap   `yaml:"keymap"`
	View     View     `yaml:"view"`
	Rules    []Rule   `yaml:"rules"`
}

// Endpoint is the configuration for the GitHub API endpoint.
type Endpoint struct {
	// Pull all notifications from the endpoint.
	// By default, only the unread notifications are fetched.
	// This maps to `?all=true|false` in the GitHub API.
	// See https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#list-notifications-for-the-authenticated-user
	All bool `yaml:"all"`

	// The maximum number of retries to fetch notifications.
	// The Notifications API is notably flaky, retrying HTTP requests is
	// definitely needed.
	MaxRetry int `yaml:"max_retry"`

	// The number of notification pages to fetch.
	// This will cap the `?page=X` parameter in the GitHub API.
	// See https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#list-notifications-for-the-authenticated-user
	MaxPage int `yaml:"max_page"`
}

// Cache is the configuration for the cache file.
type Cache struct {
	// The path to the cache file.
	Path string `yaml:"path"`

	// The time-to-live of the cache in hours.
	TTLInHours int `yaml:"ttl_in_hours"`
}

// View is the configuration for the terminal view.
type View struct {
	// Number of notifications to display at once.
	Height int `yaml:"height"`
}

func Default(path string) *viper.Viper {
	slog.Debug("loading default configuration")
	v := viper.New()

	for key, value := range Defaults {
		v.SetDefault(key, value)
	}

	slog.Debug("setting config name and path", "path", path)
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
	slog.Debug("loading configuration", "path", path)
	c := &Config{viper: Default(path), Path: path}

	if err := c.viper.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) ||
			errors.Is(err, fs.ErrNotExist) {
			slog.Warn("Config file not found, using default")
		} else {
			slog.Error("Failed to read config file", "err", err)
			return nil, err
		}
	}
	c.Path = c.viper.ConfigFileUsed()

	if err := c.viper.Unmarshal(&c.Data); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Marshal() ([]byte, error) {
	marshalled, err := yaml.Marshal(c.Data)
	if err != nil {
		slog.Error("Failed to marshall config", "err", err)
		return nil, err
	}

	return marshalled, nil
}
