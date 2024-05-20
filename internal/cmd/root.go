package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbosity  int
	configPath string
	refresh    bool
	noRefresh  bool

	rootCmd = &cobra.Command{
		Use:   "gh-not",
		Short: "Manage your GitHub notifications",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initLogger()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(manageCmd)

	rootCmd.PersistentFlags().IntVarP(&verbosity, "verbosity", "v", 2, "Change logger verbosity")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yaml", "Path to the YAML config file")
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
