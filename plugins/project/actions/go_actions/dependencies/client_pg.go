package dependencies

import (
	"path"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions/renamer"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/clients_generators"
)

type Postgres struct {
	dependencyBase
}

func postgresClient(dep dependencyBase) Dependency {
	return &Postgres{dependencyBase: dep}
}

// GetFolderName returns the resource/instance name for this Postgres
// dependency (used to key into DataSources and derive the
// ConnectTo{Name}/MigrateTo{Name} function names). It is NOT the folder the
// client package lives in - unlike Redis, every Postgres instance shares a
// single "postgres" client package, see applyClientFolder.
func (p Postgres) GetFolderName() string {
	if p.Name != "" {
		return p.Name
	}

	return DependencyNamePostgres
}

func (p Postgres) AppendToProject(proj Project) error {
	err := p.applyClientFolder(proj)
	if err != nil {
		return rerrors.Wrap(err, "error applying client folder")
	}

	p.applyConfig(proj)

	err = p.applyInstances(proj)
	if err != nil {
		return rerrors.Wrap(err, "error applying instances file")
	}

	return nil
}

// applyClientFolder ensures the shared internal/clients/postgres package
// (conn.go + driver.go) exists. All named Postgres instances live in this one
// package, so it's only generated once regardless of how many instances are
// added.
func (p Postgres) applyClientFolder(proj Project) error {
	if len(p.Cfg.Env.PathsToClients) == 0 {
		return ErrNoFolderInConfig
	}

	ok, err := containsDependencyFolder(p.Cfg.Env.PathsToClients, proj.GetFolder(), DependencyNamePostgres)
	if err != nil {
		return rerrors.Wrap(err, "error finding Dependency path")
	}

	if ok {
		return nil
	}

	connContent, err := clients_generators.GeneratePostgresConn()
	if err != nil {
		return rerrors.Wrap(err, "error generating postgres conn file")
	}

	connFile := &folder.Folder{
		Name:    path.Join(p.Cfg.Env.PathsToClients[0], DependencyNamePostgres, patterns.ConnFileName),
		Content: connContent,
	}

	renamer.ReplaceProjectName(proj.GetName(), connFile)

	proj.GetFolder().Add(connFile)

	driverContent, err := clients_generators.GeneratePostgresDriver()
	if err != nil {
		return rerrors.Wrap(err, "error generating postgres driver file")
	}

	driverFile := &folder.Folder{
		Name:    path.Join(p.Cfg.Env.PathsToClients[0], DependencyNamePostgres, "driver.go"),
		Content: driverContent,
	}

	proj.GetFolder().Add(driverFile)

	return nil
}

func (p Postgres) applyConfig(proj Project) {
	cfg := proj.GetConfig()

	for _, item := range cfg.DataSources {
		if item.GetName() == p.GetFolderName() {
			return
		}
	}

	appNameInfo := proj.GetShortName()

	cfg.DataSources = append(cfg.DataSources,
		&resources.Postgres{
			Name:             resources.Name(p.GetFolderName()),
			Host:             "localhost",
			Port:             defaultPgPort,
			DbName:           appNameInfo,
			User:             appNameInfo,
			MigrationsFolder: "./migrations",
		})
}

// applyInstances (re)generates internal/clients/postgres/instances.go so it
// reflects every currently configured Postgres data source, not just the one
// being added right now.
func (p Postgres) applyInstances(proj Project) error {
	if len(p.Cfg.Env.PathsToClients) == 0 {
		return ErrNoFolderInConfig
	}

	instances := postgresInstances(proj.GetConfig().DataSources)

	content, err := clients_generators.GeneratePostgresInstances(clients_generators.PostgresInstancesArgs{
		Instances: instances,
	})
	if err != nil {
		return rerrors.Wrap(err, "error generating postgres instances file")
	}

	instancesFile := &folder.Folder{
		Name:    path.Join(p.Cfg.Env.PathsToClients[0], DependencyNamePostgres, "instances.go"),
		Content: content,
	}

	renamer.ReplaceProjectName(proj.GetName(), instancesFile)

	proj.GetFolder().Add(instancesFile)

	return nil
}

func postgresInstances(dataSources matreshka.DataSources) []clients_generators.PostgresInstance {
	instances := make([]clients_generators.PostgresInstance, 0, len(dataSources))

	for _, ds := range dataSources {
		if ds.GetType() != resources.PostgresResourceName {
			continue
		}

		instances = append(instances, clients_generators.PostgresInstance{
			PascalName: generators.NormalizeResourceName(ds.GetName()),
		})
	}

	return instances
}
