package cmd

import (
	"log/slog"

	"github.com/nobe4/gh-not/internal/api"
	"github.com/nobe4/gh-not/internal/api/file"
	"github.com/nobe4/gh-not/internal/api/github"
	configPkg "github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/logger"
	managerPkg "github.com/nobe4/gh-not/internal/manager"
	"github.com/nobe4/gh-not/internal/version"
	"github.com/spf13/cobra"
)

var (
	verbosityFlag        int
	configPathFlag       string
	notificationDumpPath string
	refreshFlag          bool
	noRefreshFlag        bool

	config  *configPkg.Config
	manager *managerPkg.Manager

	rootCmd = &cobra.Command{
		Use:   "gh-not",
		Short: "Manage your GitHub notifications",
		Example: `
  gh-not --config list
  gh-not --no-refresh list
  gh-not --from-file notifications.json list
  gh-not sync --refresh --verbosity 4
`,
		PersistentPreRunE: setupGlobals,
		SilenceErrors:     true,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Version = version.String()

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().IntVarP(&verbosityFlag, "verbosity", "v", 1, "Change logger verbosity")
	rootCmd.PersistentFlags().StringVarP(&configPathFlag, "config", "c", "", "Path to the YAML config file")

	rootCmd.PersistentFlags().StringVarP(&notificationDumpPath, "from-file", "", "", "Path to notification dump in JSON (generate with 'gh api /notifications')")

	rootCmd.PersistentFlags().BoolVarP(&refreshFlag, "refresh", "r", false, "Force a refresh")
	rootCmd.PersistentFlags().BoolVarP(&noRefreshFlag, "no-refresh", "R", false, "Prevent a refresh")
	rootCmd.MarkFlagsMutuallyExclusive("refresh", "no-refresh")
}

func setupGlobals(cmd *cobra.Command, args []string) error {
	if err := logger.Init(verbosityFlag); err != nil {
		slog.Error("Failed to init the logger", "err", err)
		return err
	}

	var err error

	config, err = configPkg.New(configPathFlag)
	if err != nil {
		slog.Error("Failed to load the config", "path", configPathFlag, "err", err)
		return err
	}

	var caller api.Caller
	if notificationDumpPath != "" {
		caller = file.New(notificationDumpPath)
	} else {
		caller, err = github.New()
		if err != nil {
			slog.Error("Failed to create an API REST client", "err", err)
			return err
		}
	}

	manager = managerPkg.New(config.Data, caller, refreshFlag, noRefreshFlag)

	return nil
}
