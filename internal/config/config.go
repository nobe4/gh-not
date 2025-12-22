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

	"github.com/nobe4/dent.go"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/nobe4/gh-not/internal/gh"
)

const validationErrorStr = `Invalid rule (index %d): 
%s
Errors: 
%s`

var errRuleValidation = errors.New("invalid rules")

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

	if err = c.ValidateRules(); err != nil {
		return nil, err
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



func (c *Config) ValidateRules() error {
	validationErrors := []string{}

	for i, rule := range c.Data.Rules {
		violations := rule.Validate()
		if len(violations) == 0 {
			continue
		}

		errorStr := dent.IndentString(strings.Join(violations, "\n"), "  - ")

		yml, err := rule.Marshal()
		if err != nil {
			slog.Error("failed to marshal rule", "err", err)
		}

		valErr := fmt.Sprintf(validationErrorStr,
			i,
			dent.IndentString(string(yml), "  "), errorStr)
		validationErrors = append(validationErrors, valErr)
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("%w\n\n%s", errRuleValidation, strings.Join(validationErrors, "\n\n"))
	}

	return nil
}
