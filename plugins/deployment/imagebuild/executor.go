package imagebuild

import (
	"context"
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

type Command struct {
	Tool    string
	Args    []string
	WorkDir string
	Env     []string
	Timeout time.Duration
}

type Executor interface {
	Run(ctx context.Context, command Command) (output string, err error)
}

type executor struct{}

func NewExecutor() Executor {
	return executor{}
}

// cmd.Execute returns empty output on failure and carries stderr/stdout only inside the
// error text, so on failure that text is handed back as the output for the failure tail.
func (executor) Run(ctx context.Context, command Command) (string, error) {
	err := ctx.Err()
	if err != nil {
		return "", rerrors.Wrap(err, "error before running command")
	}

	request := cmd.Request{
		Tool:    command.Tool,
		Args:    command.Args,
		WorkDir: command.WorkDir,
		Timeout: command.Timeout,
		Env:     command.Env,
	}

	output, err := cmd.Execute(request)
	if err != nil {
		return err.Error(), rerrors.Wrap(errCommandFailed)
	}

	return output, nil
}
