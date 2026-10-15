package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/plugins/deployment/imagebuild"
	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

const (
	testServiceName = "billing"
	testEnvironment = "prod"
	testProjectDir  = "/work/billing"
)

var errBoom = rerrors.New("boom")

type fixture struct {
	log     *callLog
	node    *fakeNode
	project *fakeProject
	builder *fakeBuilder
	printer *nopIO
}

func newFixture() *fixture {
	log := &callLog{}

	return &fixture{
		log: log,
		node: &fakeNode{
			log:      log,
			platform: nodeclient.Platform{Version: "1.0.0", Os: "linux", Arch: "arm64"},
			names:    []string{"main", "backup"},
			credentials: nodeclient.RegistryCredentials{
				Username:    "user",
				Password:    "secret",
				RegistryUrl: "registry.example.com",
			},
		},
		project: &fakeProject{log: log, serviceName: testServiceName, tag: "v1-abc"},
		builder: &fakeBuilder{log: log},
		printer: &nopIO{},
	}
}

func (f *fixture) deps() pipeline.Deps {
	return pipeline.Deps{
		Node:    f.node,
		Project: f.project,
		Builder: f.builder,
		Printer: f.printer,
	}
}

func newParams() pipeline.Params {
	return pipeline.Params{
		ProjectDir:  testProjectDir,
		Dockerfile:  "Dockerfile.prod",
		Environment: testEnvironment,
	}
}

func Test_Run_BuildsAndDeploysInOrder(t *testing.T) {
	t.Parallel()

	f := newFixture()
	wantRef := imagebuild.ImageRef("registry.example.com", testServiceName, "v1-abc")

	err := pipeline.Run(context.Background(), f.deps(), newParams())

	require.NoError(t, err)
	require.Equal(t, []string{
		"node.Version",
		"node.ListRegistryNames",
		"node.RegistryCredentials",
		"project.ServiceName",
		"project.Tag",
		"builder.Build",
		"node.CreateDeploy",
	}, f.log.calls)
	require.Equal(t, "main", f.node.credentialsName)

	wantBuild := imagebuild.Request{
		ProjectDir: testProjectDir,
		Dockerfile: "Dockerfile.prod",
		Platform:   "linux/arm64",
		ImageRef:   wantRef,
		Credentials: imagebuild.Credentials{
			Username:    "user",
			Password:    "secret",
			RegistryUrl: "registry.example.com",
		},
	}
	require.Equal(t, []imagebuild.Request{wantBuild}, f.builder.requests)

	wantDeploy := nodeclient.DeployRequest{
		ServiceName: testServiceName,
		Environment: testEnvironment,
		Image:       wantRef,
	}
	require.Equal(t, []nodeclient.DeployRequest{wantDeploy}, f.node.deployRequests)

	require.Contains(t, f.printer.lines, "✅ Node platform: linux/arm64")
	require.Contains(t, f.printer.lines, "✅ Registry: registry.example.com")
	require.Contains(t, f.printer.lines, "✅ Image: "+wantRef)
	require.Contains(t, f.printer.lines, "✅ Deployed billing to prod")
}

func Test_Run_TagOverrideSkipsProjectTag(t *testing.T) {
	t.Parallel()

	f := newFixture()
	params := newParams()
	params.Tag = "manual-1"

	err := pipeline.Run(context.Background(), f.deps(), params)

	require.NoError(t, err)
	require.NotContains(t, f.log.calls, "project.Tag")
	require.Equal(t, imagebuild.ImageRef("registry.example.com", testServiceName, "manual-1"), f.builder.requests[0].ImageRef)
}

func Test_Run_NoRegistry(t *testing.T) {
	t.Parallel()

	f := newFixture()
	f.node.names = nil

	err := pipeline.Run(context.Background(), f.deps(), newParams())

	require.ErrorIs(t, err, pipeline.ErrNoRegistry)
	require.Empty(t, f.builder.requests)
}

func Test_Run_NodePlatformUnknown(t *testing.T) {
	t.Parallel()

	f := newFixture()
	f.node.platform = nodeclient.Platform{Version: "0.9.0"}

	err := pipeline.Run(context.Background(), f.deps(), newParams())

	require.ErrorIs(t, err, pipeline.ErrNodePlatformUnknown)
	require.Equal(t, []string{"node.Version"}, f.log.calls)
}

func Test_Run_EnvironmentRequiredMakesNoCalls(t *testing.T) {
	t.Parallel()

	f := newFixture()
	params := newParams()
	params.Environment = ""

	err := pipeline.Run(context.Background(), f.deps(), params)

	require.ErrorIs(t, err, pipeline.ErrEnvironmentRequired)
	require.Empty(t, f.log.calls)
}

func Test_Run_BuildFailureStopsBeforeDeploy(t *testing.T) {
	t.Parallel()

	f := newFixture()
	f.builder.err = errBoom

	err := pipeline.Run(context.Background(), f.deps(), newParams())

	require.ErrorIs(t, err, errBoom)
	require.Empty(t, f.node.deployRequests)
	require.NotContains(t, f.log.calls, "node.CreateDeploy")
}

func Test_Run_DeployFailureIsReturned(t *testing.T) {
	t.Parallel()

	f := newFixture()
	f.node.deployErr = errBoom

	err := pipeline.Run(context.Background(), f.deps(), newParams())

	require.ErrorIs(t, err, errBoom)
}
