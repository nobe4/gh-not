package manager

import (
	"errors"
	"fmt"
	"log/slog"
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

var errNotAllowed = errors.New("not allowed")

func (r *RefreshStrategy) String() string {
	switch *r {
	case AutoRefresh:
		return "auto"
	case ForceRefresh:
		return "force"
	case PreventRefresh:
		return "prevent"
	}

	return "unknown"
}

func (*RefreshStrategy) Allowed() string {
	return "auto, force, prevent"
}

func (r *RefreshStrategy) ShouldRefresh(expired bool) bool {
	switch *r {
	case ForceRefresh:
		slog.Info("forcing a refresh")

		return true

	case PreventRefresh:
		slog.Info("preventing a refresh")

		return false

	case AutoRefresh:
		fallthrough
	default:
		slog.Debug("refresh based on cache expiration", "expired", expired)

		return expired
	}
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
		return fmt.Errorf(`%s must be one of %s: %w`, value, r.Allowed(), errNotAllowed)
	}

	return nil
}

func (*RefreshStrategy) Type() string {
	return "RefreshStrategy"
}

// ForceStrategy is an enum for the force strategy.
// It implements https://pkg.go.dev/github.com/spf13/pflag#Value.
type ForceStrategy int

const (
	// ForceApply forces the application of the ruleset on all notifications,
	// even the ones marked as Done.
	ForceApply ForceStrategy = 1 << iota

	// ForceNoop prevents any Action from being executed.
	ForceNoop

	// ForceApply forces the enrichment of all notifications, even the ones
	// marked as Done.
	ForceEnrich
)

func (r *ForceStrategy) Has(s ForceStrategy) bool {
	return *r&s != 0
}

func (r *ForceStrategy) String() string {
	s := []string{}

	if r.Has(ForceApply) {
		s = append(s, "apply")
	}

	if r.Has(ForceNoop) {
		s = append(s, "noop")
	}

	if r.Has(ForceEnrich) {
		s = append(s, "enrich")
	}

	return strings.Join(s, ", ")
}

func (*ForceStrategy) Allowed() string {
	return "apply, noop, enrich"
}

func (r *ForceStrategy) Set(value string) error {
	v := strings.Split(value, ",")

	for _, s := range v {
		switch s {
		case "apply":
			*r |= ForceApply
		case "noop":
			*r |= ForceNoop
		case "enrich":
			*r |= ForceEnrich
		default:
			return fmt.Errorf(`%s must be one of %s: %w`, s, r.Allowed(), errNotAllowed)
		}
	}

	return nil
}

func (*ForceStrategy) Type() string {
	return "ForceStrategy"
}
