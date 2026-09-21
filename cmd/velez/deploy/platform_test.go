package deploy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PlatformAllowed_Scenarios(t *testing.T) {
	cases := []struct {
		name  string
		goos  string
		debug bool
		want  bool
	}{
		{"linux without debug", "linux", false, true},
		{"linux with debug", "linux", true, true},
		{"darwin without debug", "darwin", false, false},
		{"darwin with debug", "darwin", true, true},
		{"windows without debug", "windows", false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := platformAllowed(tc.goos, tc.debug)
			require.Equal(t, tc.want, got)
		})
	}
}
