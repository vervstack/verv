package pipeline

import (
	"context"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/plugins/deployment/imagebuild"
	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
)

type Node interface {
	Version(ctx context.Context) (nodeclient.Platform, error)
	ListRegistryNames(ctx context.Context) ([]string, error)
	RegistryCredentials(ctx context.Context, name string) (nodeclient.RegistryCredentials, error)
	CreateDeploy(ctx context.Context, req nodeclient.DeployRequest) error
}

type Project interface {
	ServiceName(projectDir string) (string, error)
	Tag(ctx context.Context, projectDir string) (string, error)
}

type Builder interface {
	Build(ctx context.Context, req imagebuild.Request) error
}

type Deps struct {
	Node    Node
	Project Project
	Builder Builder
	Printer io.IO
}

type Params struct {
	ProjectDir  string
	Dockerfile  string
	Environment string
	Tag         string
}

type Status struct {
	Platform          nodeclient.Platform
	IsRegistryEnabled bool
	Err               error
}
