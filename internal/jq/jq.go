package jq

import (
	"errors"
	"fmt"

	"github.com/itchyny/gojq"

	"github.com/nobe4/gh-not/internal/notifications"
)

var errInvalidID = errors.New("invalid ID")

// TODO: refactor this as a callback to be called on n.Filter(flt) and have
// n.Filter call .Compact
// Filter applies a `.[] | select(filter)` on the notifications.
//
//revive:disable:cognitive-complexity // TODO: simplify.
func Filter(filter string, n notifications.Notifications) (notifications.Notifications, error) {
	if filter == "" {
		return n, nil
	}

	// Extract the IDs from the notifications that match the filter.
	// FilterFromIds will be used to get back the notifications from the IDs.
	query, err := gojq.Parse(fmt.Sprintf(".[] | select(%s) | .id", filter))
	if err != nil {
		return nil, fmt.Errorf("failed to parse filter: %w", err)
	}

	// gojq works only on `any` data, so we need to convert Notifications to
	// interface{}. This also gives us back the JSON fields from the API.
	notificationsRaw, err := n.Interface()
	if err != nil {
		return nil, fmt.Errorf("failed to convert notifications to raw interface: %w", err)
	}

	filteredIDs := []string{}

	iter := query.Run(notificationsRaw)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok = v.(error); ok {
			haltError := &gojq.HaltError{}
			if ok = errors.As(err, &haltError); ok && haltError.Value() == nil {
				break
			}

			return nil, err
		}

		newID, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("invalid filtered id %#v: %w", v, errInvalidID)
		}

		filteredIDs = append(filteredIDs, newID)
	}

	return n.FilterFromIDs(filteredIDs), nil
}

func Validate(filter string) error {
	if _, err := gojq.Parse(filter); err != nil {
		return fmt.Errorf("failed to parse filter: %w", err)
	}

	return nil
}
