package pipelines

import (
	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

type TidyGithubWorkflowAction struct {
}

func (a TidyGithubWorkflowAction) Do(p project.IProject) error {
	ghF := p.GetFolder().GetByPath(patterns.GithubFolder, patterns.WorkflowsFolder)
	if ghF == nil {
		ghF = &folder.Folder{
			Name: patterns.GithubFolder,
			Inner: []*folder.Folder{
				{
					Name: patterns.WorkflowsFolder,
				},
			},
		}
		p.GetFolder().Add(ghF)
		ghF = ghF.Inner[0]
	}

	if ghF.GetByPath(patterns.GithubWorkflowRelease.Name) == nil {
		ghF.Add(patterns.GithubWorkflowRelease.Copy())
	}

	switch p.GetType() {
	case project.TypeGo:
		if ghF.GetByPath(patterns.GithubWorkflowGoBranchPush.Name) == nil {
			ghF.Add(patterns.GithubWorkflowGoBranchPush.Copy())
		}
	}

	return nil
}
func (a TidyGithubWorkflowAction) NameInAction() string {
	return "Tiding github workflows"
}
