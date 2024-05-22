package cmd

import (
	"fmt"
	"log/slog"

	"github.com/nobe4/gh-not/internal/jq"
	"github.com/spf13/cobra"
)

// TODO: move this into the list command with the flag --jq
var (
	jqCmd = &cobra.Command{
		Use:   "jq",
		Short: "Filter notifications using jq",
		RunE: func(cmd *cobra.Command, args []string) error {
			notifications, err := client.Notifications()
			if err != nil {
				slog.Error("Failed to list notifications", "err", err)
				return err
			}

			notificationsList, err := jq.Filter(args[0], notifications.ToSlice())
			if err != nil {
				return err
			}

			out, err := notificationsList.ToMap().ToTable()
			if err != nil {
				return err
			}
			fmt.Println(out)

			return nil
		},
	}
)
