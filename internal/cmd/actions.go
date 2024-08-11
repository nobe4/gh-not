package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed actions-help.txt
var longHelp string

var (
	actionsCmd = &cobra.Command{
		Use:   "actions",
		Short: "Show information about the actions",
		Long:  "'gh-not' has multiple actions that perform different actions:\n\n" + longHelp,
	}
)

func init() {
	rootCmd.AddCommand(actionsCmd)
}
