package main

import (
	"time"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
)

const cachePath = "./cache-test.json"
const configPath = "./config.json"
const cacheTTL = time.Hour * 24 * 7

func main() {
	cache := cache.NewFileCache(cacheTTL, cachePath)

	client, err := gh.NewClient(cache)
	if err != nil {
		panic(err)
	}

	allNotifications, err := client.Notifications()
	if err != nil {
		panic(err)
	}

	config, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	allNotifications, err = config.Apply(allNotifications, actors.Map(client))
	if err != nil {
		panic(err)
	}

	if err := cache.Write(allNotifications); err != nil {
		panic(err)
	}
}
