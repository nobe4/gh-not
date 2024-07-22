package notifications

import (
	"testing"
	"time"
)

func TestIDList(t *testing.T) {
	n := Notifications{}
	ids := []string{"0", "1", "2"}
	for _, id := range ids {
		n = append(n, &Notification{Id: id})
	}

	got := n.IDList()

	for i, id := range ids {
		if id != got[i] {
			t.Fatalf("expected %d: %s but got %s", i, id, got[i])
		}
	}
}

func TestCompact(t *testing.T) {
	n0 := &Notification{Id: "0"}
	n1 := &Notification{Id: "1"}
	n := Notifications{nil, nil, n0, nil, n1, nil}

	got := n.Compact()

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

func TestSort(t *testing.T) {
	n0 := &Notification{UpdatedAt: time.Unix(0, 2)}
	n1 := &Notification{UpdatedAt: time.Unix(0, 1)}
	n2 := &Notification{UpdatedAt: time.Unix(0, 0)}
	n := Notifications{n2, n0, n1}

	n.Sort()

	if n[0] != n0 {
		t.Fatalf("expected %#v but got %#v", n0, n[0])
	}

	if n[1] != n1 {
		t.Fatalf("expected %#v but got %#v", n1, n[1])
	}

	if n[2] != n2 {
		t.Fatalf("expected %#v but got %#v", n2, n[2])
	}
}

func TestMap(t *testing.T) {
	n0 := &Notification{Id: "0", UpdatedAt: time.Unix(0, 2)}
	n1 := &Notification{Id: "1", UpdatedAt: time.Unix(0, 1)}
	n2 := &Notification{Id: "2", UpdatedAt: time.Unix(0, 0)}
	n := Notifications{n0, n1, n2}

	got := n.Map()

	if len(got) != 3 {
		t.Fatalf("expected 3 elements but got %d\n%+v", len(got), got)
	}

	// Testing map is challenging because the order is not guaranteed.
	// Instead, we will compare the sorted list.
	l := got.List()
	l.Sort()

	for i, l := range l {
		if n[i] != l {
			t.Fatalf("expected %+v but got %+v", n[i], got)
		}
	}
}

func TestList(t *testing.T) {
	n0 := &Notification{Id: "0"}
	n1 := &Notification{Id: "1"}
	n2 := &Notification{Id: "2"}
	n := NotificationMap{
		"0": n0,
		"1": n1,
		"2": n2,
	}

	got := n.List()

	if len(got) != 3 {
		t.Fatalf("expected 3 elements but got %d\n%+v", len(got), got)
	}

	for _, got := range got {
		if got != n0 && got != n1 && got != n2 {
			t.Fatalf("expected %+v to be in the list", got)
		}
	}
}

func TestUniq(t *testing.T) {
	n0 := &Notification{Id: "0"}
	n1 := &Notification{Id: "1"}
	n2 := &Notification{Id: "2"}
	n := Notifications{n0, n1, n2, n0, n1, n2}

	got := n.Uniq()

	if len(got) != 3 {
		t.Fatalf("expected 3 elements but got %d\n%+v", len(got), got)
	}

	if got[0] != n0 {
		t.Fatalf("expected %+v but got %+v", n0, got[0])
	}

	if got[1] != n1 {
		t.Fatalf("expected %+v but got %+v", n1, got[1])
	}

	if got[2] != n2 {
		t.Fatalf("expected %+v but got %+v", n2, got[2])
	}
}

func TestFilterFromIds(t *testing.T) {
	n0 := &Notification{Id: "0"}
	n1 := &Notification{Id: "1"}
	n2 := &Notification{Id: "2"}
	n := Notifications{n0, n1, n2}

	got := n.FilterFromIds([]string{"0", "2"})

	if len(got) != 2 {
		t.Fatalf("expected 2 elements but got %d\n%+v", len(got), got)
	}

	if got[0] != n0 {
		t.Fatalf("expected %+v but got %+v", n0, got[0])
	}

	if got[1] != n2 {
		t.Fatalf("expected %+v but got %+v", n2, got[1])
	}
}
