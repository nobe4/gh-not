package config

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

type (
	Keymap      map[string]KeyBindings
	KeyBindings map[string]KeyBinding
)

type KeyBinding []string

//nolint:gochecknoglobals // This replacement list is used a lot.
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

func (k Keymap) Binding(mode, action string) key.Binding {
	keys := k[mode][action]

	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys.Help(), action),
	)
}

func (k KeyBinding) Help() string {
	return strings.NewReplacer(unicodeReplacement...).Replace(strings.Join(k, "|"))
}
