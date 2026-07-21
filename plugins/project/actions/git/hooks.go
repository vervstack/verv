package git

import (
	"os"
	"path/filepath"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

type InstallHooksAction struct {
}

func (a InstallHooksAction) Do(p project.IProject) error {
	return InstallHooks(p.GetProjectPath())
}

func (a InstallHooksAction) NameInAction() string {
	return "Installing git hooks"
}

func InstallHooks(workDir string) error {
	hookPath := filepath.Join(workDir, patterns.GitHooksFolder, patterns.PreCommitHook.Name)

	err := os.Chmod(hookPath, 0o755)
	if err != nil {
		return rerrors.Wrap(err, "error making pre-commit hook executable")
	}

	_, err = cmd.Execute(cmd.Request{
		Tool:    bin,
		Args:    []string{"config", "core.hooksPath", patterns.GitHooksFolder},
		WorkDir: workDir,
	})
	if err != nil {
		return rerrors.Wrap(err, "error configuring git hooks path")
	}

	return nil
}
