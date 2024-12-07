package repl

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"

	"github.com/nobe4/gh-not/internal/config"
)

type Keymap struct {
	Toggle key.Binding
	All    key.Binding
	None   key.Binding
	Open   key.Binding

	CommandMode key.Binding

	CommandAccept key.Binding
	CommandCancel key.Binding
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Toggle, k.All, k.None},
		{k.CommandMode, k.Open},
		{k.CommandAccept, k.CommandCancel},
	}
}

func (m *model) initKeymap(keymap config.Keymap) {
	m.keymap = Keymap{
		Toggle:        keymap.Binding("normal", "toggle selected"),
		All:           keymap.Binding("normal", "select all"),
		None:          keymap.Binding("normal", "select none"),
		Open:          keymap.Binding("normal", "open in browser"),
		CommandMode:   keymap.Binding("normal", "command mode"),
		CommandAccept: keymap.Binding("command", "command accept"),
		CommandCancel: keymap.Binding("command", "command cancel"),
	}

	m.list.KeyMap = list.KeyMap{
		CursorUp:             keymap.Binding("normal", "cursor up"),
		CursorDown:           keymap.Binding("normal", "cursor down"),
		PrevPage:             keymap.Binding("normal", "previous page"),
		NextPage:             keymap.Binding("normal", "next page"),
		GoToStart:            keymap.Binding("normal", "go to start"),
		GoToEnd:              keymap.Binding("normal", "go to end"),
		Filter:               keymap.Binding("normal", "filter mode"),
		ShowFullHelp:         keymap.Binding("normal", "toggle help"),
		CloseFullHelp:        keymap.Binding("normal", "toggle help"),
		Quit:                 keymap.Binding("normal", "quit"),
		ForceQuit:            keymap.Binding("normal", "force quit"),
		ClearFilter:          keymap.Binding("filter", "filter cancel"),
		CancelWhileFiltering: keymap.Binding("filter", "filter cancel"),
		AcceptWhileFiltering: keymap.Binding("filter", "filter accept"),
	}

	m.result.KeyMap = viewport.KeyMap{
		Down:     keymap.Binding("normal", "cursor down"),
		Up:       keymap.Binding("normal", "cursor up"),
		PageDown: keymap.Binding("normal", "next page"),
		PageUp:   keymap.Binding("normal", "previous page"),
	}
}

type ViewportKeymap struct {
	viewport.KeyMap
}

func (k ViewportKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.PageDown, k.PageUp},
	}
}
