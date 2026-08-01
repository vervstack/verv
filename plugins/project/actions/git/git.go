package git

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/plugins/project"
)

const (
	ChangesTypeInvalid = iota
	ChangesTypeNotStaged
	ChangesTypeNotCommitted
)

type InitGit struct {
	// SkipCommit, when true, skips the initial commit while still
	// initializing the git repo and installing hooks.
	SkipCommit bool
}

func (a InitGit) Do(p project.IProject) error {
	projectPath := p.GetProjectPath()

	err := Init(projectPath)
	if err != nil {
		return rerrors.Wrap(err, "error initializing project")
	}

	err = InstallHooks(projectPath)
	if err != nil {
		return rerrors.Wrap(err, "error installing git hooks")
	}

	if !a.SkipCommit {
		err = CommitWithUntracked(projectPath, "project init via RedSock CLI")
		if err != nil {
			return rerrors.Wrap(err, "error committing changes")
		}
	}

	err = SetOrigin(projectPath, p.GetName())
	if err != nil {
		return rerrors.Wrap(err, "error setting git origin")
	}

	return nil
}
func (a InitGit) NameInAction() string {
	return "Initiating git"
}
