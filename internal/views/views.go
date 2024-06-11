package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/colors"
)

type Mode int64

const (
	NormalMode Mode = iota
	SearchMode
	CommandMode
	ResultMode
	HelpMode
)

type ChangeModeMsg struct {
	Mode Mode
}

func ChangeMode(m Mode) tea.Cmd {
	return func() tea.Msg {
		return ChangeModeMsg{Mode: m}
	}
}

type ResultMsg struct {
	Out string
	Err error
}

func (r ResultMsg) ToString() string {
	out := r.Out

	if r.Err != nil {
		out += "\n" + colors.Red(r.Err.Error())
	}

	return out + "\npress any key to continue"
}
