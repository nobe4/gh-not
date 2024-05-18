package jq

import (
	"fmt"

	"github.com/itchyny/gojq"
	"github.com/nobe4/gh-not/internal/notifications"
)

// Filter applies a `.[] | select(filter)` on the notifications.
func Filter(filter string, n notifications.Notifications) (notifications.Notifications, error) {
	if filter == "" {
		return n, nil
	}

	query, err := gojq.Parse(fmt.Sprintf(".[] | select(%s)", filter))
	if err != nil {
		panic(err)
	}

	// gojq works only on any data, so we need to convert Notifications to
	// interface{}.
	// This also gives us back the JSON fields from the API.
	notificationsRaw, err := n.ToInterface()
	if err != nil {
		return nil, err
	}

	fitleredNotificationsRaw := []interface{}{}
	iter := query.Run(notificationsRaw)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			panic(err)
		}

		fitleredNotificationsRaw = append(fitleredNotificationsRaw, v)
	}

	filteredNotifications, err := notifications.FromInterface(fitleredNotificationsRaw)
	if err != nil {
		return nil, err
	}

	return filteredNotifications, nil
}
