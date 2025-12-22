package config

import (
	"fmt"
	"strings"

	"github.com/nobe4/gh-not/internal/actions"
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

// Validate tests the rule for correctness. A rule must have an action and at least one filter
func (r Rule) Validate() error {
	reasons := []string{}

	actions := actions.GetMap(nil)
	if _, ok := actions[r.Action]; !ok {
		if r.Action == "" {
			reasons = append(reasons, "rule action is empty")
		} else {
			reasons = append(reasons, fmt.Sprintf("invalid rule action: \"%v\"", r.Action))
		}
	}

	if len(r.Filters) == 0 {
		reasons = append(reasons, "rule has no filters")
	}

	for _, filter := range r.Filters {
		if err := jq.Validate(filter); err != nil {
			reasons = append(reasons, fmt.Sprintf("failed to validate filter %s: %v", filter, err))
		}
	}

	if len(reasons) > 0 {
		return fmt.Errorf("failed to validate rule %q:\n - %s", r.Name, strings.Join(reasons, "\n - "))
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
		if n, err = jq.Filter(filter, n); err != nil {
			return nil, fmt.Errorf("failed to filter notifications: %w", err)
		}
	}

	return n, nil
}
