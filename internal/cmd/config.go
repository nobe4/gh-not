package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	configPkg "github.com/nobe4/gh-not/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	editConfigFlag = false
	initConfigFlag = false

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Show configuration information",
		RunE:  runConfig,
	}
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&editConfigFlag, "edit", "e", false, "Edit the config in $EDITOR")
	configCmd.Flags().BoolVarP(&initConfigFlag, "init", "i", false, "Create an initial config file")
}

func runConfig(cmd *cobra.Command, args []string) error {
	if initConfigFlag {
		return initConfig()
	}

	if editConfigFlag {
		return editConfig()
	}

	marshalled, err := yaml.Marshal(config)
	if err != nil {
		slog.Error("Failed to marshall config", "err", err)
	}

	fmt.Println(configPathFlag)
	fmt.Println(string(marshalled))

	return nil
}

func initConfig() error {
	slog.Debug("creating initial config file", "path", configPathFlag)

	if _, err := os.Stat(configPathFlag); err == nil {
		return fmt.Errorf("config file %s already exists", configPathFlag)
	}

	f, err := os.Create(configPathFlag)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(configPkg.Example)); err != nil {
		return err
	}

	fmt.Printf("Created config file: %s\n", configPathFlag)

	return nil
}

func editConfig() error {
	slog.Debug("editing config file", "path", configPathFlag)

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
