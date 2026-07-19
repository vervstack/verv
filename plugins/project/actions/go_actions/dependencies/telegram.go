package dependencies

import (
	"path"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions/renamer"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/clients_generators"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/transport_generators"
)

type Telegram struct {
	dependencyBase
}

func telegram(dep dependencyBase) Dependency {
	return &Telegram{
		dep,
	}
}

func (t Telegram) GetFolderName() string {
	if t.Name != "" {
		return t.Name
	}

	return patterns.TelegramServer
}

func (t Telegram) AppendToProject(proj Project) error {
	err := t.applyClient(proj)
	if err != nil {
		return rerrors.Wrap(err, "error applying tg client")
	}

	err = t.applyFolder(proj)
	if err != nil {
		return rerrors.Wrap(err, "error applying tg folder")
	}

	t.applyConfig(proj)

	return nil
}

func (t Telegram) applyClient(proj Project) error {
	ok, err := containsDependencyFolder(t.Cfg.Env.PathsToClients, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return rerrors.Wrap(err, "error finding Dependency path")
	}

	if ok {
		return nil
	}

	content, err := clients_generators.GenerateTelegramConn()
	if err != nil {
		return rerrors.Wrap(err, "error generating telegram conn file")
	}

	tgConnFile := &folder.Folder{
		Name:    path.Join(t.Cfg.Env.PathsToClients[0], t.GetFolderName(), patterns.ConnFileName),
		Content: content,
	}

	renamer.ReplaceProjectName(proj.GetName(), tgConnFile)

	proj.GetFolder().Add(
		tgConnFile,
	)

	return nil
}

func (t Telegram) applyFolder(proj Project) error {
	ok, err := containsDependencyFolder(t.Cfg.Env.PathToServers, proj.GetFolder(), t.GetFolderName())
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	tgServerContent, err := transport_generators.GenerateTelegramListener()
	if err != nil {
		return rerrors.Wrap(err, "error generating telegram listener")
	}
	tgServer := &folder.Folder{Name: patterns.TelegramListenerFileName, Content: tgServerContent}
	renamer.ReplaceProjectName(proj.GetName(), tgServer)

	tgHandlerContent, err := transport_generators.GenerateTelegramVersionHandler()
	if err != nil {
		return rerrors.Wrap(err, "error generating telegram version handler")
	}
	tgHandlerExample := &folder.Folder{Name: patterns.TelegramVersionHandlerFileName, Content: tgHandlerContent}
	renamer.ReplaceProjectName(proj.GetName(), tgHandlerExample)

	proj.GetFolder().Add(
		&folder.Folder{
			Name: path.Join(t.Cfg.Env.PathToServers[0], t.GetFolderName()),
			Inner: []*folder.Folder{
				tgServer,
				{
					Name:  path.Join(patterns.HandlersFolderName, patterns.VersionFolderName),
					Inner: []*folder.Folder{tgHandlerExample},
				},
			},
		},
	)

	return nil
}

func (t Telegram) applyConfig(proj Project) {
	for _, srv := range proj.GetConfig().DataSources {
		if srv.GetName() == t.GetFolderName() {
			return
		}
	}

	proj.GetConfig().DataSources = append(proj.GetConfig().DataSources, &resources.Telegram{
		Name: resources.Name(t.GetFolderName()),
	})
}
