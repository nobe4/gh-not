package main

import (
	"log"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
)

func main() {
	allNotifications, err := gh.Run([]string{"api", "/notifications?all=true&per_page=3"})

	if err != nil {
		panic(err)
	}

	filteredNotifications := []notifications.Notification{}

	actorsMap := map[string]actors.Actor{
		"debug": &actors.DebugActor{},
		"print": &actors.PrintActor{},
	}

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
			if actor, ok := actorsMap[group.Action]; ok == true {
				actor.Run(notification)
			} else {
				log.Fatalf("unknown action '%s'", group.Action)
			}
		}
	}
}
