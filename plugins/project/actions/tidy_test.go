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
		name       string
		fast       bool
		dirty      bool
		wantHooks  bool
		wantCommit bool
	}{
		{name: "normal run installs hooks and commits", fast: false, dirty: false, wantHooks: true, wantCommit: true},
		{name: "dirty run installs hooks but skips commit", fast: false, dirty: true, wantHooks: true, wantCommit: false},
		{name: "fast run skips hooks and commit", fast: true, dirty: false, wantHooks: false, wantCommit: false},
		{name: "fast+dirty run skips hooks and commit", fast: true, dirty: true, wantHooks: false, wantCommit: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			acts := GetTidyActionsForProject(project.TypeGo, tt.fast, tt.dirty)
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

			require.Equal(t, tt.wantHooks, hasInstallHooks)
			require.Equal(t, tt.wantCommit, hasCommit)
		})
	}
}

func Test_GetTidyActionsForProject_UnknownType(t *testing.T) {
	t.Parallel()

	acts := GetTidyActionsForProject(project.Type("unknown"), false, false)
	require.Nil(t, acts)
}
