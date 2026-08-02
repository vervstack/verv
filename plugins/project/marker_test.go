package project_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/plugins/project"
)

func Test_IsVervProject_NoMarker(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	require.False(t, project.IsVervProject(dir))
}

func Test_IsVervProject_ValidMarker(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeMarker(t, dir, "name: my-proj\n")

	require.True(t, project.IsVervProject(dir))

	v, err := project.ReadVervonomicon(dir)
	require.NoError(t, err)
	require.Equal(t, "my-proj", v.Name)
}

func Test_ReadVervonomicon_Malformed(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeMarker(t, dir, "name: [unterminated\n")

	_, err := project.ReadVervonomicon(dir)
	require.Error(t, err)
}

func Test_ReadVervonomicon_Missing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	_, err := project.ReadVervonomicon(dir)
	require.Error(t, err)
}

func writeMarker(t *testing.T, dir, content string) {
	t.Helper()

	markerDir := path.Join(dir, project.VervMarkerDir)
	require.NoError(t, os.MkdirAll(markerDir, io.DefaultDirPerm))
	require.NoError(t, os.WriteFile(path.Join(markerDir, project.VervonomiconFile), []byte(content), io.DefaultFilePerm))
}
