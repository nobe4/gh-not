package notifications

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

func Uniq(notificationsToFilter []Notification) []Notification {
	seen := make(map[string]bool)
	unique := []Notification{}

	for _, notification := range notificationsToFilter {
		if _, ok := seen[notification.Id]; !ok {
			seen[notification.Id] = true
			unique = append(unique, notification)
		}
	}

	return unique
}
