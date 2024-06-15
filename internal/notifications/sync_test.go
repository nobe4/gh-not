package notifications

import "testing"

func TestSync(t *testing.T) {
	n0 := &Notification{Id: "0"}
	n0Hidden := &Notification{Id: "0", Meta: Meta{Hidden: true}}
	n0Updated := &Notification{Id: "0"}
	n1 := &Notification{Id: "1"}
	n1ToDelete := &Notification{Id: "1", Meta: Meta{ToDelete: true}}

	tests := []struct {
		name     string
		local    Notifications
		remote   Notifications
		expected Notifications
	}{
		// (1)Insert
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

		// (3)Noop
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

		// (4)Drop
		{
			name:     "cleanup notifications",
			local:    Notifications{n0Hidden, n1ToDelete},
			expected: Notifications{n0Hidden},
		},

		// Testing (2)Update latest so previous tests are not impacted my
		// modified notifications.
		// (2)Update
		{
			name:     "hidden notification present in remote",
			local:    Notifications{n0Hidden, n1ToDelete},
			remote:   Notifications{n0, n1},
			expected: Notifications{n0Hidden, n1ToDelete},
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
				if n.Id != test.expected[i].Id {
					t.Fatalf("expected %+v but got %+v", test.expected[i], n)
				}
			}
		})
	}
}
