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
		"up":       []string{"up", "k"},
		"down":     []string{"down", "j"},
		"next":     []string{"right", "l"},
		"previous": []string{"left", "h"},
		"toggle":   []string{" "},
		"all":      []string{"a"},
		"none":     []string{"A"},
		"open":     []string{"o"},
		"filter":   []string{"/"},
		"command":  []string{":"},
		"help":     []string{"?"},
		"quit":     []string{"q", "esc", "ctrl+c"},
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
}

func (k KeyBinding) Help() string {
	return strings.NewReplacer(unicodeReplacement...).Replace(strings.Join(k, "|"))
}

func (k KeyBinding) Binding(help string) key.Binding {
	return key.NewBinding(
		key.WithKeys(k...),
		key.WithHelp(k.Help(), help),
	)
}
