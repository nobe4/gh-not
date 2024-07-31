package config

import "path"

var Defaults = map[string]any{
	"cache.ttl_in_hours": 1,
	"cache.path":         path.Join(StateDir(), "cache.json"),

	"endpoint.all":       true,
	"endpoint.max_retry": 10,
	"endpoint.max_page":  5,

	"view.height": 40,

	"rules": []Rule{},

	"keymap.normal.cursor up":       []string{"up", "k"},
	"keymap.normal.cursor down":     []string{"down", "j"},
	"keymap.normal.next page":       []string{"right", "l"},
	"keymap.normal.previous page":   []string{"left", "h"},
	"keymap.normal.go to start":     []string{"home", "g"},
	"keymap.normal.go to end":       []string{"end", "G"},
	"keymap.normal.toggle selected": []string{" "},
	"keymap.normal.select all":      []string{"a"},
	"keymap.normal.select none":     []string{"A"},
	"keymap.normal.open in browser": []string{"o"},
	"keymap.normal.filter mode":     []string{"/"},
	"keymap.normal.command mode":    []string{":"},
	"keymap.normal.toggle help":     []string{"?"},
	"keymap.normal.quit":            []string{"q", "esc"},
	"keymap.normal.force quit":      []string{"ctrl+c"},

	"keymap.filter.confirm":       []string{"enter"},
	"keymap.filter.cancel":        []string{"esc"},
	"keymap.filter.accept filter": []string{"enter"},
	"keymap.filter.cancel filter": []string{"esc"},

	"keymap.command.confirm": []string{"enter"},
	"keymap.command.cancel":  []string{"esc"},
}
