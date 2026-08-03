package go_actions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/tests/project_mock"
)

func Test_GenerateProjectConfig_ServerEnvVars(t *testing.T) {
	t.Parallel()

	t.Run("with server: allowed_origins/cookie_secure are registered and round-trip through Marshal", func(t *testing.T) {
		t.Parallel()

		proj := project_mock.GetMockProject(t, project_mock.WithGrpcServer(50051))

		err := GenerateProjectConfig{}.Do(proj)
		require.NoError(t, err)

		envVars := map[string]bool{}
		for _, v := range proj.GetConfig().Environment {
			envVars[v.Name] = true
		}

		require.True(t, envVars[AllowedOriginsEvonName], "allowed_origins must be registered when servers are present")
		require.True(t, envVars[CookieSecureEvonName], "cookie_secure must be registered when servers are present")

		// Prove environment.Variable.Comment actually round-trips into the marshalled
		// config YAML via matreshka.AppConfig.Marshal() (plain yaml.Marshal, no custom
		// MarshalYAML on Variable that would drop the field).
		appConfig := proj.GetConfig().AppConfig

		marshalled, err := appConfig.Marshal()
		require.NoError(t, err)
		require.Contains(t, string(marshalled), "comment:")
	})

	t.Run("without server: allowed_origins/cookie_secure are not registered", func(t *testing.T) {
		t.Parallel()

		proj := project_mock.GetMockProject(t)

		err := GenerateProjectConfig{}.Do(proj)
		require.NoError(t, err)

		for _, v := range proj.GetConfig().Environment {
			require.NotEqual(t, AllowedOriginsEvonName, v.Name, "allowed_origins must not be registered without servers")
			require.NotEqual(t, CookieSecureEvonName, v.Name, "cookie_secure must not be registered without servers")
		}
	})
}
