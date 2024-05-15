package jq

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
	"github.com/nobe4/ghnot/internal/notifications"
)

// Filter applies a `.[] | select(filter)` on the notifications.
func Filter(filter string, n []notifications.Notification) ([]notifications.Notification, error) {
	query, err := gojq.Parse(fmt.Sprintf(".[] | select(%s)", filter))
	if err != nil {
		panic(err)
	}

	notificationsRaw, err := notificationsToInterface(n)
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

	filteredNotifications, err := interfaceToNotifications(fitleredNotificationsRaw)
	if err != nil {
		return nil, err
	}

	return filteredNotifications, nil
}

func notificationsToInterface(n []notifications.Notification) (interface{}, error) {
	// gojq works only on any data, so we need to convert []Notifications to
	// interface{}
	// This also gives us back the JSON fields from the API.
	marshalled, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal notifications: %w", err)
	}

	var i interface{}
	if err := json.Unmarshal(marshalled, &i); err != nil {
		return nil, fmt.Errorf("cannot unmarshal interface: %w", err)
	}

	return i, nil
}

func interfaceToNotifications(n interface{}) ([]notifications.Notification, error) {
	marshalled, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("cannot marshall interface: %w", err)
	}

	var notifications []notifications.Notification
	if err := json.Unmarshal(marshalled, &notifications); err != nil {
		return nil, fmt.Errorf("cannot unmarshall into notification: %w", err)
	}

	return notifications, nil
}
