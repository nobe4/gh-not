package cmd

import (
	"log/slog"

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
	if err := manager.Load(); err != nil {
		slog.Error("Failed to load the notifications", "err", err)
		return err
	}

	notifications := manager.Notifications.Visible()

	renderCache, err := notifications.Table()
	if err != nil {
		return err
	}

	model := normal.New(manager.Actors, notifications, renderCache, config.Data.Keymap, config.Data.View)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	if err := manager.Save(); err != nil {
		slog.Error("Failed to save the notifications", "err", err)
		return err
	}

	return nil
}
