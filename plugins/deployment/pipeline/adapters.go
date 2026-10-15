package pipeline

import (
	"context"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/plugins/deployment/imagebuild"
	"go.vervstack.ru/verv/plugins/deployment/projectinfo"
)

type project struct {
	executor projectinfo.Executor
}

func NewProject(executor projectinfo.Executor) Project {
	return project{executor: executor}
}

func (p project) ServiceName(projectDir string) (string, error) {
	name, err := projectinfo.ServiceName(projectDir)
	if err != nil {
		return "", rerrors.Wrap(err, "error reading service name")
	}

	return name, nil
}

func (p project) Tag(ctx context.Context, projectDir string) (string, error) {
	tag, err := projectinfo.Tag(ctx, p.executor, projectDir)
	if err != nil {
		return "", rerrors.Wrap(err, "error resolving tag")
	}

	return tag, nil
}

type builder struct {
	printer  io.IO
	executor imagebuild.Executor
}

func NewBuilder(printer io.IO, executor imagebuild.Executor) Builder {
	return builder{printer: printer, executor: executor}
}

func (b builder) Build(ctx context.Context, req imagebuild.Request) error {
	err := imagebuild.Build(ctx, b.printer, b.executor, req)
	if err != nil {
		return rerrors.Wrap(err, "error building image")
	}

	return nil
}
