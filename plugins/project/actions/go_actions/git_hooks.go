package go_actions

import (
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

type PrepareGitHooks struct {
}

func (a PrepareGitHooks) Do(p project.IProject) error {
	hooksF := p.GetFolder().GetByPath(patterns.GitHooksFolder)
	if hooksF == nil {
		hooksF = &folder.Folder{Name: patterns.GitHooksFolder}
		p.GetFolder().Add(hooksF)
	}

	if hooksF.GetByPath(patterns.PreCommitHook.Name) == nil {
		hooksF.Add(patterns.PreCommitHook.Copy())
	}

	return nil
}

func (a PrepareGitHooks) NameInAction() string {
	return "Preparing git hooks"
}
