package patterns

import (
	_ "embed"

	"go.vervstack.ru/verv/internal/io/folder"
)

const GitHooksFolder = ".githooks"

var (
	//go:embed static/hooks/pre-commit
	preCommitHook []byte
	PreCommitHook = &folder.Folder{
		Name:    "pre-commit",
		Content: preCommitHook,
	}
)
