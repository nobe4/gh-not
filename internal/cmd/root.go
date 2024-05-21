package cmd

import (
	"log/slog"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	cachePkg "github.com/nobe4/gh-not/internal/cache"
	configPkg "github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/spf13/cobra"
)

var (
	verbosity  int
	configPath string
	refresh    bool
	noRefresh  bool

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

	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.PersistentFlags().IntVarP(&verbosity, "verbosity", "v", 2, "Change logger verbosity")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yaml", "Path to the YAML config file")

	rootCmd.PersistentFlags().BoolVarP(&refresh, "refresh", "r", false, "Force a refresh")
	rootCmd.PersistentFlags().BoolVarP(&noRefresh, "no-refresh", "R", false, "Prevent a refresh")
	rootCmd.MarkFlagsMutuallyExclusive("refresh", "no-refresh")
}

func setupGlobals(cmd *cobra.Command, args []string) error {
	var err error

	config, err = configPkg.New(configPath)
	if err != nil {
		slog.Error("Failed to load the cache", "path", configPath, "err", err)
		return err
	}

	cache = cachePkg.NewFileCache(config.Cache.TTLInHours, config.Cache.Path)

	apiCaller, err := api.DefaultRESTClient()
	if err != nil {
		slog.Error("Failed to create an API REST client", "err", err)
		return err
	}

	client = gh.NewClient(apiCaller, cache, refresh, noRefresh)

	if err := initLogger(); err != nil {
		slog.Error("Failed to init the logger", "err", err)
		return err
	}

	return nil
}

func initLogger() error {
	opts := &slog.HandlerOptions{}

	switch verbosity {
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
