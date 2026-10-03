package upgrade

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func newBinaryServer(t *testing.T, body string) *httptest.Server {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return server
}

func Test_Install_KeepsDestinationMode(t *testing.T) {
	t.Parallel()

	server := newBinaryServer(t, "new")
	dest := filepath.Join(t.TempDir(), "verv")

	err := os.WriteFile(dest, []byte("old"), 0o750)
	require.NoError(t, err)

	err = os.Chmod(dest, 0o750)
	require.NoError(t, err)

	err = install(t.Context(), server.URL, dest)
	require.NoError(t, err)

	info, err := os.Stat(dest)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o750), info.Mode().Perm())

	content, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, "new", string(content))
}

func Test_Install_MissingDestinationUsesDefaultMode(t *testing.T) {
	t.Parallel()

	server := newBinaryServer(t, "new")
	dest := filepath.Join(t.TempDir(), "verv")

	err := install(t.Context(), server.URL, dest)
	require.NoError(t, err)

	info, err := os.Stat(dest)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(binaryMode), info.Mode().Perm())
}
