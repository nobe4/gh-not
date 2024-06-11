package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	editConfigFlag = false

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Show configuration information",
		RunE:  runConfig,
	}
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&editConfigFlag, "edit", "e", false, "Edit the config in $EDITOR")
}

func runConfig(cmd *cobra.Command, args []string) error {
	if editConfigFlag {
		return edit()
	}

	marshalled, err := yaml.Marshal(config)
	if err != nil {
		slog.Error("Failed to marshall config", "err", err)
	}

	fmt.Println(configPathFlag)
	fmt.Println(string(marshalled))

	return nil
}

func edit() error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("EDITOR environment variable not set")
	}

	cmd := exec.Command(editor, configPathFlag)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
