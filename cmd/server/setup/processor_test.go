package setup

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolveValue_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		isFlagChanged bool
		flagValue     string
		envValue      string
		want          string
	}{
		{"changed flag wins over env", true, "from-flag", "from-env", "from-flag"},
		{"env wins over unchanged flag default", false, "default", "from-env", "from-env"},
		{"default when env is empty", false, "default", "", "default"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := resolveValue(tc.isFlagChanged, tc.flagValue, tc.envValue)

			require.Equal(t, tc.want, got)
		})
	}
}

func Test_ResolveBool_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		isFlagChanged bool
		flagValue     bool
		envValue      string
		want          bool
	}{
		{"changed flag wins over env", true, false, "1", false},
		{"env wins over unchanged flag default", false, false, "1", true},
		{"invalid env falls back to flag default", false, false, "yes", false},
		{"default when env is empty", false, false, "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := resolveBool(tc.isFlagChanged, tc.flagValue, tc.envValue)

			require.Equal(t, tc.want, got)
		})
	}
}
