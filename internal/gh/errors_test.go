package gh

import "testing"

func TestRetryError(t *testing.T) {
	e := RetryError{
		verb: "verb",
		url:  "endpoint",
	}
	want := "retry exceeded for verb endpoint"
	got := e.Error()

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
