package app_struct_generators

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/tests/project_mock"
)

var (
	//go:embed test_snapshots/basic_app.go.snapshot
	basicAppFile string
	//go:embed test_snapshots/with_sqlite_app.go.snapshot
	withSqliteAppFile string
	//go:embed test_snapshots/with_server_app.go.snapshot
	withServerAppFile string
	//go:embed test_snapshots/with_sqlite_and_server_app.go.snapshot
	withSqliteAndServerAppFile string

	//go:embed test_snapshots/with_sqlite_data_sources.go.snapshot
	dataSourcesGoFile []byte
	//go:embed test_snapshots/with_server_server.go.snapshot
	serverGoFile []byte
)

func Test_GenerateAppFiles(t *testing.T) {
	t.Parallel()

	type testCase struct {
		genProj      func() project.IProject
		expectedApp  string
		expectedKeys map[string][]byte
	}

	testCases := map[string]testCase{
		"basic": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t)
			},
			expectedApp: basicAppFile,
			expectedKeys: map[string][]byte{
				patterns.AppConfigFileName: appConfigPattern,
				patterns.AppCustomFileName: customPattern,
			},
		},
		"with_postgres": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t, project_mock.WithPostgres("test"))
			},
			// Postgres manages its own connections via internal/clients/postgres,
			// so a Postgres-only project produces the same app.go as basic - no
			// InitDataSources wiring, no extra generated file.
			expectedApp: basicAppFile,
			expectedKeys: map[string][]byte{
				patterns.AppConfigFileName: appConfigPattern,
				patterns.AppCustomFileName: customPattern,
			},
		},
		"with_sqlite": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t, project_mock.WithSqlite("test"))
			},
			expectedApp: withSqliteAppFile,
			expectedKeys: map[string][]byte{
				patterns.AppConfigFileName:          appConfigPattern,
				patterns.AppCustomFileName:          customPattern,
				patterns.AppInitDataSourcesFileName: dataSourcesGoFile,
			},
		},
		"with_server": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t, project_mock.WithGrpcServer(50051))
			},
			expectedApp: withServerAppFile,
			expectedKeys: map[string][]byte{
				patterns.AppConfigFileName:     appConfigPattern,
				patterns.AppCustomFileName:     customPattern,
				patterns.AppInitServerFileName: serverGoFile,
			},
		},
		"with_sqlite_and_server": {
			genProj: func() project.IProject {
				return project_mock.GetMockProject(t,
					project_mock.WithSqlite("test"),
					project_mock.WithGrpcServer(50051))
			},
			expectedApp: withSqliteAndServerAppFile,
			expectedKeys: map[string][]byte{
				patterns.AppConfigFileName:          appConfigPattern,
				patterns.AppCustomFileName:          customPattern,
				patterns.AppInitDataSourcesFileName: dataSourcesGoFile,
				patterns.AppInitServerFileName:      serverGoFile,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			proj := tc.genProj()

			out, err := GenerateAppFiles(proj)
			require.NoError(t, err)

			require.Equal(t, tc.expectedApp, string(out[patterns.AppFileName]))

			expectedKeys := tc.expectedKeys

			expectedKeys[patterns.AppFileName] = out[patterns.AppFileName]
			require.Len(t, out, len(expectedKeys))

			for k, v := range expectedKeys {
				require.Equal(t, string(v), string(out[k]), "file %s", k)
			}
		})
	}
}

func Test_GenerateAppFiles_CustomFileAlreadyExists(t *testing.T) {
	t.Parallel()

	proj := project_mock.GetMockProject(t)

	proj.GetFolder().Add(&folder.Folder{
		Name:    patterns.InternalFolder + "/" + patterns.AppFolder + "/" + patterns.AppCustomFileName,
		Content: []byte("package app\n"),
	})

	out, err := GenerateAppFiles(proj)
	require.NoError(t, err)

	_, ok := out[patterns.AppCustomFileName]
	require.False(t, ok, "custom.go should not be regenerated if it already exists")
}
