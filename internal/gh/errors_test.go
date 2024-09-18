package gh

import "testing"

func TestRetryError(t *testing.T) {
	e := RetryError{
		verb: "verb",
		url:  "url",
	}
	want := "retry exceeded for verb url"
	got := e.Error()

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
