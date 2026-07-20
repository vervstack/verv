package server_generators

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.vervstack.ru/matreshka/pkg/matreshka/server"

	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

func Test_GenerateFileServer(t *testing.T) {
	testCases := map[string]struct {
		fs       server.FS
		expected string
	}{
		"basic": {
			fs: server.FS{Dist: "dist"},
			expected: `package web

import (
	"embed"
	"io/fs"
	"net/http"

	"go.redsock.ru/rerrors"
)

//go:embed all:dist
var distFS embed.FS

func NewServer() (http.Handler, error) {
	mux := http.NewServeMux()

	distSub, err := fs.Sub(distFS, dist)
	if err != nil {
	return nil, rerrors.Wrap(err, "error creating dist fs")
	}

	ffs := http.FileServer(http.FS(distSub))
	mux.Handle("/*", ffs)

	return mux, nil
}
`,
		},
		"nested_dist_path": {
			fs: server.FS{Dist: "web/build"},
			expected: `package web

import (
	"embed"
	"io/fs"
	"net/http"

	"go.redsock.ru/rerrors"
)

//go:embed all:web/build
var distFS embed.FS

func NewServer() (http.Handler, error) {
	mux := http.NewServeMux()

	distSub, err := fs.Sub(distFS, web/build)
	if err != nil {
	return nil, rerrors.Wrap(err, "error creating dist fs")
	}

	ffs := http.FileServer(http.FS(distSub))
	mux.Handle("/*", ffs)

	return mux, nil
}
`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			f, err := GenerateFileServer(tc.fs)
			require.NoError(t, err)

			require.Equal(t, patterns.AppInitServerFileName, f.Name)
			require.Equal(t, tc.expected, string(f.Content))
		})
	}
}
