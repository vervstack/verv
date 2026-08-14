package go_actions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/tests/project_mock"
)

// Test_PrepareServer_AttachesTransportFolder proves the fix for a pre-existing bug: when
// internal/transport doesn't already exist in the project's folder tree, PrepareServer.Do
// built a transportFolder, populated it, but never attached it back into the tree, so the
// generated files were silently discarded. The mock project here starts with no
// internal/transport folder at all (WithGrpcServer only sets cfg.Servers), which is exactly
// the code path the old code got wrong. It also proves the parallel internal/middleware
// folder (new in this change) is attached the same way.
func Test_PrepareServer_AttachesTransportFolder(t *testing.T) {
	t.Parallel()

	proj := project_mock.GetMockProject(t, project_mock.WithGrpcServer(50051))

	require.Nil(t, proj.GetFolder().GetByPath(patterns.InternalFolder, patterns.TransportFolder))
	require.Nil(t, proj.GetFolder().GetByPath(patterns.InternalFolder, patterns.MiddlewareFolder))

	err := PrepareServer{}.Do(proj)
	require.NoError(t, err)

	transportFolder := proj.GetFolder().GetByPath(patterns.InternalFolder, patterns.TransportFolder)
	require.NotNil(t, transportFolder, "internal/transport must be attached to the project's folder tree")

	for _, fileName := range []string{
		patterns.ServerManagerFileName,
		patterns.GrpcServerFileName,
		patterns.HttpServerFileName,
		patterns.GatewayMuxFileName,
	} {
		f := transportFolder.GetByPath(fileName)
		require.NotNilf(t, f, "internal/transport/%s must exist", fileName)
		require.NotEmptyf(t, f.Content, "internal/transport/%s must have content", fileName)
	}

	middlewareFolder := proj.GetFolder().GetByPath(patterns.InternalFolder, patterns.MiddlewareFolder)
	require.NotNil(t, middlewareFolder, "internal/middleware must be attached to the project's folder tree")

	for _, fileName := range []string{
		patterns.CookieNamesFileName,
		patterns.CookieAnnotatorFileName,
		patterns.CookieResponseFileName,
		patterns.CSRFInterceptorFileName,
		patterns.RequestSchemeAnnotatorFileName,
	} {
		f := middlewareFolder.GetByPath(fileName)
		require.NotNilf(t, f, "internal/middleware/%s must exist", fileName)
		require.NotEmptyf(t, f.Content, "internal/middleware/%s must have content", fileName)
	}
}
