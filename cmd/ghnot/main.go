package main

import (
	"fmt"

	"github.com/nobe4/ghnot/internal/config"
	"github.com/nobe4/ghnot/internal/gh"
	"github.com/nobe4/ghnot/internal/jq"
	"github.com/nobe4/ghnot/internal/notifications"
)

func main() {
	allNotifications, err := gh.Run([]string{"api", "/notifications?all=true&per_page=3"})

	if err != nil {
		panic(err)
	}

	filteredNotifications := []notifications.Notification{}

	config, err := config.Load("config.json")
	if err != nil {
		panic(err)
	}

	for _, group := range config.Groups {
		for _, filter := range group.Filters {
			selectedNotifications, err := jq.Filter(filter, allNotifications)
			if err != nil {
				panic(err)
			}

			filteredNotifications = append(filteredNotifications, selectedNotifications...)
		}
		filteredNotifications = notifications.Uniq(filteredNotifications)

		for _, notification := range filteredNotifications {
			fmt.Printf("%s %v\n", group.Action, notification)
		}
	}
}
