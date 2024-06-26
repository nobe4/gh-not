package jq

import (
	"testing"

	"github.com/nobe4/gh-not/internal/notifications"
)

func notificationsEqual(a notifications.Notifications, ids []string) bool {
	if len(a) != len(ids) {
		return false
	}

	for i := range a {
		if a[i].Id != ids[i] {
			return false
		}
	}

	return true
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter string
		n      notifications.Notifications
		want   []string
		err    string
	}{
		{
			name:   "empty filter",
			filter: "",
			n: notifications.Notifications{
				&notifications.Notification{Id: "0"},
				&notifications.Notification{Id: "1"},
				&notifications.Notification{Id: "2"},
			},
			want: []string{"0", "1", "2"},
		},
		{
			name:   "invalid filter",
			filter: "!!!",
			err:    `unexpected token "!"`,
		},
		{
			name:   "filter on specific id",
			filter: `.id == "1"`,
			n: notifications.Notifications{
				&notifications.Notification{Id: "0"},
				&notifications.Notification{Id: "1"},
				&notifications.Notification{Id: "2"},
			},
			want: []string{"1"},
		},
		{
			name:   "composite filter",
			filter: `(.id == "1" or .id == "2") and (.unread == true)`,
			n: notifications.Notifications{
				&notifications.Notification{Id: "0"},
				&notifications.Notification{Id: "1", Unread: true},
				&notifications.Notification{Id: "2", Unread: false},
			},
			want: []string{"1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Filter(test.filter, test.n)
			if err != nil && err.Error() != test.err {
				t.Fatalf("expected error %s but got %s", test.err, err.Error())
			}

			if !notificationsEqual(got, test.want) {
				t.Fatalf("expected %#v but got %#v", test.want, got.IDList())
			}
		})
	}
}
