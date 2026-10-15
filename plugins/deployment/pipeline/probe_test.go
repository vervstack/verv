package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

func Test_Probe_VersionFails(t *testing.T) {
	t.Parallel()

	f := newFixture()
	f.node.versionErr = errBoom

	status := pipeline.Probe(context.Background(), f.node)

	require.ErrorIs(t, status.Err, errBoom)
	require.False(t, status.IsRegistryEnabled)
	require.Equal(t, []string{"node.Version"}, f.log.calls)
}

func Test_Probe_RegistryListFails(t *testing.T) {
	t.Parallel()

	f := newFixture()
	f.node.namesErr = errBoom

	status := pipeline.Probe(context.Background(), f.node)

	require.ErrorIs(t, status.Err, errBoom)
	require.False(t, status.IsRegistryEnabled)
	require.Equal(t, f.node.platform, status.Platform)
}

func Test_Probe_NoRegistries(t *testing.T) {
	t.Parallel()

	f := newFixture()
	f.node.names = nil

	status := pipeline.Probe(context.Background(), f.node)

	require.NoError(t, status.Err)
	require.False(t, status.IsRegistryEnabled)
	require.Equal(t, f.node.platform, status.Platform)
}

func Test_Probe_Ok(t *testing.T) {
	t.Parallel()

	f := newFixture()

	status := pipeline.Probe(context.Background(), f.node)

	require.NoError(t, status.Err)
	require.True(t, status.IsRegistryEnabled)
	require.Equal(t, "arm64", status.Platform.Arch)
}
