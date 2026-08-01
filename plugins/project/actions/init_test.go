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
		name           string
		fast           bool
		dirty          bool
		wantGit        bool
		wantSkipCommit bool
	}{
		{name: "normal run includes git init", fast: false, dirty: false, wantGit: true, wantSkipCommit: false},
		{name: "dirty run includes git init but skips commit", fast: false, dirty: true, wantGit: true, wantSkipCommit: true},
		{name: "fast run skips git init", fast: true, dirty: false, wantGit: false},
		{name: "fast+dirty run skips git init", fast: true, dirty: true, wantGit: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			acts := InitProject(project.TypeGo, tt.fast, tt.dirty)
			require.NotEmpty(t, acts)

			hasGitInit := false

			for _, a := range acts {
				if ig, ok := a.(git.InitGit); ok {
					hasGitInit = true

					require.Equal(t, tt.wantSkipCommit, ig.SkipCommit)
				}
			}

			require.Equal(t, tt.wantGit, hasGitInit)
		})
	}
}

func Test_InitProject_UnknownType(t *testing.T) {
	t.Parallel()

	acts := InitProject(project.Type("unknown"), false, false)
	require.Nil(t, acts)
}
