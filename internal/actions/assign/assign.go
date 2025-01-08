/*
Package assign implements an [actions.Runner] that assigns the subject of a
notification.

It only works when the notifications has an issue or pull request for subject.

It takes as arguments the usernames to assign.

Usage in the config:

	rules:
	  - action: assign
	    args: [user0, user1]

Usage in the REPL:

	:assign user0 user1

Refs: https://docs.github.com/en/rest/issues/assignees?apiVersion=2022-11-28#add-assignees-to-an-issue
*/
package assign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct {
	Client *gh.Client
}

type Body struct {
	Assignees []string `json:"assignees"`
}

func (a *Runner) Run(n *notifications.Notification, assignees []string, w io.Writer) error {
	slog.Debug("assigning notification", "notification", n, "assignees", assignees)

	if len(assignees) == 0 {
		return errors.New("no assignees provided")
	}

	url, ok := issueURL(n.Subject.URL)
	if !ok {
		slog.Warn("not an issue or pull", "notification", n)
		return nil
	}

	assigneesURL := url + "/assignees"

	body, err := json.Marshal(Body{Assignees: assignees})
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	r, err := a.Client.API.Request(http.MethodPost, assigneesURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to request assignees: %w", err)
	}
	defer r.Body.Close()

	fmt.Fprint(w, colors.Red("ASSIGN ")+n.String()+" to "+strings.Join(assignees, ", "))

	return nil
}

func issueURL(url string) (string, bool) {
	re := regexp.MustCompile(`^(https://api\.github\.com/repos/.+/.+/)(issues|pulls)(/\d+)$`)
	matches := re.FindStringSubmatch(url)

	if len(matches) == 0 {
		return "", false
	}

	return fmt.Sprintf("%sissues%s", matches[1], matches[3]), true
}
