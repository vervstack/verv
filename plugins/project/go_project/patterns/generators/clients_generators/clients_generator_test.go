package clients_generators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateGRPCConn(t *testing.T) {
	t.Parallel()

	got, err := GenerateGRPCConn()
	require.NoError(t, err)
	require.Equal(t, grpcConnPattern, string(got))
}

func Test_GenerateRedisConn(t *testing.T) {
	t.Parallel()

	got, err := GenerateRedisConn()
	require.NoError(t, err)
	require.Equal(t, redisConnPattern, string(got))
}

func Test_GenerateSQLConn(t *testing.T) {
	t.Parallel()

	got, err := GenerateSQLConn()
	require.NoError(t, err)
	require.Equal(t, sqlConnPattern, string(got))
}

func Test_GeneratePostgresDriver(t *testing.T) {
	t.Parallel()

	got, err := GeneratePostgresDriver()
	require.NoError(t, err)
	require.Equal(t, postgresDriverPattern, string(got))
}

func Test_GenerateSqliteDriver(t *testing.T) {
	t.Parallel()

	got, err := GenerateSqliteDriver()
	require.NoError(t, err)
	require.Equal(t, sqliteDriverPattern, string(got))
}

func Test_GenerateTelegramConn(t *testing.T) {
	t.Parallel()

	got, err := GenerateTelegramConn()
	require.NoError(t, err)
	require.Equal(t, telegramConnPattern, string(got))
}

func Test_GeneratePostgresConn(t *testing.T) {
	t.Parallel()

	got, err := GeneratePostgresConn()
	require.NoError(t, err)
	require.Equal(t, postgresConnPattern, string(got))
}

func Test_GeneratePostgresDriver_PostgresPackage(t *testing.T) {
	t.Parallel()

	got, err := GeneratePostgresDriver()
	require.NoError(t, err)
	require.Equal(t, postgresDriverPattern, string(got))
	require.Contains(t, string(got), "package postgres")
}

func Test_GeneratePostgresInstances(t *testing.T) {
	t.Parallel()

	got, err := GeneratePostgresInstances(PostgresInstancesArgs{
		Instances: []PostgresInstance{
			{PascalName: "Postgres"},
			{PascalName: "PostgresReplica"},
		},
	})
	require.NoError(t, err)

	content := string(got)
	require.Contains(t, content, "func ConnectToPostgres() (*sql.DB, error)")
	require.Contains(t, content, "func MigrateToPostgres() error")
	require.Contains(t, content, "func ConnectToPostgresReplica() (*sql.DB, error)")
	require.Contains(t, content, "func MigrateToPostgresReplica() error")
}

func Test_GeneratePostgresInstances_Empty(t *testing.T) {
	t.Parallel()

	got, err := GeneratePostgresInstances(PostgresInstancesArgs{})
	require.NoError(t, err)
	require.Contains(t, string(got), "package postgres")
	require.NotContains(t, string(got), "func ConnectTo")
}
