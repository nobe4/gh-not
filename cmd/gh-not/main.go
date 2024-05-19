package main

import (
	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
)

const configPath = "./config.yaml"

func main() {
	config, err := config.New(configPath)
	if err != nil {
		panic(err)
	}

	cache := cache.NewFileCache(config.Cache.TTLInHours, config.Cache.Path)

	client, err := gh.NewClient(cache)
	if err != nil {
		panic(err)
	}

	notifications, err := client.Notifications()
	if err != nil {
		panic(err)
	}

	notifications, err = config.Apply(notifications, actors.Map(client))
	if err != nil {
		panic(err)
	}

	if err := cache.Write(notifications); err != nil {
		panic(err)
	}
}
