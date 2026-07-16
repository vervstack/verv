package git

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
)

func Init(workingDir string) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    bin,
		Args:    []string{"init"},
		WorkDir: workingDir,
	})
	if err != nil {
		return rerrors.Wrap(err, "error initiating git repository")
	}

	return nil
}
