package git

import (
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
	"go.vervstack.ru/verv/plugins/project"
)

// commitTimeout is longer than the default command timeout because `git commit`
// runs the project's pre-commit hook, which lints and tests the whole project.
const commitTimeout = 5 * time.Minute

type CommitWithUntrackedAction struct {
}

func (a CommitWithUntrackedAction) Do(p project.IProject) error {
	err := CommitWithUntracked(p.GetProjectPath(), "verv auto-commit")
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
		Timeout: commitTimeout,
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
