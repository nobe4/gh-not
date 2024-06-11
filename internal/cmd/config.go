package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Show configuration information",
		RunE:  runConfig,
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	marshalled, err := yaml.Marshal(config)
	if err != nil {
		slog.Error("Failed to marshall config", "err", err)
	}

	fmt.Println(configPathFlag)
	fmt.Println(string(marshalled))

	return nil
}
