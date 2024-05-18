package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/notifications"
)

const CacheTTL = time.Hour * 4

func main() {
	cache := cache.NewFileCache(CacheTTL, "/tmp/cache-test.json")

	client, err := gh.NewClient(cache)
	if err != nil {
		panic(err)
	}

	allNotifications, err := client.Notifications()
	if err != nil {
		panic(err)
	}

	fmt.Printf("all notifications %v\n", allNotifications)

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
