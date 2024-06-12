package config

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

type Keymap map[string]KeyBindings
type KeyBindings map[string]KeyBinding

type KeyBinding []string

var defaultKeymap = Keymap{
	"normal": KeyBindings{
		"cursor up":       []string{"up", "k"},
		"cursor down":     []string{"down", "j"},
		"next page":       []string{"right", "l"},
		"previous page":   []string{"left", "h"},
		"toggle selected": []string{" "},
		"select all":      []string{"a"},
		"select none":     []string{"A"},
		"open in browser": []string{"o"},
		"filter mode":     []string{"/"},
		"command mode":    []string{":"},
		"toggle help":     []string{"?"},
		"quit":            []string{"q", "esc", "ctrl+c"},
	},
	"filter": KeyBindings{
		"confirm": []string{"enter"},
		"cancel":  []string{"esc", "ctrl+c"},
	},
	"command": KeyBindings{
		"confirm": []string{"enter"},
		"cancel":  []string{"esc", "ctrl+c"},
	},
}

var unicodeReplacement = []string{
	"up", "↑",
	"down", "↓",
	"left", "←",
	"right", "→",
	" ", "␣",
	"tab", "⇥",
	"enter", "↵",
	"esc", "⎋",
	"ctrl", "C",
}

func (k Keymap) Keybinding(mode, action string) key.Binding {
	keys := k[mode][action]
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys.Help(), action),
	)
}

func (k KeyBinding) Help() string {
	return strings.NewReplacer(unicodeReplacement...).Replace(strings.Join(k, "|"))
}
