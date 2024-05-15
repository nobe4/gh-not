package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Name    string   `json:"name"`
	Filters []string `json:"filters"`
	Action  string   `json:"action"`
}

func Load(path string) (Config, error) {
	var config Config

	content, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(content, &config); err != nil {
		return config, err
	}

	return config, nil
}
