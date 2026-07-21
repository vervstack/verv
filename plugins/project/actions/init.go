package actions

import (
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/git"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions"
)

func InitProject(pt project.Type) []Action {
	switch pt {
	case project.TypeGo:
		return initVirtualGoProject()
	default:
		return nil
	}
}

func initVirtualGoProject() []Action {
	return []Action{
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

		git.InitGit{},
	}
}
