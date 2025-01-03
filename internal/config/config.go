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
	"path/filepath"

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
	Cache    Cache    `mapstructure:"cache"`
	Endpoint Endpoint `mapstructure:"endpoint"`
	Keymap   Keymap   `mapstructure:"keymap"`
	View     View     `mapstructure:"view"`
	Rules    []Rule   `mapstructure:"rules"`
}

// Endpoint is the configuration for the GitHub API endpoint.
//
//nolint:lll
type Endpoint struct {
	// Pull all notifications from the endpoint.
	// By default, only the unread notifications are fetched.
	// This maps to `?all=true|false` in the GitHub API.
	// See https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#list-notifications-for-the-authenticated-user
	All bool `mapstructure:"all"`

	// The maximum number of retries to fetch notifications.
	// The Notifications API is notably flaky, retrying HTTP requests is
	// definitely needed.
	MaxRetry int `mapstructure:"max_retry"`

	// The number of notification pages to fetch.
	// This will cap the `?page=X` parameter in the GitHub API.
	// See https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#list-notifications-for-the-authenticated-user
	MaxPage int `mapstructure:"max_page"`

	// The number of notifications to fetch per page.
	// This maps to `?per_page=X` in the GitHub API.
	// See https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28#list-notifications-for-the-authenticated-user
	PerPage int `mapstructure:"per_page"`
}

// Cache is the configuration for the cache file.
type Cache struct {
	// The path to the cache file.
	Path string `mapstructure:"path"`

	// The time-to-live of the cache in hours.
	TTLInHours int `mapstructure:"ttl_in_hours"`
}

// View is the configuration for the terminal view.
type View struct {
	// Number of notifications to display at once.
	Height int `mapstructure:"height"`
}

func Default(path string) (*viper.Viper, string) {
	slog.Debug("loading default configuration")

	if path == "" {
		path = filepath.Join(ConfigDir(), "config.yaml")
		slog.Debug("path is empty, setting default path", "default path", path)
	}

	v := viper.New()

	for key, value := range Defaults {
		v.SetDefault(key, value)
	}

	slog.Debug("setting config name and path", "path", path)

	v.SetConfigFile(path)

	return v, path
}

func New(path string) (*Config, error) {
	v, path := Default(path)
	slog.Debug("loading configuration", "path", path)
	c := &Config{viper: v, Path: path}

	if err := c.viper.ReadInConfig(); err != nil {
		var viperNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &viperNotFoundError) ||
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

	for _, rule := range c.Data.Rules {
		if err := rule.Test(); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Config) Marshal() ([]byte, error) {
	marshaled, err := yaml.Marshal(c.Data)
	if err != nil {
		slog.Error("Failed to marshall config", "err", err)
		return nil, err
	}

	return marshaled, nil
}
