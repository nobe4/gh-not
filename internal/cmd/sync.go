package cmd

import (
	"log/slog"

	"github.com/nobe4/gh-not/internal/actors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			notifications, err := client.Notifications()
			if err != nil {
				slog.Error("Failed to list notifications", "err", err)
				return err
			}

			notifications, err = config.Apply(notifications, actors.Map(client), noop)
			if err != nil {
				slog.Error("Failed to applying rules", "err", err)
				return err
			}

			if err := cache.Write(notifications); err != nil {
				slog.Error("Failed to write the cache", "err", err)
				return err
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVarP(&noop, "noop", "n", false, "Doesn't execute any action")
}
