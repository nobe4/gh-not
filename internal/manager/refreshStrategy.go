package manager

import "fmt"

// RefreshStrategy is an enum for the refresh strategy.
// It implements https://pkg.go.dev/github.com/spf13/pflag#Value.
type RefreshStrategy int

const (
	AutoRefresh RefreshStrategy = iota
	ForceRefresh
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
		*r = ForceRefresh
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
