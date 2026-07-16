package go_actions

import (
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project/config"
)

type Action interface {
	Do(p Project) error
	NameInAction() string
}

type Project interface {
	GetName() string
	GetShortName() string

	GetConfig() *config.Config

	GetFolder() *folder.Folder
	GetProjectPath() string
}
