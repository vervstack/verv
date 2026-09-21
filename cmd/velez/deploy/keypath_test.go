package deploy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExpandKeyPath_Scenarios(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"bare tilde", "~", home},
		{"tilde slash path", "~/velez", filepath.Join(home, "velez")},
		{"absolute path untouched", "/data/velez", "/data/velez"},
		{"relative path untouched", "velez-keys", "velez-keys"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := expandKeyPath(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
