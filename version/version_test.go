package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolveVersion_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		baked         string
		moduleVersion string
		want          string
	}{
		{"baked tag wins", "v1.2.3", "v0.0.1", "v1.2.3"},
		{"clean module version from go install", "", "v0.0.53", "v0.0.53"},
		{"pseudo-version is a local build", "", "v0.0.50-0.20261002205412-a7a430486147", devVersion},
		{"dirty build", "", "v0.0.53+dirty", devVersion},
		{"devel build", "", "(devel)", devVersion},
		{"nothing known", "", "", devVersion},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, resolveVersion(tc.baked, tc.moduleVersion))
		})
	}
}
