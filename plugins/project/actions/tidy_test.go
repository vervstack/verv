package actions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/git"
)

func Test_GetTidyActionsForProject_Fast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fast    bool
		wantGit bool
	}{
		{name: "normal run installs hooks and commits", fast: false, wantGit: true},
		{name: "fast run skips hooks and commit", fast: true, wantGit: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			acts := GetTidyActionsForProject(project.TypeGo, tt.fast)
			require.NotEmpty(t, acts)

			hasInstallHooks, hasCommit := false, false

			for _, a := range acts {
				switch a.(type) {
				case git.InstallHooksAction:
					hasInstallHooks = true
				case git.CommitWithUntrackedAction:
					hasCommit = true
				}
			}

			require.Equal(t, tt.wantGit, hasInstallHooks)
			require.Equal(t, tt.wantGit, hasCommit)
		})
	}
}

func Test_GetTidyActionsForProject_UnknownType(t *testing.T) {
	t.Parallel()

	acts := GetTidyActionsForProject(project.Type("unknown"), false)
	require.Nil(t, acts)
}
