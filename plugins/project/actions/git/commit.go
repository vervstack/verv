package git

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
	"go.vervstack.ru/verv/plugins/project"
)

type CommitWithUntrackedAction struct {
}

func (a CommitWithUntrackedAction) Do(p project.IProject) error {
	err := CommitWithUntracked(p.GetProjectPath(), "rscli auto-commit")
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

func (a CommitWithUntrackedAction) NameInAction() string {
	return "Committing changes"
}

func Commit(workingDir, msg string) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    bin,
		Args:    []string{"commit", "-m", "\"" + msg + "\""},
		WorkDir: workingDir,
	})
	if err != nil {
		return rerrors.Wrap(err, "error committing files to git repository")
	}

	return nil
}

func CommitWithUntracked(workDir, msg string) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    bin,
		Args:    []string{"add", "."},
		WorkDir: workDir,
	})
	if err != nil {
		return rerrors.Wrap(err, "error adding files to git repository")
	}

	status, err := Status(workDir)
	if err != nil {
		return rerrors.Wrap(err, "error getting git status")
	}
	if len(status) == 0 {
		return nil
	}

	err = Commit(workDir, msg)
	if err != nil {
		return rerrors.Wrap(err, "error performing commit")
	}

	return nil
}
