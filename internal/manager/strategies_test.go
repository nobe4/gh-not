package manager

import "testing"

func TestRefreshStrategy(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		tests := []struct {
			strategy RefreshStrategy
			want     string
		}{
			{AutoRefresh, "auto"},
			{ForceRefresh, "force"},
			{PreventRefresh, "prevent"},
			{RefreshStrategy(-1), "unknown"},
		}

		for _, test := range tests {
			if got := test.strategy.String(); got != test.want {
				t.Errorf("RefreshStrategy.String() = %q, want %q", got, test.want)
			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		tests := []struct {
			value string
			want  RefreshStrategy
			error bool
		}{
			{"auto", AutoRefresh, false},
			{"force", ForceRefresh, false},
			{"prevent", PreventRefresh, false},
			{"test", 0, true},
		}

		for _, test := range tests {
			t.Run(test.value, func(t *testing.T) {
				var got RefreshStrategy
				error := got.Set(test.value)

				if test.error {
					if error == nil {
						t.Errorf("expected an error but got none")
					}
				} else {
					if error != nil {
						t.Errorf("expected no error but got %#v", error)
					}

					if got != test.want {
						t.Errorf("RefreshStrategy.Set(%s) = %q, want %q", test.value, got, test.want)
					}
				}
			})
		}
	})
}

func TestForceStrategy(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		tests := []struct {
			strategy ForceStrategy
			want     string
		}{
			{ForceStrategy(0), ""},
			{ForceApply, "apply"},
			{ForceEnrich, "enrich"},
			{ForceApply | ForceEnrich, "apply, enrich"},
		}

		for _, test := range tests {
			if got := test.strategy.String(); got != test.want {
				t.Errorf("ForceStrategy.String() = %q, want %q", got, test.want)
			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		tests := []struct {
			value string
			want  ForceStrategy
			error bool
		}{
			{"apply", ForceApply, false},
			{"enrich", ForceEnrich, false},
			{"apply,enrich", ForceApply | ForceEnrich, false},
			{"test", 0, true},
		}

		for _, test := range tests {
			t.Run(test.value, func(t *testing.T) {
				var got ForceStrategy
				error := got.Set(test.value)

				if test.error {
					if error == nil {
						t.Errorf("expected an error but got none")
					}
				} else {
					if error != nil {
						t.Errorf("expected no error but got %#v", error)
					}

					if got != test.want {
						t.Errorf("RefreshStrategy.Set(%s) = %q, want %q", test.value, got, test.want)
					}
				}
			})
		}
	})
}
