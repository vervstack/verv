package actions

import (
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/git"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions"
)

func InitProject(pt project.Type, fast bool) []Action {
	switch pt {
	case project.TypeGo:
		return initVirtualGoProject(fast)
	default:
		return nil
	}
}

// initVirtualGoProject builds the init pipeline. When fast is true, git.InitGit
// is skipped — it contributes nothing to whether the generated project compiles,
// and callers that only want to validate codegen shouldn't pay for it.
func initVirtualGoProject(fast bool) []Action {
	acts := []Action{
		go_actions.PrepareProjectStructure{}, // basic go project structure
		go_actions.InitGoProjectApp{},
		go_actions.GenerateProjectConfig{},
		go_actions.PrepareConfigFolder{}, // generates config keys
		go_actions.PrepareClients{},
		go_actions.PrepareServer{},
		go_actions.PrepareGitHooks{},

		go_actions.BuildProjectAction{}, // build project in file system

		go_actions.InitGoMod{}, // executes go mod

		go_actions.BuildProjectAction{}, // builds project to file system

		go_actions.RunGoTidyAction{}, // resolves real dependencies and formats go code
	}

	if !fast {
		acts = append(acts, git.InitGit{})
	}

	return acts
}
