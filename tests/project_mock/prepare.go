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

	p := &MockProject{
		Project: &project.Project{
			Name: "github.com/" + t.Name(),
			Cfg: &config.Config{
				AppConfig: matreshka.NewEmptyConfig(),
			},
			Root: folder.Folder{},
		},
	}

	require.NoError(t, p.Cfg.Unmarshal(basicConfigFile))

	for _, o := range opts {
		o(p)
	}

	cfgMarshalled, err := p.Cfg.Marshal()
	require.NoError(t, err)

	p.Cfg.AppConfig = matreshka.NewEmptyConfig()
	// to be sure in types of env variables
	require.NoError(t, p.Cfg.Unmarshal(cfgMarshalled))

	masterConfigPath := path.Join(patterns.ConfigsFolder, patterns.ConfigMasterYamlFile)
	if p.Root.GetByPath(masterConfigPath) == nil {
		p.Root.Add(
			&folder.Folder{
				Name:    masterConfigPath,
				Content: cfgMarshalled,
			},
		)
	}

	if p.Path != "" {
		require.NoError(t, os.RemoveAll(p.Path))
		require.NoError(t, os.MkdirAll(p.Path, io.DefaultDirPerm))
	}

	return p
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
