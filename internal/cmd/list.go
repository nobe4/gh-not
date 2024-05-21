package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List notifications",
		RunE: func(cmd *cobra.Command, args []string) error {
			notifications, err := client.Notifications()
			if err != nil {
				slog.Error("Failed to list notifications", "err", err)
				return err
			}

			out, err := notifications.ToTable()
			if err != nil {
				slog.Error("Failed to list the notifications", "err", err)
				return err
			}

			fmt.Printf(out)
			return nil
		},
	}
)
