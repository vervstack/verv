package menu_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	vervconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/plugins/project"
)

func Test_BuildHeader_NoMarker(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	h := menu.BuildHeader(dir, &vervconfig.VervConfig{})

	require.False(t, h.IsVervProject)
	require.Empty(t, h.Name)
	require.Empty(t, h.Version)
	require.NotEmpty(t, h.Path)
}

func Test_BuildHeader_WithMarker(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeVervonomicon(t, dir, "name: my-app\n")

	h := menu.BuildHeader(dir, &vervconfig.VervConfig{})

	require.True(t, h.IsVervProject)
	require.Equal(t, "my-app", h.Name)
	require.NotEmpty(t, h.Emoji)
}

func Test_BuildHeader_MalformedMarker_TreatedAsNotAProject(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeVervonomicon(t, dir, "name: [broken\n")

	h := menu.BuildHeader(dir, &vervconfig.VervConfig{})

	require.False(t, h.IsVervProject)
}

func writeVervonomicon(t *testing.T, dir, content string) {
	t.Helper()

	markerDir := path.Join(dir, project.VervMarkerDir)
	require.NoError(t, os.MkdirAll(markerDir, io.DefaultDirPerm))
	require.NoError(t, os.WriteFile(path.Join(markerDir, project.VervonomiconFile), []byte(content), io.DefaultFilePerm))
}
