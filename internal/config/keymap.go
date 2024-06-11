package config

import "strings"

type Keymap map[string]KeyBindings
type KeyBindings map[string]KeyBinding

type KeyBinding []string

var defaultKeymap = Keymap{
	"normal": KeyBindings{
		"up":   []string{"up", "k"},
		"down": []string{"down", "j"},
	},
}

func (k KeyBinding) Help() string {
	return strings.Join(k, "|")
}
