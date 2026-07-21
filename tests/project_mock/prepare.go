package project_mock

import (
	_ "embed"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"go.vervstack.ru/matreshka/pkg/matreshka"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/config"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

type MockProject struct {
	*project.Project
}

type Opt func(m *MockProject)

func GetMockProject(t *testing.T, opts ...Opt) *MockProject {
	t.Helper()

	projectMock := &MockProject{
		Project: &project.Project{
			Name: "github.com/" + t.Name(),
			Cfg: &config.Config{
				AppConfig: matreshka.NewEmptyConfig(),
			},
			Root: folder.Folder{},
		},
	}

	require.NoError(t, projectMock.Cfg.Unmarshal(basicConfigFile))

	for _, o := range opts {
		o(projectMock)
	}

	cfgMarshalled, err := projectMock.Cfg.Marshal()
	require.NoError(t, err)

	projectMock.Cfg.AppConfig = matreshka.NewEmptyConfig()
	// to be sure in types of env variables
	require.NoError(t, projectMock.Cfg.Unmarshal(cfgMarshalled))

	masterConfigPath := path.Join(patterns.ConfigsFolder, patterns.ConfigMasterYamlFile)
	if projectMock.Root.GetByPath(masterConfigPath) == nil {
		projectMock.Root.Add(
			&folder.Folder{
				Name:    masterConfigPath,
				Content: cfgMarshalled,
			},
		)
	}

	if projectMock.Path != "" {
		require.NoError(t, os.RemoveAll(projectMock.Path))
		require.NoError(t, os.MkdirAll(projectMock.Path, io.DefaultDirPerm))
	}

	return projectMock
}

func (m *MockProject) WriteFile(t *testing.T, relativePath string, data []byte) {
	t.Helper()

	relativePath = path.Join(m.Path, relativePath)

	require.NoError(t, os.MkdirAll(path.Dir(relativePath), io.DefaultDirPerm))

	cfgFile, err := os.Create(relativePath)
	require.NoError(t, err)

	_, err = cfgFile.Write(data)
	require.NoError(t, err)
	require.NoError(t, cfgFile.Close())
}
