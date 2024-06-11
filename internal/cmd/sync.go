package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	noop bool

	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync notifications based on your config",
		Long: `
'gh-not sync' applies the ruleset on all the notifications in the cache.

Use this command when you want to make sure that your notification list is up to
date with your ruleset.
`,
		RunE: runSync,
	}
)

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVarP(&noop, "noop", "n", false, "Doesn't execute any action")
}

func runSync(cmd *cobra.Command, args []string) error {
	if err := manager.Apply(noop); err != nil {
		slog.Error("Failed to applying rules", "err", err)
		return err
	}

	if err := manager.Save(); err != nil {
		slog.Error("Failed to save the notifications", "err", err)
		return err
	}

	return nil
}
