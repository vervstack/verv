package dependencies

import (
	"path"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions/renamer"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/clients_generators"
)

type Redis struct {
	dependencyBase
}

func redisClient(dep dependencyBase) Dependency {
	return &Redis{
		dep,
	}
}

func (p Redis) GetFolderName() string {
	if p.Name != "" {
		return p.Name
	}

	return DependencyNameRedis
}

func (p Redis) AppendToProject(proj Project) error {
	err := p.applyClientFolder(proj)
	if err != nil {
		return rerrors.Wrap(err, "error applying client folder")
	}

	p.applyConfig(proj)

	return nil
}

func (p Redis) applyClientFolder(proj Project) error {
	ok, err := containsDependencyFolder(p.Cfg.Env.PathsToClients, proj.GetFolder(), p.GetFolderName())
	if err != nil {
		return rerrors.Wrap(err, "error finding Dependency path")
	}

	if ok {
		return nil
	}

	content, err := clients_generators.GenerateRedisConn()
	if err != nil {
		return rerrors.Wrap(err, "error generating redis conn file")
	}

	redisConn := &folder.Folder{
		Name:    path.Join(p.Cfg.Env.PathsToClients[0], p.GetFolderName(), patterns.ConnFileName),
		Content: content,
	}

	renamer.ReplaceProjectName(proj.GetName(), redisConn)

	proj.GetFolder().Add(redisConn)

	return nil
}

func (p Redis) applyConfig(proj Project) {
	for _, item := range proj.GetConfig().DataSources {
		if item.GetName() == p.GetFolderName() {
			return
		}
	}

	proj.GetConfig().DataSources = append(proj.GetConfig().DataSources,
		&resources.Redis{
			Name: resources.Name(p.GetFolderName()),
			Host: "localhost",
			Port: 6379,
		})
}
