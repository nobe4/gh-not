package config

import (
	"slices"
	"testing"

	"github.com/nobe4/gh-not/internal/notifications"
)

func TestFilterIds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		r    Rule
		n    notifications.Notifications
		want []string
	}{
		{
			r:    Rule{Name: "no filter"},
			n:    notifications.Notifications{{ID: "0"}, {ID: "1"}, {ID: "2"}},
			want: []string{"0", "1", "2"},
		},

		{
			r: Rule{
				Name:    "filter for a specific id",
				Filters: []string{`.id == "1"`},
			},
			n:    notifications.Notifications{{ID: "0"}, {ID: "1"}, {ID: "2"}},
			want: []string{"1"},
		},

		// Or
		{
			r: Rule{
				Name: "filter for a specific set of ids",
				Filters: []string{
					// All three filters are equivalent
					`.id == "1" or .id == "2"`,
					`(.id == "1" or .id == "2")`,
					`(.id == "1") or (.id == "2")`,
				},
			},
			n:    notifications.Notifications{{ID: "0"}, {ID: "1"}, {ID: "2"}},
			want: []string{"1", "2"},
		},

		// And
		{
			r: Rule{
				Name: "filters are joined by and",
				Filters: []string{
					`.reason == "test"`,
					`.unread == true`,
				},
			},
			n: notifications.Notifications{
				{ID: "0", Reason: "test", Unread: true},
				{ID: "1", Reason: "test", Unread: false},
				{ID: "2", Reason: "ci_activity", Unread: true},
			},
			want: []string{"0"},
		},
		// Can also be written as
		{
			r: Rule{
				Name: "use and in filter",
				Filters: []string{
					`.reason == "test" and .unread == true`,
				},
			},
			n: notifications.Notifications{
				{ID: "0", Reason: "test", Unread: true},
				{ID: "1", Reason: "test", Unread: false},
				{ID: "2", Reason: "ci_activity", Unread: true},
			},
			want: []string{"0"},
		},

		// Order of operations
		{
			r: Rule{
				Name: "and is evaluated before or",
				Filters: []string{
					`.reason == "test" and .unread == true or .reason == "test2"`,
				},
			},
			n: notifications.Notifications{
				{ID: "0", Reason: "test", Unread: true},
				{ID: "1", Reason: "test", Unread: false},
				{ID: "2", Reason: "test2", Unread: true},
			},
			want: []string{"0", "2"},
		},
		{
			r: Rule{
				Name: "parenthesis can force the order of operations",
				Filters: []string{
					`.reason == "test" and (.unread == true or .id == "1")`,
				},
			},
			n: notifications.Notifications{
				{ID: "0", Reason: "test", Unread: true},
				{ID: "1", Reason: "test", Unread: false},
				{ID: "2", Reason: "test2", Unread: true},
			},
			want: []string{"0", "1"},
		},
		{
			r: Rule{
				Name: "parenthesis work also across filters",
				Filters: []string{
					`(.reason == "test" or .id == "2")`,
					`(.unread == true or .id == "1")`,
				},
			},
			n: notifications.Notifications{
				{ID: "0", Reason: "test", Unread: true},
				{ID: "1", Reason: "test", Unread: false},
				{ID: "2", Reason: "test2", Unread: true},
				{ID: "3", Reason: "test2", Unread: false},
			},
			want: []string{"0", "1", "2"},
		},

		// Issue-related tests
		{
			r: Rule{
				Name: "https://github.com/nobe4/gh-not/issues/86",
				Filters: []string{
					`(.repository.full_name == "org/repo1" or .repository.full_name == "org/repo2")`,
					`.reason == "review_requested"`,
				},
			},
			n: notifications.Notifications{
				{ID: "0", Repository: notifications.Repository{FullName: "org/repo1"}, Reason: "review_requested"},
				{ID: "1", Repository: notifications.Repository{FullName: "org/repo1"}, Reason: "test"},
				{ID: "2", Repository: notifications.Repository{FullName: "org/repo2"}, Reason: "review_requested"},
				{ID: "3", Repository: notifications.Repository{FullName: "org/repo2"}, Reason: "test"},
				{ID: "4", Repository: notifications.Repository{FullName: "org/repo3"}, Reason: "review_requested"},
				{ID: "5", Repository: notifications.Repository{FullName: "org/repo3"}, Reason: "test"},
			},
			want: []string{"0", "2"},
		},
		{
			r: Rule{
				Name: "https://github.com/nobe4/gh-not/issues/86",
				Filters: []string{
					`(.repository.full_name == "org/repo1" or .repository.full_name == "org/repo2") and .reason == "review_requested"`,
				},
			},
			n: notifications.Notifications{
				{ID: "0", Repository: notifications.Repository{FullName: "org/repo1"}, Reason: "review_requested"},
				{ID: "1", Repository: notifications.Repository{FullName: "org/repo1"}, Reason: "test"},
				{ID: "2", Repository: notifications.Repository{FullName: "org/repo2"}, Reason: "review_requested"},
				{ID: "3", Repository: notifications.Repository{FullName: "org/repo2"}, Reason: "test"},
				{ID: "4", Repository: notifications.Repository{FullName: "org/repo3"}, Reason: "review_requested"},
				{ID: "5", Repository: notifications.Repository{FullName: "org/repo3"}, Reason: "test"},
			},
			want: []string{"0", "2"},
		},
	}

	for _, test := range tests {
		t.Run(test.r.Name, func(t *testing.T) {
			t.Parallel()

			got, err := test.r.Filter(test.n)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(got) != len(test.want) {
				t.Fatalf("want %#v, but got %#v", test.want, got)
			}

			if !slices.Equal(got.IDList(), test.want) {
				t.Fatalf("want %#v, but got %#v", test.want, got)
			}
		})
	}
}

func TestValidationEmptyAction(t *testing.T) {
	r := Rule{
		Name: "parenthesis work also across filters",
		Filters: []string{
			`(.reason == "test" or .id == "2")`,
			`(.unread == true or .id == "1")`,
		},
	}
	err := r.Validate()
	if err == nil {
		t.Fatal("rule with no action passed validation (it shouldn't)")
	}
}

func TestValidationEmptyFilters(t *testing.T) {
	r := Rule{
		Name:    "parenthesis work also across filters",
		Filters: []string{},
		Action:  "done",
	}
	if r.Validate() == nil {
		t.Fatal("rule with no filter passed validation (it shouldn't)")
	}
}

func TestValidationPass(t *testing.T) {
	r := Rule{
		Name: "parenthesis work also across filters",
		Filters: []string{
			`(.reason == "test" or .id == "2")`,
			`(.unread == true or .id == "1")`,
		},
		Action: "done",
	}
	if r.Validate() != nil {
		t.Fatal("rule that should pass validation didn't")
	}
}

func TestValidationPassEvenWithNoName(t *testing.T) {
	r := Rule{
		Filters: []string{
			`(.reason == "test" or .id == "2")`,
			`(.unread == true or .id == "1")`,
		},
		Action: "done",
	}
	if r.Validate() != nil {
		t.Fatal("rule that should pass validation didn't")
	}
}
