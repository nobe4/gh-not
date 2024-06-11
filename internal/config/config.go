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

const Example = `
---

cache:
  ttl_in_hours: 1
  path: ./cache.json

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
`

var (
	defaultCache = Cache{
		TTLInHours: 1,
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
