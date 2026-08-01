package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/internal/cmd"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/git"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

// Test_InitGit_Do_SkipCommit verifies that InitGit still initializes the repo
// and installs hooks when SkipCommit is set, but does not create a commit.
func Test_InitGit_Do_SkipCommit(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	hooksDir := filepath.Join(tmpDir, patterns.GitHooksFolder)
	require.NoError(t, os.MkdirAll(hooksDir, 0o755)) //nolint:mnd

	hookPath := filepath.Join(hooksDir, patterns.PreCommitHook.Name)
	require.NoError(t, os.WriteFile(hookPath, []byte("#!/bin/sh\nexit 0\n"), 0o644)) //nolint:mnd

	// a tracked-worthy file so there'd be something to commit if SkipCommit
	// were not honored
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("hello"), 0o644)) //nolint:mnd

	proj := &project.Project{
		Name: "test/dirty-init-project",
		Path: tmpDir,
	}

	err := git.InitGit{SkipCommit: true}.Do(proj)
	require.NoError(t, err)

	// git repo was still initialized
	require.DirExists(t, filepath.Join(tmpDir, ".git"))

	// hooks path was still configured
	hooksPath, err := cmd.Execute(cmd.Request{
		Tool:    "git",
		Args:    []string{"config", "core.hooksPath"},
		WorkDir: tmpDir,
	})
	require.NoError(t, err)
	require.Contains(t, hooksPath, patterns.GitHooksFolder)

	// no commit was created
	_, err = cmd.Execute(cmd.Request{
		Tool:    "git",
		Args:    []string{"log"},
		WorkDir: tmpDir,
	})
	require.Error(t, err)
}
