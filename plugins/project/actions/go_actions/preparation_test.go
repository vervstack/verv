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

// Test_PrepareServer_KeepsExistingTransportAndMiddlewareFiles proves the fix for a bug where
// PrepareServer.Do unconditionally re-generated manager.go, http.go and the cookie middleware
// files on every tidy run, silently discarding hand edits - none of these files carry a
// "Code generated" header, so (like custom.go, see app_struct_generators.GenerateAppFiles) they
// are meant to be scaffolded once and then left alone.
func Test_PrepareServer_KeepsExistingTransportAndMiddlewareFiles(t *testing.T) {
	t.Parallel()

	const customManager = "package transport\n// hand-edited, must survive tidy\n"

	const customCookieResponse = "package middleware\n// hand-edited, must survive tidy\n"

	proj := project_mock.GetMockProject(t, project_mock.WithGrpcServer(50051),
		project_mock.WithFile(
			patterns.InternalFolder+"/"+patterns.TransportFolder+"/"+patterns.ServerManagerFileName,
			[]byte(customManager),
		),
		project_mock.WithFile(
			patterns.InternalFolder+"/"+patterns.MiddlewareFolder+"/"+patterns.CookieResponseFileName,
			[]byte(customCookieResponse),
		),
	)

	err := PrepareServer{}.Do(proj)
	require.NoError(t, err)

	transportFolder := proj.GetFolder().GetByPath(patterns.InternalFolder, patterns.TransportFolder)
	manager := transportFolder.GetByPath(patterns.ServerManagerFileName)
	require.Equal(t, customManager, string(manager.Content), "manager.go must not be regenerated if it already exists")

	middlewareFolder := proj.GetFolder().GetByPath(patterns.InternalFolder, patterns.MiddlewareFolder)
	cookieResponse := middlewareFolder.GetByPath(patterns.CookieResponseFileName)
	require.Equal(t, customCookieResponse, string(cookieResponse.Content),
		"cookie_response.go must not be regenerated if it already exists")

	// Files that weren't already present are still scaffolded normally.
	stillScaffolded := []string{patterns.GrpcServerFileName, patterns.HttpServerFileName, patterns.GatewayMuxFileName}
	for _, fileName := range stillScaffolded {
		f := transportFolder.GetByPath(fileName)
		require.NotNilf(t, f, "internal/transport/%s must still be generated", fileName)
		require.NotEmptyf(t, f.Content, "internal/transport/%s must have content", fileName)
	}

	for _, fileName := range []string{
		patterns.CookieNamesFileName, patterns.CookieAnnotatorFileName,
		patterns.CSRFInterceptorFileName, patterns.RequestSchemeAnnotatorFileName,
	} {
		f := middlewareFolder.GetByPath(fileName)
		require.NotNilf(t, f, "internal/middleware/%s must still be generated", fileName)
		require.NotEmptyf(t, f.Content, "internal/middleware/%s must have content", fileName)
	}
}
