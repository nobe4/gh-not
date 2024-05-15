package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/itchyny/gojq"
)

type NotificationResponse []Notification

type Notification struct {
	Title      string     `json:"title"`
	Id         string     `json:"id"`
	Unread     bool       `json:"unread"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`

	Owner Owner `json:"owner"`
}

type Owner struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type Config struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Name    string   `json:"name"`
	Filters []string `json:"filters"`
	Action  string   `json:"action"`
}

func main() {
	notifications, err := gh([]string{"api", "/notifications?all=true&per_page=3"})

	if err != nil {
		panic(err)
	}

	filteredNotifications := []Notification{}

	config, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}
	for _, group := range config.Groups {
		fmt.Println("Group: ", group.Name)
		for _, filter := range group.Filters {
			fmt.Println("Filter: ", filter)
			selectedNotifications, err := jq(filter, notifications)
			if err != nil {
				panic(err)
			}
			fmt.Println("Filtered: ", len(selectedNotifications))

			filteredNotifications = append(filteredNotifications, selectedNotifications...)
		}
		filteredNotifications = uniq(filteredNotifications)

		for _, notification := range filteredNotifications {
			fmt.Printf("%s %v\n", group.Action, notification)
		}
	}

}

func loadConfig(path string) (Config, error) {
	var config Config

	content, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(content, &config); err != nil {
		return config, err
	}

	return config, nil
}

func uniq(notifications []Notification) []Notification {
	seen := make(map[string]bool)
	unique := []Notification{}

	for _, notification := range notifications {
		if _, ok := seen[notification.Id]; !ok {
			seen[notification.Id] = true
			unique = append(unique, notification)
		}
	}

	return unique
}

func gh(args []string) ([]Notification, error) {
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

	var notifications []Notification
	if err := json.NewDecoder(stdout).Decode(&notifications); err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	return notifications, nil
}

func jq(q string, notifications []Notification) ([]Notification, error) {
	query, err := gojq.Parse(fmt.Sprintf(".[] | select(%s)", q))
	if err != nil {
		panic(err)
	}

	// gojq works only on any data, so we need to convert []Notifications to
	// interface{}
	// This also gives us back the JSON fields from the API.
	marshalled, err := json.Marshal(notifications)
	if err != nil {
		panic(err)
	}

	var notificationsRaw interface{}
	if err := json.Unmarshal(marshalled, &notificationsRaw); err != nil {
		panic(err)
	}

	fitleredNotificationsRaw := []interface{}{}
	iter := query.Run(notificationsRaw)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			panic(err)
		}

		fitleredNotificationsRaw = append(fitleredNotificationsRaw, v)
	}

	marshalled, err = json.Marshal(fitleredNotificationsRaw)
	if err != nil {
		panic(err)
	}

	var filteredNotifications []Notification
	if err := json.Unmarshal(marshalled, &filteredNotifications); err != nil {
		panic(err)
	}

	return filteredNotifications, nil
}
