package actions

import (
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/ci_cd"
	"go.vervstack.ru/verv/plugins/project/actions/git"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions"
)

// GetTidyActionsForProject builds the tidy pipeline. When fast is true, the git
// hooks-install and commit steps are skipped — they contribute nothing to
// whether the generated project compiles, and callers that only want to
// validate codegen shouldn't pay for them.
func GetTidyActionsForProject(pt project.Type, fast bool) []Action {
	out := commonProjectTidyPreActions()

	switch pt {
	case project.TypeGo:
		out = append(out, goProjectTidyActions()...)
	default:
		return unknownProjectActions()
	}

	if !fast {
		out = append(out, commonProjectTidyPostActions()...)
	}

	return out
}

func goProjectTidyActions() []Action {
	return []Action{
		go_actions.PrepareConfigFolder{},
		go_actions.PrepareClients{},
		go_actions.PrepareServer{},
		go_actions.PrepareDockerfile{},
		go_actions.PrepareGitHooks{},
		go_actions.BuildProjectAction{},
		go_actions.InitGoProjectApp{},

		go_actions.BuildProjectAction{},
		go_actions.RunGoTidyAction{},
	}
}

func commonProjectTidyPreActions() []Action {
	return []Action{
		ci_cd.TidyGithubWorkflowAction{},
	}
}

func commonProjectTidyPostActions() []Action {
	return []Action{
		git.InstallHooksAction{},
		git.CommitWithUntrackedAction{},
	}
}

func unknownProjectActions() []Action {
	return nil
}
