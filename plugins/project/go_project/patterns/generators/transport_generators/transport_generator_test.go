package transport_generators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateServerManager(t *testing.T) {
	t.Parallel()

	got, err := GenerateServerManager()
	require.NoError(t, err)
	require.Equal(t, serverManagerPattern, string(got))
}

func Test_GenerateGrpcServer(t *testing.T) {
	t.Parallel()

	got, err := GenerateGrpcServer()
	require.NoError(t, err)
	require.Equal(t, grpcServerPattern, string(got))
}

func Test_GenerateHttpServer(t *testing.T) {
	t.Parallel()

	got, err := GenerateHttpServer()
	require.NoError(t, err)
	require.Equal(t, httpServerPattern, string(got))
}

func Test_GenerateTelegramListener(t *testing.T) {
	t.Parallel()

	got, err := GenerateTelegramListener()
	require.NoError(t, err)
	require.Equal(t, telegramListenerPattern, string(got))
}

func Test_GenerateTelegramVersionHandler(t *testing.T) {
	t.Parallel()

	got, err := GenerateTelegramVersionHandler()
	require.NoError(t, err)
	require.Equal(t, telegramVersionHandlerPattern, string(got))
}
