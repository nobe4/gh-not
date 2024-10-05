/*
Package tag implements an [actions.Runner] that manages tags in a notification.

Tags are stored sorted and unique.

It takes as arguments the tags to add prefixed by `+` and the tags to remove
prefixed by `-`. If a tag is not prefixed, it is added.

E.g.: `tag0 -tag1 +tag2` will add `tag0` and `tag2` and remove `tag1`.

Adding an existing tag is a no-op.
Removing a missing tag is a no-op.

Usage in the config:

	rules:
	  - action: tag
	    args: ["+tag0", "-tag1"]

Usage in the REPL:

	:tag +tag0 -tag1
*/
package tag

import (
	"fmt"
	"io"
	"log/slog"
	"slices"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/notifications"
)

type Runner struct{}

func (_ *Runner) Run(n *notifications.Notification, tags []string, w io.Writer) error {
	slog.Debug("tagging notification", "notification", n.Id, "tags", tags)

	tagsToAdd := []string{}
	tagsToRemove := []string{}

	for _, tag := range tags {
		switch tag[0] {
		case '+':
			tagsToAdd = append(tagsToAdd, tag[1:])
		case '-':
			tagsToRemove = append(tagsToRemove, tag[1:])
		default:
			tagsToAdd = append(tagsToAdd, tag)

		}
	}

	for _, tag := range tagsToAdd {
		n.Meta.Tags = append(n.Meta.Tags, tag)
	}

	newTags := []string{}
	for _, tag := range n.Meta.Tags {
		if !slices.Contains(tagsToRemove, tag) {
			newTags = append(newTags, tag)
		}
	}

	// Ensure tags are sorted and unique
	slices.Sort(newTags)
	n.Meta.Tags = slices.Compact(newTags)

	fmt.Fprint(w, colors.Red("TAGGED ")+n.String())

	return nil
}
