package projectinfo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/deployment/projectinfo"
)

func writeEmptyFile(t *testing.T, path string) {
	t.Helper()

	err := os.MkdirAll(filepath.Dir(path), 0o700)
	require.NoError(t, err)

	err = os.WriteFile(path, nil, 0o600)
	require.NoError(t, err)
}

func Test_Environments_UnionOrderedByPriority(t *testing.T) {
	dir := t.TempDir()

	writeEmptyFile(t, filepath.Join(dir, "config", "config.yaml"))
	writeEmptyFile(t, filepath.Join(dir, "config", "config_template.yaml"))
	writeEmptyFile(t, filepath.Join(dir, "config", "local.yaml"))
	writeEmptyFile(t, filepath.Join(dir, "config", "zeta.yaml"))
	writeEmptyFile(t, filepath.Join(dir, "config", "notes.txt"))
	writeEmptyFile(t, filepath.Join(dir, ".verv", "nodes.local.yaml"))
	writeEmptyFile(t, filepath.Join(dir, ".verv", "stage", "values.yaml"))
	writeEmptyFile(t, filepath.Join(dir, ".verv", "alpha", "values.yaml"))
	writeEmptyFile(t, filepath.Join(dir, ".verv", ".cache", "values.yaml"))
	writeEmptyFile(t, filepath.Join(dir, ".verv", "local", "values.yaml"))

	got, err := projectinfo.Environments(dir)

	require.NoError(t, err)
	require.Equal(t, []string{"stage", "local", "prod", "alpha", "zeta"}, got)
}

func Test_Environments_MissingDirectories(t *testing.T) {
	got, err := projectinfo.Environments(t.TempDir())

	require.NoError(t, err)
	require.Empty(t, got)
}
