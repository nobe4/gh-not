package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "repl",
		Short: "deprecated: use `gh-not --repl` instead",
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "deprecated: use 'gh-not' instead",
	})
}
