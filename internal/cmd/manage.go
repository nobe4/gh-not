package cmd

import (
	"log/slog"

	"github.com/nobe4/gh-not/internal/actors"
	"github.com/nobe4/gh-not/internal/cache"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/spf13/cobra"
)

var (
	noop bool

	manageCmd = &cobra.Command{
		Use:   "manage",
		Short: "Manage your notification based on your config",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := config.New(configPath)
			if err != nil {
				slog.Error("Failed to load the cache", "path", configPath, "err", err)
				return err
			}

			cache := cache.NewFileCache(config.Cache.TTLInHours, config.Cache.Path)

			client, err := gh.NewClient(cache)
			if err != nil {
				slog.Error("Failed to create a gh client", "err", err)
				return err
			}

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
	manageCmd.Flags().BoolVarP(&noop, "no-op", "n", false, "Doesn't execute any action")
}
