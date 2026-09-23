package menu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsHeadlessLinux_Scenarios(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		goos           string
		display        string
		waylandDisplay string
		want           bool
	}{
		{"linux, no session", "linux", "", "", true},
		{"linux, X11 session", "linux", ":0", "", false},
		{"linux, wayland session", "linux", "", "wayland-0", false},
		{"darwin", "darwin", "", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := isHeadlessLinux(tc.goos, tc.display, tc.waylandDisplay)

			require.Equal(t, tc.want, got)
		})
	}
}
