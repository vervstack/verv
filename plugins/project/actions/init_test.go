package actions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/git"
)

func Test_InitProject_Fast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fast    bool
		wantGit bool
	}{
		{name: "normal run includes git init", fast: false, wantGit: true},
		{name: "fast run skips git init", fast: true, wantGit: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			acts := InitProject(project.TypeGo, tt.fast)
			require.NotEmpty(t, acts)

			hasGitInit := false

			for _, a := range acts {
				if _, ok := a.(git.InitGit); ok {
					hasGitInit = true
				}
			}

			require.Equal(t, tt.wantGit, hasGitInit)
		})
	}
}

func Test_InitProject_UnknownType(t *testing.T) {
	t.Parallel()

	acts := InitProject(project.Type("unknown"), false)
	require.Nil(t, acts)
}
