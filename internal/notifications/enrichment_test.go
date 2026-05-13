package notifications

import (
	"testing"
	"time"
)

func TestMergeUpdatedNotification(t *testing.T) {
	t.Parallel()

	t0 := time.Unix(0, 0)
	t1 := time.Unix(0, 1)

	enrichedAuthor := User{Login: "author", Type: "User"}
	enrichedCommentor := User{Login: "commentor", Type: "User"}
	enrichedAssignees := []User{{Login: "assignee", Type: "User"}}
	enrichedReviewers := []User{{Login: "reviewer", Type: "User"}}
	enrichedTeams := []Team{{Name: "team", ID: 1}}
	enrichedMergedBy := User{Login: "merger", Type: "User"}
	enrichedState := "closed"
	enrichedHTMLURL := "https://github.com/test/1"

	tests := []struct {
		name string
		n    *Notification
		o    *Notification
		want *Notification
	}{
		{
			name: "remote is newer, clears enrichment and resets meta",
			n: &Notification{
				UpdatedAt:       t0,
				Author:          enrichedAuthor,
				LatestCommentor: enrichedCommentor,
				Meta: Meta{
					Done:     true,
					Enriched: true,
				},
			},
			o: &Notification{
				UpdatedAt:       t1,
				Author:          User{Login: "old", Type: "User"},
				LatestCommentor: User{Login: "old", Type: "User"},
			},
			want: &Notification{
				UpdatedAt: t1,
				Meta: Meta{
					Done:         false,
					Enriched:     false,
					RemoteExists: true,
				},
			},
		},
		{
			name: "remote is same age, enriched, copies enrichment",
			n: &Notification{
				UpdatedAt:       t0,
				Author:          enrichedAuthor,
				LatestCommentor: enrichedCommentor,
				Assignees:       enrichedAssignees,
				Reviewers:       enrichedReviewers,
				ReviewersTeams:  enrichedTeams,
				MergedBy:        enrichedMergedBy,
				Subject:         Subject{State: enrichedState, HTMLURL: enrichedHTMLURL},
				Meta: Meta{
					Enriched: true,
				},
			},
			o: &Notification{
				UpdatedAt: t0,
			},
			want: &Notification{
				UpdatedAt:       t0,
				Author:          enrichedAuthor,
				LatestCommentor: enrichedCommentor,
				Assignees:       enrichedAssignees,
				Reviewers:       enrichedReviewers,
				ReviewersTeams:  enrichedTeams,
				MergedBy:        enrichedMergedBy,
				Subject:         Subject{State: enrichedState, HTMLURL: enrichedHTMLURL},
				Meta: Meta{
					Enriched:     true,
					RemoteExists: true,
				},
			},
		},
		{
			name: "remote is older, enriched, copies enrichment",
			n: &Notification{
				UpdatedAt:       t1,
				Author:          enrichedAuthor,
				LatestCommentor: enrichedCommentor,
				Assignees:       enrichedAssignees,
				Reviewers:       enrichedReviewers,
				ReviewersTeams:  enrichedTeams,
				MergedBy:        enrichedMergedBy,
				Subject:         Subject{State: enrichedState, HTMLURL: enrichedHTMLURL},
				Meta: Meta{
					Enriched: true,
				},
			},
			o: &Notification{
				UpdatedAt: t0,
			},
			want: &Notification{
				UpdatedAt:       t0,
				Author:          enrichedAuthor,
				LatestCommentor: enrichedCommentor,
				Assignees:       enrichedAssignees,
				Reviewers:       enrichedReviewers,
				ReviewersTeams:  enrichedTeams,
				MergedBy:        enrichedMergedBy,
				Subject:         Subject{State: enrichedState, HTMLURL: enrichedHTMLURL},
				Meta: Meta{
					Enriched:     true,
					RemoteExists: true,
				},
			},
		},
		{
			name: "remote is same age, not enriched, no enrichment copy",
			n: &Notification{
				UpdatedAt: t0,
				Meta: Meta{
					Enriched: false,
					Done:     true,
				},
			},
			o: &Notification{
				UpdatedAt: t0,
			},
			want: &Notification{
				UpdatedAt: t0,
				Meta: Meta{
					Done:         true,
					Enriched:     false,
					RemoteExists: true,
				},
			},
		},
		{
			name: "meta fields preserved through merge",
			n: &Notification{
				UpdatedAt: t0,
				Meta: Meta{
					Hidden: true,
					Done:   true,
				},
			},
			o: &Notification{
				UpdatedAt: t0,
			},
			want: &Notification{
				UpdatedAt: t0,
				Meta: Meta{
					Hidden:       true,
					Done:         true,
					RemoteExists: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.n.Update(tt.o)

			if got != tt.o {
				t.Fatal("expected returned notification to be o")
			}

			if !got.Equal(tt.want) {
				t.Fatalf("got %#v, want %#v", *got, *tt.want)
			}
		})
	}
}
