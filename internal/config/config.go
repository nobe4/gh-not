/*
Package config provides a way to load the configuration from a file.
It also comes with a default configuration that can be used if no file is found.

See individual types for more information on the configuration.
*/
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/nobe4/gh-not/internal/gh"
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
	Cache    Cache       `mapstructure:"cache"`
	Endpoint gh.Endpoint `mapstructure:"endpoint"`
	Keymap   Keymap      `mapstructure:"keymap"`
	View     View        `mapstructure:"view"`
	Rules    []Rule      `mapstructure:"rules"`
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
	// Where to write logs when REPL is showing.
	LogPath string `mapstructure:"log_path"`
}

func Default(path string) (*viper.Viper, string) {
	slog.Debug("loading default configuration")

	if path == "" {
		path = filepath.Join(Dir(), "config.yaml")
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

//revive:disable:cognitive-complexity // TODO: simplify.
func New(path string) (*Config, error) {
	path, err := ExpandPathWithoutTilde(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand config path: %w", err)
	}

	v, path := Default(path)
	slog.Debug("loading configuration", "path", path)
	c := &Config{viper: v, Path: path}

	if err = c.viper.ReadInConfig(); err != nil {
		var viperNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &viperNotFoundError) &&
			!errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		slog.Warn("Config file not found, using default")
	}

	c.Path = c.viper.ConfigFileUsed()

	if err = c.viper.Unmarshal(&c.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validationErrors := []string{}
	for _, rule := range c.Data.Rules {
		if err = rule.Validate(); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}
	}
	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("invalid rules\n%s", strings.Join(validationErrors, "\n"))
	}

	c.Data.Cache.Path, err = ExpandPathWithoutTilde(c.Data.Cache.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand cache path: %w", err)
	}

	return c, nil
}

func (c *Config) Marshal() ([]byte, error) {
	//nolint:musttag // The struct is annotated with `mapstructure` tags already
	marshaled, err := yaml.Marshal(c.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	return marshaled, nil
}
