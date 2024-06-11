package config

import "github.com/charmbracelet/bubbles/key"

type Keymap map[string]KeyBindings
type KeyBindings map[string]KeyBinding

type KeyBinding struct {
	Keys []string `yaml:"keys"`
	Help Help     `yaml:"help"`
}

type Help struct {
	Keys        string `yaml:"keys"`
	Description string `yaml:"description"`
}

var defaultKeymap = Keymap{
	"normal": KeyBindings{
		"up": {
			Keys: []string{"up", "k"},
			Help: Help{
				Keys:        "↑/k",
				Description: "move test up",
			},
		},
		"down": {
			Keys: []string{"down", "j"},
			Help: Help{
				Keys:        "↓/j",
				Description: "move down",
			},
		},
	},
}

func (k Keymap) Binding(mode, key string) key.Binding {
	return k[mode][key].Binding()
}

func (k KeyBinding) Binding() key.Binding {
	return key.NewBinding(
		key.WithKeys(k.Keys...),
		key.WithHelp(k.Help.Keys, k.Help.Description),
	)
}
