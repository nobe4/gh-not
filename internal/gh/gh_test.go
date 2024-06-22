package gh

import (
	"io"
	"maps"
	"strings"
	"testing"
)

func TestIsRetryable(t *testing.T) { t.Skip("TODO") }
func TestRetry(t *testing.T)       { t.Skip("TODO") }

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		body io.ReadCloser
		// Using a map[string]string type for the expected value to ease the comparison.
		// This does not matter for the tests' validity
		expected map[string]string
		fails    bool
	}{
		{
			name: "empty body",
			body: io.NopCloser(strings.NewReader(`{}`)),
		},
		{
			name:     "map body",
			body:     io.NopCloser(strings.NewReader(`{"a": "b"}`)),
			expected: map[string]string{"a": "b"},
		},
		{
			name:  "invalid body",
			body:  io.NopCloser(strings.NewReader(`{`)),
			fails: true,
		},
		{
			name:  "invalid JSON",
			body:  io.NopCloser(strings.NewReader(`1`)),
			fails: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := map[string]string{}
			err := decode(test.body, &out)

			if test.fails && err == nil {
				t.Errorf("expected test to fails")
			} else if !test.fails && err != nil {
				t.Errorf("expected test to pass, got %v", err)
			}

			if !maps.Equal(out, test.expected) {
				t.Errorf("expected %#v, got %#v", test.expected, out)
			}
		})
	}
}
