package actions

import (
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/git"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions"
)

func InitProject(pt project.Type, fast, dirty, deploy bool) []Action {
	switch pt {
	case project.TypeGo:
		return initVirtualGoProject(fast, dirty, deploy)
	default:
		return nil
	}
}

// initVirtualGoProject builds the init pipeline. When fast is true, git.InitGit
// is skipped — it contributes nothing to whether the generated project compiles,
// and callers that only want to validate codegen shouldn't pay for it. When
// dirty is true (and fast is false), git.InitGit still runs but skips its final
// commit so the caller can review/amend the generated changes before
// committing manually. When deploy is true, a generic deploy marker is written
// under .verv/deploy for future deploy tooling to key off of.
func initVirtualGoProject(fast, dirty, deploy bool) []Action {
	acts := []Action{
		go_actions.PrepareProjectStructure{}, // basic go project structure
		go_actions.InitGoProjectApp{},
		go_actions.GenerateProjectConfig{},
		go_actions.PrepareConfigFolder{}, // generates config keys
		go_actions.PrepareClients{},
		go_actions.PrepareServer{},
		go_actions.PrepareGitHooks{},
		go_actions.PrepareVervMarker{},
	}

	if deploy {
		acts = append(acts, go_actions.PrepareDeployMarker{})
	}

	acts = append(acts,
		go_actions.BuildProjectAction{}, // build project in file system

		go_actions.InitGoMod{}, // executes go mod

		go_actions.BuildProjectAction{}, // builds project to file system

		go_actions.RunGoTidyAction{}, // resolves real dependencies and formats go code
	)

	if !fast {
		acts = append(acts, git.InitGit{SkipCommit: dirty})
	}

	return acts
}
