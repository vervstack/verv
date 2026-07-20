package init

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions"
)

func (p *Proc) createProject(args project.CreateArgs) (project.IProject, error) {
	proj, err := project.CreateProject(args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error during project creation")
	}

	p.IO.Println("Starting project constructor")

	initActions := actions.InitProject(project.TypeGo)
	for _, act := range initActions {
		err = act.Do(proj)
		if err != nil {
			return nil, rerrors.Wrap(err, "error performing init actions")
		}
	}

	p.IO.Println("Project actions performed")

	return proj, nil
}
