package main

import (
	"fmt"
	"time"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
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

	fmt.Println(allNotifications.ToString())

	actorsMap := map[string]actors.Actor{
		"debug": &actors.DebugActor{},
		"print": &actors.PrintActor{},
		"hide":  &actors.HideActor{},
	}

	config, err := config.New("config.json")
	if err != nil {
		panic(err)
	}

	allNotifications, err = config.Apply(allNotifications, actorsMap)
	if err != nil {
		panic(err)
	}

	fmt.Println(allNotifications.ToString())

	if err := cache.Write(allNotifications); err != nil {
		panic(err)
	}
}
