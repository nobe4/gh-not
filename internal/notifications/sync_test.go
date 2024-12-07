package notifications

import (
	"testing"
	"time"
)

func TestSync(t *testing.T) {
	n0 := &Notification{ID: "0", UpdatedAt: time.Unix(0, 1)}
	n0Hidden := &Notification{ID: "0", UpdatedAt: time.Unix(0, 1), Meta: Meta{Hidden: true}}
	n0Updated := &Notification{ID: "0", UpdatedAt: time.Unix(0, 1)}
	n1 := &Notification{ID: "1", UpdatedAt: time.Unix(0, 0)}
	n1Done := &Notification{ID: "1", Meta: Meta{Done: true}, UpdatedAt: time.Unix(0, 0)}

	tests := []struct {
		name     string
		local    Notifications
		remote   Notifications
		expected Notifications
	}{
		// (1) Insert
		{
			name:     "one new notification",
			remote:   Notifications{n0},
			expected: Notifications{n0},
		},
		{
			name:     "two new notifications",
			remote:   Notifications{n0, n1},
			expected: Notifications{n0, n1},
		},

		// (3) Keep
		{
			name: "no notifications",
		},
		{
			name:     "missing one notification",
			local:    Notifications{n0},
			expected: Notifications{n0},
		},
		{
			name:     "one new, one missing notification",
			local:    Notifications{n0},
			remote:   Notifications{n1},
			expected: Notifications{n0, n1},
		},
		{
			name:     "keep hidden notification",
			local:    Notifications{n0Hidden},
			remote:   Notifications{n0},
			expected: Notifications{n0Hidden},
		},

		// (4) Drop
		{
			name:     "cleanup notifications",
			local:    Notifications{n0Hidden, n1Done},
			expected: Notifications{},
		},

		// Testing (2) Update latest so previous tests are not impacted my
		// modified notifications.
		// (2) Update
		{
			name:     "hidden notification present in remote",
			local:    Notifications{n0Hidden, n1Done},
			remote:   Notifications{n0Updated, n1},
			expected: Notifications{n0Hidden, n1Done},
		},
		{
			name:     "updated notification present in remote",
			local:    Notifications{n0, n1},
			remote:   Notifications{n0Updated, n1},
			expected: Notifications{n0Updated, n1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Sync(test.local, test.remote)

			if len(got) != len(test.expected) {
				t.Fatalf("expected %d notifications but got %d", len(test.expected), len(got))
			}

			for i, n := range got {
				if n.ID != test.expected[i].ID {
					t.Fatalf("expected %+v but got %+v", test.expected[i], n)
				}
			}
		})
	}

	t.Run("update Meta.RemoteExist", func(t *testing.T) {
		got := Sync(
			Notifications{
				&Notification{ID: "0", UpdatedAt: time.Unix(0, 2), Meta: Meta{RemoteExists: false}},
				&Notification{ID: "1", UpdatedAt: time.Unix(0, 1), Meta: Meta{RemoteExists: true}},
			},
			Notifications{
				&Notification{ID: "0", UpdatedAt: time.Unix(0, 2)},
				&Notification{ID: "2", UpdatedAt: time.Unix(0, 0)},
			},
		)

		if !got[0].Meta.RemoteExists {
			t.Fatalf("expected RemoteExists to be true but got false")
		}
		if got[1].Meta.RemoteExists {
			t.Fatalf("expected RemoteExists to be false but got true")
		}
		if !got[2].Meta.RemoteExists {
			t.Fatalf("expected RemoteExists to be true but got false")
		}
	})

	t.Run("update Meta.Done", func(t *testing.T) {
		tests := []struct {
			name         string
			local        *Notification
			remote       *Notification
			expectedDone bool
		}{
			{
				name:         "Not done && Not updated",
				local:        &Notification{UpdatedAt: time.Unix(0, 0), Meta: Meta{Done: false}},
				remote:       &Notification{UpdatedAt: time.Unix(0, 0)},
				expectedDone: false,
			},
			{
				name:         "Not done && updated",
				local:        &Notification{UpdatedAt: time.Unix(0, 0), Meta: Meta{Done: false}},
				remote:       &Notification{UpdatedAt: time.Unix(0, 1)},
				expectedDone: false,
			},
			{
				name:         "Done && Not updated",
				local:        &Notification{UpdatedAt: time.Unix(0, 0), Meta: Meta{Done: true}},
				remote:       &Notification{UpdatedAt: time.Unix(0, 0)},
				expectedDone: true,
			},
			{
				name:         "Done && updated",
				local:        &Notification{UpdatedAt: time.Unix(0, 0), Meta: Meta{Done: true}},
				remote:       &Notification{UpdatedAt: time.Unix(0, 1)},
				expectedDone: false,
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got := Sync(Notifications{test.local}, Notifications{test.remote})
				if got[0].Meta.Done != test.expectedDone {
					t.Fatalf("expected Done to be %v but got %v", test.expectedDone, got[0].Meta.Done)
				}
			})
		}
	})
}
