package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cli/go-gh/v2/pkg/text"
	"github.com/spf13/cobra"

	"github.com/nobe4/gh-not/internal/api/github"
	configpkg "github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/jq"
	"github.com/nobe4/gh-not/internal/logger"
	managerpkg "github.com/nobe4/gh-not/internal/manager"
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

	config  *configpkg.Config
	manager *managerpkg.Manager

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

	errRuleNotFound = errors.New("rule not found")
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute the root command: %w", err)
	}

	return nil
}

//nolint:lll // Having the whole flag definition on a single line is OK.
func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().IntVarP(&verbosityFlag, "verbosity", "v", 1, "Change logger verbosity")
	rootCmd.PersistentFlags().StringVarP(&configPathFlag, "config", "c", "", "Path to the YAML config file")

	// Filter
	rootCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "List all the notifications, even the hidden/done ones.")
	rootCmd.Flags().StringVarP(&ruleFlag, "rule", "r", "", "Filter based on a rule name")
	rootCmd.Flags().StringVarP(&filterFlag, "filter", "f", "",
		"Filter with a jq expression passed into a select(...) call")
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

	slog.Debug("flags",
		"verbosity", verbosityFlag,
		"config", configPathFlag,
		"all", allFlag,
		"rule", ruleFlag,
		"filter", filterFlag,
		"tag", tagFlag,
		"tags", tagsFlag,
		"json", jsonFlag,
		"repl", replFlag,
	)

	var err error

	config, err = configpkg.New(configPathFlag)
	if err != nil {
		return fmt.Errorf("failed to load the config: %w", err)
	}

	manager = managerpkg.New(config.Data)

	return nil
}

func runRoot(_ *cobra.Command, _ []string) error {
	if err := manager.Load(); err != nil {
		return fmt.Errorf("failed to load the notifications: %w", err)
	}

	n := load()

	n, err := filter(n)
	if err != nil {
		slog.Error("Failed to filter the notifications", "err", err)

		return err
	}

	if err := display(n); err != nil {
		slog.Error("Failed to display the notifications", "err", err)

		return err
	}

	return nil
}

func load() notifications.Notifications {
	var n notifications.Notifications

	if allFlag {
		n = manager.Notifications
	} else {
		n = manager.Notifications.Visible()
	}

	n.Sort()

	return n
}

func filter(n notifications.Notifications) (notifications.Notifications, error) {
	var err error

	if filterFlag != "" {
		if n, err = jq.Filter(filterFlag, n); err != nil {
			return nil, fmt.Errorf("failed to filter the notifications: %w", err)
		}
	}

	if ruleFlag != "" {
		found := false

		for _, rule := range config.Data.Rules {
			if rule.Name == ruleFlag {
				found = true

				if n, err = rule.Filter(n); err != nil {
					return nil, fmt.Errorf("failed to filter the notifications: %w", err)
				}
			}
		}

		if !found {
			slog.Error("Rule not found", "rule", ruleFlag)

			return nil, fmt.Errorf("invalid rule '%s': %w", ruleFlag, errRuleNotFound)
		}
	}

	if tagFlag != "" {
		filter := fmt.Sprintf(`select(.meta.tags | index("%s"))`, tagFlag)

		if n, err = jq.Filter(filter, n); err != nil {
			return nil, fmt.Errorf("failed to filter the notifications: %w", err)
		}
	}

	return n, nil
}

func display(n notifications.Notifications) error {
	if tagsFlag {
		return displayTags(n)
	}

	if jsonFlag {
		return displayJSON(n)
	}

	if err := n.Render(); err != nil {
		slog.Warn("Failed to generate a table, using toString", "err", err)
	}

	if replFlag {
		return displayRepl(n)
	}

	displayTable(n)

	return nil
}

func displayTable(n notifications.Notifications) {
	//nolint:forbidigo // This is an expected print statement.
	fmt.Println(n)

	//nolint:forbidigo // This is an expected print statement.
	fmt.Printf("Found %d notifications %s\n",
		len(n),
		text.RelativeTimeAgo(
			time.Now(),
			manager.Cache.RefreshedAt(),
		),
	)
}

func displayJSON(n notifications.Notifications) error {
	marshaled, err := n.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal the notifications: %w", err)
	}

	//nolint:forbidigo // This is an expected print statement.
	fmt.Printf("%s\n", marshaled)

	return nil
}

func displayTags(n notifications.Notifications) error {
	for tag, count := range n.TagsMap() {
		//nolint:forbidigo // This is an expected print statement.
		fmt.Printf("%s: %d\n", tag, count)
	}

	return nil
}

func displayRepl(n notifications.Notifications) error {
	caller, err := github.New()
	if err != nil {
		return fmt.Errorf("failed to create an API REST client: %w", err)
	}

	manager.SetCaller(caller)

	// Launching bubbletea will occupy STDOUT and STDERR, so we need to redirect
	// the logs to a file.
	f, err := logger.InitWithFile(verbosityFlag, "/tmp/gh-not-debug.log")
	if err != nil {
		return fmt.Errorf("failed to init the logger: %w", err)
	}
	defer f.Close()

	if err := repl.Init(n, manager.Actions, config.Data.Keymap, config.Data.View); err != nil {
		return fmt.Errorf("failed to init the REPL: %w", err)
	}

	if err := manager.Save(); err != nil {
		return fmt.Errorf("failed to save the notifications: %w", err)
	}

	return nil
}
