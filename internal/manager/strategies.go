package manager

import (
	"fmt"
	"strings"
)

// RefreshStrategy is an enum for the refresh strategy.
// It implements https://pkg.go.dev/github.com/spf13/pflag#Value.
type RefreshStrategy int

const (
	// AutoRefresh refreshes the notifications if the cache is expired.
	AutoRefresh RefreshStrategy = iota

	// ForceRefresh always refreshes the notifications.
	ForceRefresh

	// PreventRefresh never refreshes the notifications.
	PreventRefresh
)

func (r RefreshStrategy) String() string {
	switch r {
	case AutoRefresh:
		return "auto"
	case ForceRefresh:
		return "force"
	case PreventRefresh:
		return "prevent"
	}
	return "unknown"
}

func (r *RefreshStrategy) Allowed() string {
	return "auto, force, prevent"
}

func (r *RefreshStrategy) Set(value string) error {
	switch value {
	case "auto":
		*r = AutoRefresh
	case "force":
		*r = ForceRefresh
	case "prevent":
		*r = PreventRefresh
	default:
		return fmt.Errorf(`must be one of %s`, r.Allowed())
	}

	return nil
}

func (r RefreshStrategy) Type() string {
	return "RefreshStrategy"
}

// ForceStrategy is an enum for the force strategy.
// It implements https://pkg.go.dev/github.com/spf13/pflag#Value.
type ForceStrategy int

const (
	// ForceApply forces the application of the ruleset on all notifications,
	// even the ones marked as Done.
	ForceApply ForceStrategy = 1 << iota

	// ForceApply forces the enrichment of all notifications, even the ones
	// marked as Done.
	ForceEnrich
)

func (r ForceStrategy) Has(s ForceStrategy) bool {
	return r&s != 0
}

func (r ForceStrategy) String() string {
	s := []string{}

	if r.Has(ForceApply) {
		s = append(s, "apply")
	}

	if r.Has(ForceEnrich) {
		s = append(s, "enrich")
	}

	return strings.Join(s, ", ")
}

func (r ForceStrategy) Allowed() string {
	return "apply, enrich"
}

func (r *ForceStrategy) Set(value string) error {
	v := strings.Split(value, ",")

	for _, s := range v {
		switch s {
		case "apply":
			*r |= ForceApply
		case "enrich":
			*r |= ForceEnrich
		default:
			return fmt.Errorf(`must be one of %s`, r.Allowed())
		}
	}

	return nil
}

func (r ForceStrategy) Type() string {
	return "ForceStrategy"
}
