package go_actions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/tests/project_mock"
)

func Test_GenerateProjectConfig_ServerEnvVars(t *testing.T) {
	t.Parallel()

	t.Run("with server: allowed_origins is registered and round-trips through Marshal", func(t *testing.T) {
		t.Parallel()

		proj := project_mock.GetMockProject(t, project_mock.WithGrpcServer(50051))

		err := GenerateProjectConfig{}.Do(proj)
		require.NoError(t, err)

		envVars := map[string]bool{}
		for _, v := range proj.GetConfig().Environment {
			envVars[v.Name] = true
		}

		require.True(t, envVars[AllowedOriginsEvonName], "allowed_origins must be registered when servers are present")

		// Prove environment.Variable.Comment actually round-trips into the marshalled
		// config YAML via matreshka.AppConfig.Marshal() (plain yaml.Marshal, no custom
		// MarshalYAML on Variable that would drop the field).
		appConfig := proj.GetConfig().AppConfig

		marshalled, err := appConfig.Marshal()
		require.NoError(t, err)
		require.Contains(t, string(marshalled), "comment:")
	})

	t.Run("without server: allowed_origins is not registered", func(t *testing.T) {
		t.Parallel()

		proj := project_mock.GetMockProject(t)

		err := GenerateProjectConfig{}.Do(proj)
		require.NoError(t, err)

		for _, v := range proj.GetConfig().Environment {
			require.NotEqual(t, AllowedOriginsEvonName, v.Name, "allowed_origins must not be registered without servers")
		}
	})
}

func Test_PrepareConfigFolder_DotEnv(t *testing.T) {
	t.Parallel()

	// .env presence is checked by the generated app at runtime (os.Stat), not
	// by tidy, so load.go must wire matreshka.WithEnvFile the same way whether
	// or not the project happens to have a .env file when tidy runs.
	t.Run("with .env: generated load.go wires matreshka.WithEnvFile", func(t *testing.T) {
		t.Parallel()

		proj := project_mock.GetMockProject(t, project_mock.WithFile(".env", []byte("FOO=bar\n")))

		err := PrepareConfigFolder{}.Do(proj)
		require.NoError(t, err)

		loadGoContent := resolveLoadGoContent(t, proj)
		require.Contains(t, string(loadGoContent), "WithEnvFile")
	})

	t.Run("without .env: generated load.go still wires matreshka.WithEnvFile", func(t *testing.T) {
		t.Parallel()

		proj := project_mock.GetMockProject(t)

		err := PrepareConfigFolder{}.Do(proj)
		require.NoError(t, err)

		loadGoContent := resolveLoadGoContent(t, proj)
		require.Contains(t, string(loadGoContent), "WithEnvFile")
	})
}

func resolveLoadGoContent(t *testing.T, proj *project_mock.MockProject) []byte {
	t.Helper()

	cfgFolder := proj.GetFolder().GetByPath(patterns.InternalFolder + "/" + patterns.ConfigsFolder)
	require.NotNil(t, cfgFolder, "internal/config folder must be generated")

	loadFile := cfgFolder.GetByPath(patterns.ConfigLoadFileName)
	require.NotNil(t, loadFile, "load.go must be generated")

	return loadFile.Content
}
