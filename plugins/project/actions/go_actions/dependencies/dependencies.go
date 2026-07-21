package dependencies

import (
	"strings"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	vervconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io/folder"
)

var (
	ErrNoFolderInConfig = rerrors.New("no folder path in verv config")

	nameToDependencyConstructor = map[string]func(dep dependencyBase) Dependency{
		DependencyNamePostgres: postgresClient,
		DependencyNameRedis:    redisClient,
		DependencyNameTelegram: telegram,
		DependencyNameSqlite:   sqlite,

		DependencyEnvVariable: envVariable,
	}
)

type Dependency interface {
	AppendToProject(proj Project) error
}

type dependencyBase struct {
	Name string
	Cfg  *vervconfig.VervConfig
}

const (
	DependencyNameRedis    = "redis"
	DependencyNamePostgres = "postgres"
	DependencyNameTelegram = "telegram"
	DependencyNameSqlite   = "sqlite"

	DependencyEnvVariable = "env"

	defaultPgPort    = 5432
	defaultRedisPort = 6379
)

func HelpWithDependencyNames(passedDeps ...string) (missingDeps []string) {
	passedDepsMap := map[string]struct{}{}
	for _, dep := range passedDeps {
		passedDepsMap[dep] = struct{}{}
	}

	for depName := range nameToDependencyConstructor {
		_, ok := passedDepsMap[depName]
		if !ok {
			missingDeps = append(missingDeps, depName)
		}
	}

	return missingDeps
}

func GetDependencies(c *vervconfig.VervConfig, args []string) []Dependency {
	serverOpts := make([]Dependency, 0, len(args))

	for _, name := range args {
		resourceName, _, _ := strings.Cut(name, "_")

		depConstr, ok := nameToDependencyConstructor[resourceName]
		if !ok {
			continue
		}

		base := dependencyBase{
			Name: name,
			Cfg:  c,
		}

		serverOpts = append(serverOpts, depConstr(base))
	}

	return serverOpts
}

// TODO method checks if the root folder of dependecy is presented. In case it's empty - nothing is generated

// containsDependencyFolder - searches through VERV_PATH_TO_CLIENTS
// folders in order to find depName
// IF Dependency already placed - returns path to it
func containsDependencyFolder(paths []string, rootF *folder.Folder, depName string) (ok bool, err error) {
	if len(paths) == 0 {
		return false, rerrors.Wrap(ErrNoFolderInConfig, "no client")
	}

	for _, clientPath := range paths {
		clientFolder := rootF.GetByPath(clientPath)
		if clientFolder == nil {
			continue
		}

		for _, cF := range clientFolder.Inner {
			if cF.Name == depName {
				return true, nil
			}
		}
	}

	return false, nil
}

func containsDependency(dataSources matreshka.DataSources, resource resources.Resource) bool {
	for _, ds := range dataSources {
		if ds.GetName() == resource.GetName() {
			return true
		}
	}

	return false
}
