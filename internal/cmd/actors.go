package cmd

import "github.com/spf13/cobra"

var (
	actionsCmd = &cobra.Command{
		Use:   "actions",
		Short: "Show information about the actions",
		Long: `
'gh-not' has multiple actions that perform different actions:

open: Open the notification in a web browser.

hide: Mark the notification as hidden in the cache.
      It won't show the notification again.

done: Mark the notification as done in the cache and on the API.
      It hides the notification until an update happens.

read: Mark the notification as read in the cache and on the API.
      It hides the notification until an update happens.

TODO: remove debug/print/pass
`,
	}
)

func init() {
	rootCmd.AddCommand(actionsCmd)
}
