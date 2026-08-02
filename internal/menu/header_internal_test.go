package menu

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_displayPath(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	require.Equal(t, "~", displayPath(home))
	require.Equal(t, filepath.Join("~", "projects", "verv"), displayPath(filepath.Join(home, "projects", "verv")))
	require.Equal(t, "/tmp/outside-home", displayPath("/tmp/outside-home"))
}
