package notifications

import "testing"

func TestVisible(t *testing.T) {
	t.Parallel()

	n := Notifications{
		&Notification{Meta: Meta{Done: false, Hidden: false}},
		&Notification{Meta: Meta{Done: true, Hidden: false}},
		&Notification{Meta: Meta{Done: false, Hidden: true}},
		&Notification{Meta: Meta{Done: true, Hidden: true}},
	}

	visible := n.Visible()

	if len(visible) != 1 {
		t.Errorf("Expected 1, got %d", len(visible))
	}

	if visible[0] != n[0] {
		t.Errorf("Expected %v, got %v", n[0], visible[0])
	}
}
