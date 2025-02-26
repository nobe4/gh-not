package config

import (
	"errors"
	"fmt"
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

	tests := []struct {
		envKey   string
		envValue string
		path     string
		want     string
	}{
		{
			envKey:   "HOME",
			envValue: "/home",
			path:     "/dev",
			want:     "/dev",
		},
		{
			envKey:   "HOME",
			envValue: "/home",
			path:     "$HOME/dev",
			want:     "/home/dev",
		},
		{
			envKey:   "XDG_CONFIG_HOME",
			envValue: "/home/.config",
			path:     "$XDG_CONFIG_HOME/dev",
			want:     "/home/.config/dev",
		},
		{
			envKey:   "FOO",
			envValue: "/home",
			path:     "$FOO/dev",
			want:     "/home/dev",
		},
		{
			envKey:   "FOO",
			envValue: "/user",
			path:     "$FOO/dev$FOO",
			want:     "/user/dev/user",
		},
		{
			envKey:   "FOO",
			envValue: "",
			path:     "$FOO/dev$FOO",
			want:     "/dev",
		},
		{
			envKey:   "FOO",
			envValue: "$FOO",
			path:     "/dev$FOO",
			want:     "/dev$FOO",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("expands '%s' correctly", test.path), func(t *testing.T) {
			t.Setenv(test.envKey, test.envValue)

			got, err := ExpandPathWithoutTilde(test.path)
			if err != nil {
				t.Fatalf("want no error, got %v", err)
			}

			if got != test.want {
				t.Errorf("want %q, got %q", test.want, got)
			}
		})
	}
}
