package config_generators

import (
	"go/format"
	"testing"

	"github.com/stretchr/testify/require"
	"go.vervstack.ru/matreshka/pkg/matreshka"

	"go.vervstack.ru/verv/plugins/project/config"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

func Test_GenerateConfigFolder_HasDotEnv(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		AppConfig: matreshka.NewEmptyConfig(),
	}

	configYamlBytes := []byte("app_info:\n  name: my-app\n")

	t.Run("hasDotEnv=true: load.go wires matreshka.WithEnvFile and is valid Go", func(t *testing.T) {
		t.Parallel()

		cfgFolder, err := GenerateConfigFolder(cfg, configYamlBytes, true)
		require.NoError(t, err)

		loadFile := cfgFolder.GetByPath(patterns.ConfigLoadFileName)
		require.NotNil(t, loadFile)

		require.Contains(t, string(loadFile.Content), "WithEnvFile")

		_, err = format.Source(loadFile.Content)
		require.NoError(t, err, "generated load.go must be valid Go source")
	})

	t.Run("hasDotEnv=false: load.go does not reference matreshka.WithEnvFile and is valid Go", func(t *testing.T) {
		t.Parallel()

		cfgFolder, err := GenerateConfigFolder(cfg, configYamlBytes, false)
		require.NoError(t, err)

		loadFile := cfgFolder.GetByPath(patterns.ConfigLoadFileName)
		require.NotNil(t, loadFile)

		require.NotContains(t, string(loadFile.Content), "WithEnvFile")

		_, err = format.Source(loadFile.Content)
		require.NoError(t, err, "generated load.go must be valid Go source")
	})
}
