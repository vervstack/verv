package projectinfo

import (
	"context"
	"strings"

	"go.redsock.ru/rerrors"
)

const dirtySuffix = "-dirty"

func Tag(ctx context.Context, executor Executor, projectDir string) (string, error) {
	revision, err := executor.Run(ctx, projectDir, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", rerrors.Wrap(ErrTagUnavailable)
	}

	tag := strings.TrimSpace(revision)
	if tag == "" {
		return "", rerrors.Wrap(ErrTagUnavailable)
	}

	status, err := executor.Run(ctx, projectDir, "git", "status", "--porcelain")
	if err != nil {
		return "", rerrors.Wrap(ErrTagUnavailable)
	}

	if strings.TrimSpace(status) != "" {
		tag += dirtySuffix
	}

	return tag, nil
}
