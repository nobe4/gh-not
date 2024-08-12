package cmd

import (
	"fmt"
	"log/slog"

	"github.com/nobe4/gh-not/internal/api"
	"github.com/nobe4/gh-not/internal/api/file"
	"github.com/nobe4/gh-not/internal/api/github"
	managerPkg "github.com/nobe4/gh-not/internal/manager"
	"github.com/spf13/cobra"
)

var (
	noop                 bool
	notificationDumpPath string

	refreshStrategy managerPkg.RefreshStrategy
	forceStrategy   managerPkg.ForceStrategy

	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync notifications based on your config",
		Long: `
'gh-not sync' applies the ruleset on all the notifications in the cache.

Use this command when you want to make sure that your notification list is up to
date with your ruleset.

See synchronisation logic at https://pkg.go.dev/github.com/nobe4/gh-not/internal/notifications#Sync.
`,
		Example: `
  gh-not sync
  gh-not sync --force-strategy=noop,enrich
  gh-not sync --refresh-strategy=prevent
  gh-not sync --from-file=notifications.json
`,
		RunE: runSync,
	}
)

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().VarP(&forceStrategy, "force-strategy", "f", fmt.Sprintf("Force strategy: %s", forceStrategy.Allowed()))
	syncCmd.Flags().VarP(&refreshStrategy, "refresh-strategy", "r", fmt.Sprintf("Refresh strategy: %s", refreshStrategy.Allowed()))

	syncCmd.Flags().StringVarP(&notificationDumpPath, "from-file", "", "", "Path to notification dump in JSON (generate with 'gh api /notifications')")
}

func runSync(cmd *cobra.Command, args []string) error {
	var caller api.Requestor
	var err error

	if notificationDumpPath != "" {
		caller = file.New(notificationDumpPath)
	} else {
		caller, err = github.New()
		if err != nil {
			slog.Error("Failed to create an API REST client", "err", err)
			return err
		}
	}
	manager.SetCaller(caller)

	manager.ForceStrategy = forceStrategy
	manager.RefreshStrategy = refreshStrategy

	if err := manager.Load(); err != nil {
		slog.Error("Failed to load the notifications", "err", err)
		return err
	}
	loadedNotifications := len(manager.Notifications)

	if err := manager.Refresh(); err != nil {
		slog.Error("Failed to refresh the notifications", "err", err)
		return err
	}
	refreshedNotifications := len(manager.Notifications)

	if err := manager.Apply(); err != nil {
		slog.Error("Failed to applying rules", "err", err)
		return err
	}
	visibleNotifications := len(manager.Notifications.Visible())

	if err := manager.Save(); err != nil {
		slog.Error("Failed to save the notifications", "err", err)
		return err
	}

	fmt.Printf("Loaded %d, refreshed %d, visible %d\n", loadedNotifications, refreshedNotifications, visibleNotifications)

	return nil
}
