package velez

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PrepareKeyDir_CustomPath_NoOp(t *testing.T) {
	keyPath := filepath.Join(t.TempDir(), "x")

	err := PrepareKeyDir(keyPath)
	require.NoError(t, err)

	_, err = os.Stat(keyPath)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func Test_LinkInHome_CreatesSymlink(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	err := LinkInHome(DefaultKeyPath)
	require.NoError(t, err)

	target, err := os.Readlink(filepath.Join(home, homeLinkName))
	require.NoError(t, err)
	require.Equal(t, DefaultKeyPath, target)
}

func Test_LinkInHome_ExistingEntryUntouched(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	link := filepath.Join(home, homeLinkName)
	err := os.Mkdir(link, 0o755)
	require.NoError(t, err)

	keyFile := filepath.Join(link, "key")
	err = os.WriteFile(keyFile, []byte("secret"), 0o600)
	require.NoError(t, err)

	err = LinkInHome(DefaultKeyPath)
	require.NoError(t, err)

	info, err := os.Lstat(link)
	require.NoError(t, err)
	require.True(t, info.IsDir())
	require.Zero(t, info.Mode()&os.ModeSymlink)

	content, err := os.ReadFile(keyFile)
	require.NoError(t, err)
	require.Equal(t, "secret", string(content))
}

func Test_LinkInHome_CustomPath_NoLink(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	err := LinkInHome("/data/velez")
	require.NoError(t, err)

	_, err = os.Lstat(filepath.Join(home, homeLinkName))
	require.ErrorIs(t, err, os.ErrNotExist)
}
