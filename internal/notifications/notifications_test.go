package notifications

import "testing"

func TestDeleteNil(t *testing.T) {
	n0 := &Notification{Id: "0"}
	n1 := &Notification{Id: "1"}
	n := Notifications{nil, nil, n0, nil, n1, nil}

	got := n.DeleteNil()

	if len(got) != 2 {
		t.Fatalf("expected 2 elements but got %d\n%+v", len(got), got)
	}

    if got[0] != n0 {
		t.Fatalf("expected %+v but got %+v", n0, got[0])
    }

    if got[1] != n1 {
		t.Fatalf("expected %+v but got %+v", n1, got[1])
    }
}
