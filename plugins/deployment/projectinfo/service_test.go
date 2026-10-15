package projectinfo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/deployment/projectinfo"
)

func Test_ServiceName_Variants(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{"service name path", "service:\n  name: github.com/ruf-dev/artel\n", "artel"},
		{"legacy name", "name: go.zpotify.ru/zpotify\n", "zpotify"},
		{"service name wins over legacy", "name: legacy\nservice:\n  name: modern\n", "modern"},
		{"uppercase and junk characters", "name: Foo Bar/My_Svc!\n", "my_svc"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, ".verv", "vervonomicon.yaml"), tc.content)

			got, err := projectinfo.ServiceName(dir)

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func Test_ServiceName_MissingFile(t *testing.T) {
	_, err := projectinfo.ServiceName(t.TempDir())

	require.ErrorIs(t, err, projectinfo.ErrServiceNameUnknown)
}

func Test_ServiceName_EmptyName(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".verv", "vervonomicon.yaml"), "service:\n  name: \"\"\n")

	_, err := projectinfo.ServiceName(dir)

	require.ErrorIs(t, err, projectinfo.ErrServiceNameUnknown)
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()

	err := os.MkdirAll(filepath.Dir(path), 0o700)
	require.NoError(t, err)

	err = os.WriteFile(path, []byte(content), 0o600)
	require.NoError(t, err)
}
