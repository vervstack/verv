package project

//go:generate minimock -i IProject -o ./../../tests/mocks -g -s "_mock.go"

import (
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project/config"
)

const (
	TypeUnknown Type = "Unknown"
	TypeGo      Type = "go"
)

type Type string

type IProject interface {
	GetName() string
	GetShortName() string

	GetConfig() *config.Config

	GetFolder() *folder.Folder
	GetProjectPath() string

	GetType() Type
}
