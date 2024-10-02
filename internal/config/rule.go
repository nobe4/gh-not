package config

import (
	"log/slog"

	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
)

// Rule is a struct to filter and act on notifications.
//
//	rules:
//	  - name: showcasing conditionals
//	    action: debug
//	    filters:
//	      - .author.login == "dependabot[bot]"
//	      - >
//	        (.subject.title | contains("something unimportant")) or
//	        (.subject.title | contains("something already done"))
//
//	  - name: ignore ci failures for the current repo
//	    action: done
//	    filters:
//	      - .repository.full_name == "nobe4/gh-not"
//	      - .reason == "ci_activity"
type Rule struct {
	Name string `mapstructure:"name"`

	// Filters is a list of jq filters to filter the notifications.
	// The filters are applied in order, like they are joined by 'and'.
	// Having 'or' can be done via '(cond1) or (cond2) or ...'.
	//
	// E.g.:
	// filters: ["A", "B or C"]
	// Will filter `A and (B or C)`.
	Filters []string `mapstructure:"filters"`

	// Action is the action to take on the filtered notifications.
	// See github.com/nobe4/internal/actions for list of available actions.
	Action string `mapstructure:"action"`

	// Args is the arguments to pass to the Action.
	Args []string `mapstructure:"args"`
}

// Test tests the rule for correctness.
func (r Rule) Test() error {
	for _, filter := range r.Filters {
		if err := jq.Validate(filter); err != nil {
			slog.Error("rule failed", "rule", r.Name, "filter", filter, "err", err)
			return err
		}
	}

	return nil
}

// Filter filters the notifications with the jq filters and returns the IDs.
func (r Rule) Filter(n notifications.Notifications) (notifications.Notifications, error) {
	var err error

	// TODO: accept only a single filter so that we can use the jq.Filter function
	// once instead of looping over the filters.
	// Given that it's joined by `and`, we can just write the filters directly.
	for _, filter := range r.Filters {
		n, err = jq.Filter(filter, n)
		if err != nil {
			return nil, err
		}
	}

	return n, nil
}
