package actions

import (
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/ci_cd"
	"go.vervstack.ru/verv/plugins/project/actions/git"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions"
)

func GetTidyActionsForProject(pt project.Type) []Action {
	out := commonProjectTidyPreActions()

	switch pt {
	case project.TypeGo:
		out = append(out, goProjectTidyActions()...)
	default:
		return unknownProjectActions()
	}

	out = append(out, commonProjectTidyPostActions()...)
	return out
}

func goProjectTidyActions() []Action {
	return []Action{
		go_actions.PrepareConfigFolder{},
		go_actions.PrepareMakefile{},
		go_actions.PrepareClients{},
		go_actions.PrepareServer{},
		go_actions.PrepareDockerfile{},
		go_actions.BuildProjectAction{},
		// Takes too long
		//go_actions.RunMakeGenAction{},
		go_actions.InitGoProjectApp{},

		go_actions.BuildProjectAction{},
		// Takes too long
		// go_actions.RunGoTidyAction{},
		go_actions.GoFmt{},
	}
}

func commonProjectTidyPreActions() []Action {
	return []Action{
		ci_cd.TidyGithubWorkflowAction{},
	}
}

func commonProjectTidyPostActions() []Action {
	return []Action{
		git.CommitWithUntrackedAction{},
	}
}

func unknownProjectActions() []Action {
	return nil
}
