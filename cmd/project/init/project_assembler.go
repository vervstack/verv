package init

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions"
)

func (p *Proc) createProject(args project.CreateArgs, fast, dirty, deploy bool) (*project.Project, error) {
	proj, err := project.CreateProject(args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error during project creation")
	}

	initActions := actions.InitProject(project.TypeGo, fast, dirty, deploy)

	err = actions.RunPipeline(p.IO, proj, initActions, actions.PipelineLabels{
		StartEmoji: "🏗️",
		StartVerb:  "Constructing",
		EndEmoji:   "🎉",
		EndVerb:    "is ready",
	})
	if err != nil {
		return nil, rerrors.Wrap(err, "error performing init actions")
	}

	return proj, nil
}
