package dependencies

import (
	"path"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka/server"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/transport_generators"
)

const defaultServerPort = 80

func prepareServerConfig(proj Project) *server.Server {
	for _, s := range proj.GetConfig().Servers {
		return s
	}

	newServer := &server.Server{
		Name: "",
		GRPC: make(map[string]*server.GRPC),
		FS:   make(map[string]*server.FS),
		HTTP: make(map[string]*server.HTTP),
	}
	proj.GetConfig().Servers[defaultServerPort] = newServer

	return newServer
}

func initServerManagerFiles(proj Project) error {
	serverManagerPath := []string{
		patterns.InternalFolder,
		patterns.TransportFolder,
		patterns.ServerManagerFileName,
	}

	if proj.GetFolder().GetByPath(serverManagerPath...) == nil {
		content, err := transport_generators.GenerateServerManager()
		if err != nil {
			return rerrors.Wrap(err, "error generating server manager")
		}

		proj.GetFolder().Add(&folder.Folder{
			Name:    path.Join(serverManagerPath...),
			Content: content,
		})
	}

	return nil
}
