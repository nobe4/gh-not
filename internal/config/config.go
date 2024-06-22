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

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
	Data  *Data
}

type Data struct {
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
	c := &Config{viper: Default(path)}

	if err := c.viper.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) ||
			errors.Is(err, fs.ErrNotExist) {
			slog.Warn("Config file not found, using default")
		} else {
			slog.Error("Failed to read config file", "err", err)
			return nil, err
		}
	}

	if err := c.viper.Unmarshal(&c.Data); err != nil {
		return nil, err
	}

	return c, nil
}
