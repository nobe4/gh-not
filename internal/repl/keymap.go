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

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Toggle}
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
		Toggle: keymap.Binding("normal", "toggle selected"),
		All:    keymap.Binding("normal", "select all"),
		None:   keymap.Binding("normal", "select none"),

		Open: keymap.Binding("normal", "open in browser"),

		CommandMode:   keymap.Binding("normal", "command mode"),
		CommandCancel: key.NewBinding(key.WithKeys("esc")),
		CommandAccept: key.NewBinding(key.WithKeys("enter")),
	}

	m.list.KeyMap = list.KeyMap{
		CursorUp:   keymap.Binding("normal", "cursor up"),
		CursorDown: keymap.Binding("normal", "cursor down"),
		PrevPage:   keymap.Binding("normal", "previous page"),
		NextPage:   keymap.Binding("normal", "next page"),
		GoToStart:  keymap.Binding("normal", "go to start"),
		GoToEnd:    keymap.Binding("normal", "go to end"),
		Filter:     keymap.Binding("normal", "filter mode"),
		// TODO: move all those to "normal" keymap
		ClearFilter:          keymap.Binding("filter", "cancel fitler"),
		CancelWhileFiltering: keymap.Binding("filter", "cancel filter"),
		AcceptWhileFiltering: keymap.Binding("filter", "accept filter"),
		ShowFullHelp:         keymap.Binding("normal", "toggle help"),
		CloseFullHelp:        keymap.Binding("normal", "toggle help"),
		Quit:                 keymap.Binding("normal", "quit"),
		ForceQuit:            keymap.Binding("normal", "force quit"),
	}
}

type ViewportKeymap struct {
	viewport.KeyMap
}

func (k ViewportKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.HalfPageUp, k.HalfPageDown},
		{k.PageDown, k.PageUp},
	}
}
