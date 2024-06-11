package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/views/normal"
	"github.com/spf13/cobra"
)

var (
	replCmd = &cobra.Command{
		Use:     "repl",
		Aliases: []string{"r"},
		Short:   "Launch a REPL with notifications",
		RunE:    runRepl,
	}
)

func init() {
	rootCmd.AddCommand(replCmd)
}

func runRepl(cmd *cobra.Command, args []string) error {
	notifications := manager.Notifications

	renderCache, err := notifications.ToTable()
	if err != nil {
		return err
	}

	model := normal.New(client, notifications, renderCache)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
