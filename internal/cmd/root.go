package cmd

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
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
	verbosityFlag  int
	configPathFlag string
	ruleFlag       string
	filterFlag     string
	replFlag       bool
	jsonFlag       bool
	allFlag        bool

	config  *configPkg.Config
	manager *managerPkg.Manager

	rootCmd = &cobra.Command{
		Use:     "gh-not",
		Short:   "Lists your GitHub notifications",
		Version: version.String(),
		Example: `
  gh-not --verbosity 2
  gh-not --config /path/to/config.yaml
  gh-not --filter '(.repository.full_name | contains("nobe4")) or (.subject.title | contains("CI"))'
  gh-not --json --all --rule 'ignore CI'
  gh-not --repl  // will log in the file /tmp/gh-not-debug.log
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

	rootCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "List all the notifications.")

	rootCmd.Flags().StringVarP(&ruleFlag, "rule", "r", "", "Filter based on a rule name")
	rootCmd.Flags().StringVarP(&filterFlag, "filter", "f", "", "Filter with a jq expression passed into a select(...) call")
	rootCmd.MarkFlagsMutuallyExclusive("rule", "filter")

	rootCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output the selected notifications as JSON")

	rootCmd.Flags().BoolVarP(&replFlag, "repl", "", false, "Start a REPL with the notifications list")
}

func setupGlobals(cmd *cobra.Command, args []string) error {
	logger.Init(verbosityFlag)

	var err error
	config, err = configPkg.New(configPathFlag)
	if err != nil {
		slog.Error("Failed to load the config", "path", configPathFlag, "err", err)
		return err
	}

	manager = managerPkg.New(config.Data)

	return nil
}

func runRoot(cmd *cobra.Command, args []string) error {
	if err := manager.Load(); err != nil {
		slog.Error("Failed to load the notifications", "err", err)
		return err
	}

	notifications := load()

	notifications, err := filter(notifications)
	if err != nil {
		slog.Error("Failed to filter the notifications", "err", err)
		return err
	}

	if err := display(notifications); err != nil {
		slog.Error("Failed to display the notifications", "err", err)
		return err
	}

	return nil
}

func load() notifications.Notifications {
	var notifications notifications.Notifications

	if allFlag {
		notifications = manager.Notifications
	} else {
		notifications = manager.Notifications.Visible()
	}

	notifications.Sort()

	return notifications
}

func filter(notifications notifications.Notifications) (notifications.Notifications, error) {
	if filterFlag != "" {
		notificationsList, err := jq.Filter(filterFlag, notifications)
		if err != nil {
			return nil, err
		}
		notifications = notificationsList
	}

	if ruleFlag != "" {
		found := false

		var err error
		for _, rule := range config.Data.Rules {
			if rule.Name == ruleFlag {
				found = true
				notifications, err = rule.Filter(notifications)
				if err != nil {
					return nil, err
				}
			}
		}

		if found == false {
			slog.Error("Rule not found", "rule", ruleFlag)
			return nil, fmt.Errorf("Rule '%s' not found", ruleFlag)
		}
	}

	return notifications, nil
}

func display(notifications notifications.Notifications) error {
	if jsonFlag {
		return displayJson(notifications)
	}

	if err := notifications.Render(); err != nil {
		slog.Warn("Failed to generate a table, using toString", "err", err)
	}

	if replFlag {
		return displayRepl(notifications)
	}

	displayTable(notifications)

	return nil
}

func displayTable(n notifications.Notifications) {
	fmt.Println(n)
	fmt.Printf("Found %d notifications\n", len(n))
	// TODO: add a notice if the notifications could be refreshed
}

func displayJson(notifications notifications.Notifications) error {
	marshaled, err := notifications.Marshal()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", marshaled)

	return nil
}

func displayRepl(n notifications.Notifications) error {
	caller, err := github.New()
	if err != nil {
		slog.Error("Failed to create an API REST client", "err", err)
		return err
	}
	manager.SetCaller(caller)

	// Launching bubbletea will occupy STDOUT and STDERR, so we need to redirect
	// the logs to a file.
	f, err := logger.InitWithFile(verbosityFlag, "/tmp/gh-not-debug.log")
	if err != nil {
		slog.Error("Failed to init the logger", "err", err)
		return err
	}
	defer f.Close()

	model := normal.New(manager.Actors, n, config.Data.Keymap, config.Data.View)
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
