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
	force                bool
	notificationDumpPath string
	refreshStrategy      managerPkg.RefreshStrategy

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
	syncCmd.Flags().BoolVarP(&force, "force", "f", false, "Force the execution of the rules on Done notifications")
	syncCmd.Flags().VarP(&refreshStrategy, "refresh-strategy", "r", fmt.Sprintf("Refresh strategy: %s", refreshStrategy.Allowed()))
	syncCmd.Flags().StringVarP(&notificationDumpPath, "from-file", "", "", "Path to notification dump in JSON (generate with 'gh api /notifications')")
}

func runSync(cmd *cobra.Command, args []string) error {
	var caller api.Caller
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
	manager.WithRefresh(refreshStrategy).WithCaller(caller)

	if err := manager.Load(); err != nil {
		slog.Error("Failed to load the notifications", "err", err)
		return err
	}

	if err := manager.Apply(noop, force); err != nil {
		slog.Error("Failed to applying rules", "err", err)
		return err
	}

	if err := manager.Save(); err != nil {
		slog.Error("Failed to save the notifications", "err", err)
		return err
	}

	return nil
}
