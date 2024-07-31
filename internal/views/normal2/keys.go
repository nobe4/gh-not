package normal2

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/nobe4/gh-not/internal/config"
)

type Keymap struct {
	Toggle key.Binding
	Test   key.Binding
}

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Toggle}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Toggle},
		{k.Test},
	}
}

func (m *model) initKeymap(keymap config.Keymap) {
	m.keymap = Keymap{
		Toggle: keymap.Binding("normal", "toggle selected"),
		Test:   key.NewBinding(key.WithKeys("enter")),
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
