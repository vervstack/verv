package go_actions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/tests/project_mock"
)

func Test_PrepareDeployFolder(t *testing.T) {
	t.Parallel()

	proj := project_mock.GetMockProject(t)

	action := PrepareDeployFolder{}

	require.NoError(t, action.Do(proj))

	markerFile := proj.GetFolder().GetByPath(project.VervMarkerDir, project.VervDeployDir, project.VervonomiconFile)
	require.NotNil(t, markerFile)

	v := project.Vervonomicon{}
	require.NoError(t, yaml.Unmarshal(markerFile.Content, &v))
	require.Equal(t, proj.GetName(), v.Name)
}
