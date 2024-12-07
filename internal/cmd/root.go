package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/cli/go-gh/v2/pkg/text"
	"github.com/spf13/cobra"

	"github.com/nobe4/gh-not/internal/api/github"
	configPkg "github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/logger"
	managerPkg "github.com/nobe4/gh-not/internal/manager"
	"github.com/nobe4/gh-not/internal/notifications"
	"github.com/nobe4/gh-not/internal/repl"
	"github.com/nobe4/gh-not/internal/version"
)

var (
	verbosityFlag  int
	configPathFlag string
	ruleFlag       string
	filterFlag     string
	tagFlag        string
	tagsFlag       bool
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
  gh-not --tag tag0
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

	// Filter
	rootCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "List all the notifications, even the hidden/done ones.")
	rootCmd.Flags().StringVarP(&ruleFlag, "rule", "r", "", "Filter based on a rule name")
	rootCmd.Flags().StringVarP(&filterFlag, "filter", "f", "", "Filter with a jq expression passed into a select(...) call")
	rootCmd.Flags().StringVarP(&tagFlag, "tag", "t", "", "Filter from a single tag")
	rootCmd.MarkFlagsMutuallyExclusive("rule", "filter", "tag")

	// Display
	rootCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output the selected notifications as JSON")
	rootCmd.Flags().BoolVarP(&tagsFlag, "tags", "", false, "Show the list of tags with associated notification count")
	rootCmd.Flags().BoolVarP(&replFlag, "repl", "", false, "Start a REPL with the notifications list")
	rootCmd.MarkFlagsMutuallyExclusive("json", "repl", "tags")
}

func setupGlobals(_ *cobra.Command, _ []string) error {
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

func runRoot(_ *cobra.Command, _ []string) error {
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
	var err error

	if filterFlag != "" {
		if notifications, err = jq.Filter(filterFlag, notifications); err != nil {
			return nil, err
		}
	}

	if ruleFlag != "" {
		found := false

		for _, rule := range config.Data.Rules {
			if rule.Name == ruleFlag {
				found = true

				if notifications, err = rule.Filter(notifications); err != nil {
					return nil, err
				}
			}
		}

		if !found {
			slog.Error("Rule not found", "rule", ruleFlag)
			return nil, fmt.Errorf("Rule '%s' not found", ruleFlag)
		}
	}

	if tagFlag != "" {
		filter := fmt.Sprintf(`select(.meta.tags | index("%s"))`, tagFlag)

		if notifications, err = jq.Filter(filter, notifications); err != nil {
			return nil, err
		}
	}

	return notifications, nil
}

func display(notifications notifications.Notifications) error {
	if tagsFlag {
		return displayTags(notifications)
	}

	if jsonFlag {
		return displayJSON(notifications)
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
	fmt.Printf("Found %d notifications %s\n",
		len(n),
		text.RelativeTimeAgo(
			time.Now(),
			manager.Cache.RefreshedAt(),
		),
	)
}

func displayJSON(notifications notifications.Notifications) error {
	marshaled, err := notifications.Marshal()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", marshaled)

	return nil
}

func displayTags(n notifications.Notifications) error {
	for tag, count := range n.TagsMap() {
		fmt.Printf("%s: %d\n", tag, count)
	}

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

	if err := repl.Init(n, manager.Actions, config.Data.Keymap, config.Data.View); err != nil {
		return err
	}

	if err := manager.Save(); err != nil {
		slog.Error("Failed to save the notifications", "err", err)
		return err
	}

	return nil
}
