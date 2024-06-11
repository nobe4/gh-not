package config

import (
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Rule struct {
	Name    string   `yaml:"name"`
	Filters []string `yaml:"filters"`
	Action  string   `yaml:"action"`
}

func (r Rule) FilterIds(n notifications.Notifications) ([]string, error) {
	var err error

	for _, filter := range r.Filters {
		n, err = jq.Filter(filter, n)
		if err != nil {
			return nil, err
		}
	}

	return n.IDList(), nil
}
