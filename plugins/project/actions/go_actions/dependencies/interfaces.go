package dependencies

import (
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/config"
)

type Project interface {
	GetName() string
	GetShortName() string

	GetConfig() *config.Config

	GetFolder() *folder.Folder
	GetProjectPath() string
	GetType() project.Type
}
