package config

import (
	"errors"
	"testing"
)

//nolint:tparallel // t.Setenv forbits running tests in parallel
func TestExpandPathWithoutTilde(t *testing.T) {
	t.Run("rejects a path with a tildle", func(t *testing.T) {
		t.Parallel()

		path := "~/.config/gh-not"

		got, err := ExpandPathWithoutTilde(path)
		if !errors.Is(err, errTildeUsage) {
			t.Errorf("want %v, got %v", errTildeUsage, err)
		}

		if got != "" {
			t.Errorf("want empty string, got %q", got)
		}
	})

	t.Run("expands environment variables", func(t *testing.T) {
		t.Setenv("HOME", "/home/user")

		path := "$HOME/.config/gh-not"
		want := "/home/user/.config/gh-not"

		got, err := ExpandPathWithoutTilde(path)
		if err != nil {
			t.Fatalf("want no error, got %v", err)
		}

		if got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	})
}
