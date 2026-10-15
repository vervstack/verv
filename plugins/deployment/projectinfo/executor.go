package projectinfo

import (
	"context"
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

const commandTimeout = 30 * time.Second

type Executor interface {
	Run(ctx context.Context, workDir string, tool string, args ...string) (stdout string, err error)
}

type executor struct{}

func NewExecutor() Executor {
	return executor{}
}

func (executor) Run(ctx context.Context, workDir string, tool string, args ...string) (string, error) {
	err := ctx.Err()
	if err != nil {
		return "", rerrors.Wrap(err, "error before running command")
	}

	request := cmd.Request{
		Tool:    tool,
		Args:    args,
		WorkDir: workDir,
		Timeout: commandTimeout,
	}

	stdout, err := cmd.Execute(request)
	if err != nil {
		return "", rerrors.Wrap(errCommandFailed)
	}

	return stdout, nil
}
