package config_generators

import (
	"go/format"
	"testing"

	"github.com/stretchr/testify/require"
	"go.vervstack.ru/matreshka/pkg/matreshka"

	"go.vervstack.ru/verv/plugins/project/config"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

func Test_GenerateConfigFolder_DotEnv(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		AppConfig: matreshka.NewEmptyConfig(),
	}

	configYamlBytes := []byte("app_info:\n  name: my-app\n")

	cfgFolder, err := GenerateConfigFolder(cfg, configYamlBytes)
	require.NoError(t, err)

	loadFile := cfgFolder.GetByPath(patterns.ConfigLoadFileName)
	require.NotNil(t, loadFile)

	// .env detection happens at application runtime (os.Stat on envFilePath),
	// not at generation time, so the wiring is unconditional here.
	require.Contains(t, string(loadFile.Content), "matreshka.WithEnvFile(envFilePath)")
	require.Contains(t, string(loadFile.Content), "os.Stat(envFilePath)")

	_, err = format.Source(loadFile.Content)
	require.NoError(t, err, "generated load.go must be valid Go source")
}
