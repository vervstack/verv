package dependencies

import (
	"path"

	"go.redsock.ru/rerrors"

	vervconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions/renamer"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/clients_generators"
)

type sqlConn struct {
	Cfg *vervconfig.VervConfig
}

func (sc sqlConn) GetFolderName() string {
	return "sqldb"
}

func (sc sqlConn) applySqlConnectionFile(proj Project) error {
	if len(sc.Cfg.Env.PathsToClients) == 0 {
		return ErrNoFolderInConfig
	}

	content, err := clients_generators.GenerateSQLConn()
	if err != nil {
		return rerrors.Wrap(err, "error generating sql conn file")
	}

	fileName := path.Join(
		sc.Cfg.Env.PathsToClients[0],
		sc.GetFolderName(),
		patterns.ConnFileName)

	sqlConnFile := &folder.Folder{
		Name:    fileName,
		Content: content,
	}

	renamer.ReplaceProjectName(proj.GetName(), sqlConnFile)
	proj.GetFolder().Add(sqlConnFile)

	return nil
}

func (sc sqlConn) applySqlDriver(proj Project, driverName string, content []byte) {
	fileName := path.Join(
		sc.Cfg.Env.PathsToClients[0],
		sc.GetFolderName(),
		driverName+".go")

	sqlDriverFile := &folder.Folder{
		Name:    fileName,
		Content: content,
	}

	proj.GetFolder().Add(sqlDriverFile)
}
