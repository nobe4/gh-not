package cmd

import (
	"log/slog"
	"os"
	"path"

	"github.com/cli/go-gh/v2/pkg/api"
	cachePkg "github.com/nobe4/gh-not/internal/cache"
	configPkg "github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/spf13/cobra"
)

var (
	verbosityFlag  int
	configPathFlag string
	refreshFlag    bool
	noRefreshFlag  bool

	config *configPkg.Config
	cache  *cachePkg.FileCache
	client *gh.Client

	rootCmd = &cobra.Command{
		Use:               "gh-not",
		Short:             "Manage your GitHub notifications",
		PersistentPreRunE: setupGlobals,
		SilenceErrors:     true,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().IntVarP(&verbosityFlag, "verbosity", "v", 2, "Change logger verbosity")
	rootCmd.PersistentFlags().StringVarP(&configPathFlag, "config", "c", path.Join(configPkg.ConfigDir(), "config.yaml"), "Path to the YAML config file")

	rootCmd.PersistentFlags().BoolVarP(&refreshFlag, "refresh", "r", false, "Force a refresh")
	rootCmd.PersistentFlags().BoolVarP(&noRefreshFlag, "no-refresh", "R", false, "Prevent a refresh")
	rootCmd.MarkFlagsMutuallyExclusive("refresh", "no-refresh")
}

func setupGlobals(cmd *cobra.Command, args []string) error {
	var err error

	config, err = configPkg.New(configPathFlag)
	if err != nil {
		slog.Error("Failed to load the cache", "path", configPathFlag, "err", err)
		return err
	}

	cache = cachePkg.NewFileCache(config.Cache.TTLInHours, config.Cache.Path)

	apiCaller, err := api.DefaultRESTClient()
	if err != nil {
		slog.Error("Failed to create an API REST client", "err", err)
		return err
	}

	client = gh.NewClient(apiCaller, cache, refreshFlag, noRefreshFlag)

	if err := initLogger(); err != nil {
		slog.Error("Failed to init the logger", "err", err)
		return err
	}

	return nil
}

func initLogger() error {
	opts := &slog.HandlerOptions{}

	switch verbosityFlag {
	case 1:
		opts.Level = slog.LevelError
	case 2:
		opts.Level = slog.LevelWarn
	case 3:
		opts.Level = slog.LevelInfo
	case 4:
		opts.Level = slog.LevelDebug
	case 5:
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}
