package tag

import (
	"slices"
	"strings"
	"testing"

	"github.com/nobe4/gh-not/internal/notifications"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tags []string
		args []string
		want []string
	}{
		{
			name: "empty",
		},
		{
			name: "no new tag",
			tags: []string{"tag0", "tag1"},
			want: []string{"tag0", "tag1"},
		},
		{
			name: "add tags",
			tags: []string{"tag0", "tag1"},
			args: []string{"+tag2", "tag3"},
			want: []string{"tag0", "tag1", "tag2", "tag3"},
		},
		{
			name: "remove tags",
			tags: []string{"tag0", "tag1"},
			args: []string{"-tag1", "-tag2"},
			want: []string{"tag0"},
		},
		{
			name: "add and remove tags",
			tags: []string{"tag0", "tag1", "tag2"},
			args: []string{"-tag1", "tag3", "tag2"},
			want: []string{"tag0", "tag2", "tag3"},
		},
		{
			name: "handle duplicates",
			tags: []string{"tag0", "tag1", "tag2"},
			args: []string{"+tag1", "+tag1", "+tag2", "-tag2", "-tag2", "tag3", "-tag3"},
			want: []string{"tag0", "tag1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			r := &Runner{}
			n := &notifications.Notification{
				Meta: notifications.Meta{Tags: test.tags},
			}
			w := &strings.Builder{}

			if err := r.Run(n, test.args, w); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !slices.Equal(n.Meta.Tags, test.want) {
				t.Errorf("expected tags to be %v, got %v", test.want, n.Meta.Tags)
			}
		})
	}
}
