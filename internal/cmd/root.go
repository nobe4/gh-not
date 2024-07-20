package cmd

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobe4/gh-not/internal/api"
	"github.com/nobe4/gh-not/internal/api/file"
	"github.com/nobe4/gh-not/internal/api/github"
	configPkg "github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/logger"
	managerPkg "github.com/nobe4/gh-not/internal/manager"
	"github.com/nobe4/gh-not/internal/notifications"
	"github.com/nobe4/gh-not/internal/version"
	"github.com/nobe4/gh-not/internal/views/normal"
	"github.com/spf13/cobra"
)

var (
	verbosityFlag        int
	configPathFlag       string
	notificationDumpPath string
	refreshFlag          bool
	noRefreshFlag        bool
	filterFlag           = ""
	jqFlag               = ""
	repl                 bool

	config  *configPkg.Config
	manager *managerPkg.Manager

	rootCmd = &cobra.Command{
		Use:     "gh-not",
		Short:   "Lists your GitHub notifications",
		Version: version.String(),
		Example: `
  gh-not --no-refresh
`,
		PersistentPreRunE: setupGlobals,
		SilenceErrors:     true,
		RunE:              runRoot,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().IntVarP(&verbosityFlag, "verbosity", "v", 1, "Change logger verbosity")
	rootCmd.PersistentFlags().StringVarP(&configPathFlag, "config", "c", "", "Path to the YAML config file")

	rootCmd.Flags().StringVarP(&filterFlag, "filter", "f", "", "Filter with a jq expression passed into a select(...) call")
	rootCmd.Flags().StringVarP(&jqFlag, "jq", "q", "", "jq expression to run on the notification list")
	rootCmd.MarkFlagsMutuallyExclusive("filter", "jq")
	rootCmd.Flags().BoolVarP(&repl, "repl", "", false, "Start a REPL with the notifications list")

	// TODO: move to sync
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

	// TODO: move to sync
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

	// TODO: inject gh client only when necessary in the manager
	manager = managerPkg.New(config.Data, caller, refreshFlag, noRefreshFlag)

	return nil
}

func runRoot(cmd *cobra.Command, args []string) error {
	if err := manager.Load(); err != nil {
		slog.Error("Failed to load the notifications", "err", err)
		return err
	}

	notifications := manager.Notifications.Visible()

	if filterFlag != "" {
		notificationsList, err := jq.Filter(filterFlag, notifications)
		if err != nil {
			return err
		}
		notifications = notificationsList
	}

	if jqFlag != "" {
		return fmt.Errorf("`gh-not list --jq` implementation needed")
	}

	table, err := notifications.Table()
	if err != nil {
		slog.Warn("Failed to generate a table, using toString", "err", err)
		table = notifications.String()
	}

	if repl {
		return displayRepl(table, notifications)
	}

	displayTable(table, notifications)

	return nil
}

func displayTable(table string, notifications notifications.Notifications) {
	out := table
	out += fmt.Sprintf("\nFound %d notifications", len(notifications))

	fmt.Println(out)
}

func displayRepl(renderCache string, n notifications.Notifications) error {
	model := normal.New(manager.Actors, n, renderCache, config.Data.Keymap, config.Data.View)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	if err := manager.Save(); err != nil {
		slog.Error("Failed to save the notifications", "err", err)
		return err
	}

	return nil
}
