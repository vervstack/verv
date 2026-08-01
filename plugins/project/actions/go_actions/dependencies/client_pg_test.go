package dependencies

import (
	"testing"

	"github.com/stretchr/testify/require"

	vervconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/tests/project_mock"
)

func testVervConfig() *vervconfig.VervConfig {
	return &vervconfig.VervConfig{
		Env: vervconfig.Project{
			PathsToClients: []string{"internal/clients"},
		},
	}
}

func Test_Postgres_AppendToProject_Single(t *testing.T) {
	t.Parallel()

	cfg := testVervConfig()
	proj := project_mock.GetMockProject(t)

	dep := postgresClient(dependencyBase{Name: DependencyNamePostgres, Cfg: cfg})

	err := dep.AppendToProject(proj)
	require.NoError(t, err)

	dataSources := proj.GetConfig().DataSources
	require.Len(t, dataSources, 1)
	require.Equal(t, DependencyNamePostgres, dataSources[0].GetName())

	connFile := proj.GetFolder().GetByPath("internal/clients", "postgres", "conn.go")
	require.NotNil(t, connFile, "conn.go should be generated")
	require.NotEmpty(t, connFile.Content)

	driverFile := proj.GetFolder().GetByPath("internal/clients", "postgres", "driver.go")
	require.NotNil(t, driverFile, "driver.go should be generated")
	require.NotEmpty(t, driverFile.Content)

	instancesFile := proj.GetFolder().GetByPath("internal/clients", "postgres", "instances.go")
	require.NotNil(t, instancesFile, "instances.go should be generated")

	instancesContent := string(instancesFile.Content)
	require.Contains(t, instancesContent, "func ConnectToPostgres() (*sql.DB, error)")
	require.Contains(t, instancesContent, "func MigrateToPostgres() error")
	require.NotContains(t, instancesContent, "Replica")

	// package import path should have been renamed away from the proj_name placeholder
	require.Contains(t, instancesContent, proj.GetName()+"/internal/config")

	postgresFolder := proj.GetFolder().GetByPath("internal/clients", "postgres")
	require.NotNil(t, postgresFolder)
	require.Len(t, postgresFolder.Inner, 3, "expected conn.go, driver.go and instances.go only")
}

func Test_Postgres_AppendToProject_MultipleInstances(t *testing.T) {
	t.Parallel()

	cfg := testVervConfig()
	proj := project_mock.GetMockProject(t)

	dep := postgresClient(dependencyBase{Name: DependencyNamePostgres, Cfg: cfg})
	err := dep.AppendToProject(proj)
	require.NoError(t, err)

	replicaDep := postgresClient(dependencyBase{Name: "postgres_replica", Cfg: cfg})

	err = replicaDep.AppendToProject(proj)
	require.NoError(t, err)

	dataSources := proj.GetConfig().DataSources
	require.Len(t, dataSources, 2)

	names := []string{dataSources[0].GetName(), dataSources[1].GetName()}
	require.ElementsMatch(t, []string{"postgres", "postgres_replica"}, names)

	instancesFile := proj.GetFolder().GetByPath("internal/clients", "postgres", "instances.go")
	require.NotNil(t, instancesFile)

	instancesContent := string(instancesFile.Content)
	require.Contains(t, instancesContent, "func ConnectToPostgres() (*sql.DB, error)")
	require.Contains(t, instancesContent, "func MigrateToPostgres() error")
	require.Contains(t, instancesContent, "func ConnectToPostgresReplica() (*sql.DB, error)")
	require.Contains(t, instancesContent, "func MigrateToPostgresReplica() error")

	// conn.go/driver.go must not be duplicated when adding a second instance
	postgresFolder := proj.GetFolder().GetByPath("internal/clients", "postgres")
	require.NotNil(t, postgresFolder)
	require.Len(t, postgresFolder.Inner, 3, "expected conn.go, driver.go and instances.go only")
}

func Test_Postgres_AppendToProject_Idempotent(t *testing.T) {
	t.Parallel()

	cfg := testVervConfig()
	proj := project_mock.GetMockProject(t)

	dep := postgresClient(dependencyBase{Name: DependencyNamePostgres, Cfg: cfg})

	err := dep.AppendToProject(proj)
	require.NoError(t, err)

	connFileFirst := proj.GetFolder().GetByPath("internal/clients", "postgres", "conn.go")
	require.NotNil(t, connFileFirst)

	firstContent := string(connFileFirst.Content)

	err = dep.AppendToProject(proj)
	require.NoError(t, err)

	dataSources := proj.GetConfig().DataSources
	require.Len(t, dataSources, 1, "re-adding postgres must not duplicate the DataSources entry")

	postgresFolder := proj.GetFolder().GetByPath("internal/clients", "postgres")
	require.NotNil(t, postgresFolder)
	require.Len(t, postgresFolder.Inner, 3, "conn.go/driver.go must not be duplicated")

	connFileSecond := proj.GetFolder().GetByPath("internal/clients", "postgres", "conn.go")
	require.NotNil(t, connFileSecond)
	require.Equal(t, firstContent, string(connFileSecond.Content))

	instancesFile := proj.GetFolder().GetByPath("internal/clients", "postgres", "instances.go")
	require.NotNil(t, instancesFile)

	instancesContent := string(instancesFile.Content)
	require.Equal(t, 1, countOccurrences(instancesContent, "func ConnectToPostgres()"))
}

func countOccurrences(s, substr string) int {
	count := 0

	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			count++
		}
	}

	return count
}
