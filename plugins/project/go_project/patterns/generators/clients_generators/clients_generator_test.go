package clients_generators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateGRPCConn(t *testing.T) {
	got, err := GenerateGRPCConn()
	require.NoError(t, err)
	require.Equal(t, grpcConnPattern, string(got))
}

func Test_GenerateRedisConn(t *testing.T) {
	got, err := GenerateRedisConn()
	require.NoError(t, err)
	require.Equal(t, redisConnPattern, string(got))
}

func Test_GenerateSQLConn(t *testing.T) {
	got, err := GenerateSQLConn()
	require.NoError(t, err)
	require.Equal(t, sqlConnPattern, string(got))
}

func Test_GeneratePostgresDriver(t *testing.T) {
	got, err := GeneratePostgresDriver()
	require.NoError(t, err)
	require.Equal(t, postgresDriverPattern, string(got))
}

func Test_GenerateSqliteDriver(t *testing.T) {
	got, err := GenerateSqliteDriver()
	require.NoError(t, err)
	require.Equal(t, sqliteDriverPattern, string(got))
}

func Test_GenerateTelegramConn(t *testing.T) {
	got, err := GenerateTelegramConn()
	require.NoError(t, err)
	require.Equal(t, telegramConnPattern, string(got))
}
