package app_struct_generators

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"
)

// Postgres resources now manage their own connection lifecycle via
// internal/clients/postgres.ConnectTo{Name}/MigrateTo{Name}, so they must not
// contribute a field or an InitFuncCall to the app struct/InitDataSources
// wiring - unlike Sqlite, which still goes through sqldb.New.
func Test_generateDataSourceInitFileAndArgs_PostgresExcluded(t *testing.T) {
	t.Parallel()

	dataSources := matreshka.DataSources{
		resources.NewPostgres(resources.Name(resources.PostgresResourceName)),
		resources.NewSqlite(resources.Name(resources.SqliteResourceName + "_test")),
	}

	appContent, fileContent, err := generateDataSourceInitFileAndArgs(dataSources)
	require.NoError(t, err)
	require.NotNil(t, appContent)

	require.Len(t, appContent.Fields, 1, "postgres must not contribute a field, only sqlite should")
	require.Equal(t, "SqliteTest", appContent.Fields[0].Key)

	content := string(fileContent)
	require.Contains(t, content, "SqliteTest")
	require.NotContains(t, content, "Postgres")
}

// A Postgres-only project must not produce any InitFuncCall/field/file at
// all - Postgres manages its own connections via internal/clients/postgres.
func Test_generateDataSourceInitFileAndArgs_PostgresOnly(t *testing.T) {
	t.Parallel()

	dataSources := matreshka.DataSources{
		resources.NewPostgres(resources.Name(resources.PostgresResourceName)),
	}

	appContent, file, err := generateDataSourceInitFileAndArgs(dataSources)
	require.NoError(t, err)
	require.Nil(t, appContent)
	require.Nil(t, file)
}
