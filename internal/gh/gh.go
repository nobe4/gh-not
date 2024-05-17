package gh

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/nobe4/gh-not/internal/notifications"
)

func Run(args []string) ([]notifications.Notification, error) {
	// cmd := exec.Command("gh", args...)
	fmt.Println("mocking gh ", args)
	cmd := exec.Command("cat", "notifications.test")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	var allNotifications []notifications.Notification
	if err := json.NewDecoder(stdout).Decode(&allNotifications); err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	return allNotifications, nil
}
