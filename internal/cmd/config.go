package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	configPkg "github.com/nobe4/gh-not/internal/config"
	"github.com/spf13/cobra"
)

var (
	editConfigFlag = false
	initConfigFlag = false

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Print the config to stdout",
		RunE:  runConfig,
	}
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&editConfigFlag, "edit", "e", false, "Edit the config in $EDITOR")
	configCmd.Flags().BoolVarP(&initConfigFlag, "init", "i", false, "Create the default config file")
}

func runConfig(cmd *cobra.Command, args []string) error {
	if initConfigFlag {
		if err := initConfig(); err != nil {
			return err
		}
	}

	if editConfigFlag {
		return editConfig()
	}

	marshalled, err := config.Marshal()
	if err != nil {
		return err
	}
	fmt.Printf("Config sourced from: %s\n\n%s\n", config.Path, marshalled)

	return nil
}

func initConfig() error {
	slog.Debug("creating initial config file", "path", configPathFlag)

	if err := configPkg.Default(configPathFlag).WriteConfig(); err != nil {
		slog.Error("Failed to save initial config", "err", err)
		return err
	}
	fmt.Printf("Initial config saved to %s\n", configPathFlag)

	return nil
}

func editConfig() error {
	slog.Debug("editing config file", "path", configPathFlag)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("EDITOR environment variable not set")
	}

	cmd := exec.Command(editor, config.Path)

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
