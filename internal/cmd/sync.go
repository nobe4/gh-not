package cmd

import (
	"fmt"
	"time"

	"github.com/cli/go-gh/v2/pkg/text"
	"github.com/spf13/cobra"

	"github.com/nobe4/gh-not/internal/api"
	"github.com/nobe4/gh-not/internal/api/file"
	"github.com/nobe4/gh-not/internal/api/github"
	managerPkg "github.com/nobe4/gh-not/internal/manager"
)

var (
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

See synchronization logic at https://pkg.go.dev/github.com/nobe4/gh-not/internal/notifications#Sync.
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

//nolint:lll // Having the whole flag definition on a single line is OK.
func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().VarP(&forceStrategy, "force-strategy", "f", "Force strategy: "+forceStrategy.Allowed())
	syncCmd.Flags().VarP(&refreshStrategy, "refresh-strategy", "r", "Refresh strategy: "+refreshStrategy.Allowed())

	syncCmd.Flags().StringVarP(&notificationDumpPath, "from-file", "", "", "Path to notification dump in JSON (generate with 'gh api /notifications')")
}

func runSync(_ *cobra.Command, _ []string) error {
	var caller api.Requestor

	var err error

	if notificationDumpPath != "" {
		caller = file.New(notificationDumpPath)
	} else {
		caller, err = github.New()
		if err != nil {
			return fmt.Errorf("failed to create an API REST client: %w", err)
		}
	}

	manager.SetCaller(caller)

	manager.ForceStrategy = forceStrategy
	manager.RefreshStrategy = refreshStrategy

	if err := manager.Load(); err != nil {
		return fmt.Errorf("failed to load the notifications: %w", err)
	}

	loadedNotifications := len(manager.Notifications)

	if err := manager.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh the notifications: %w", err)
	}

	refreshedNotifications := len(manager.Notifications)

	if err := manager.Apply(); err != nil {
		return fmt.Errorf("failed to apply the rules: %w", err)
	}

	visibleNotifications := len(manager.Notifications.Visible())

	if err := manager.Save(); err != nil {
		return fmt.Errorf("failed to save the notifications: %w", err)
	}

	//nolint:forbidigo // This is an expected print statement.
	fmt.Printf("Loaded %d, refreshed %d, visible %d at %s\n",
		loadedNotifications,
		refreshedNotifications,
		visibleNotifications,
		text.RelativeTimeAgo(
			time.Now(),
			manager.Cache.RefreshedAt(),
		),
	)

	return nil
}
