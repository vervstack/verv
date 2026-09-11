package actions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions/git"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions"
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

			acts := InitProject(project.TypeGo, tt.fast, tt.dirty, false)
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

	acts := InitProject(project.Type("unknown"), false, false, false)
	require.Nil(t, acts)
}

func Test_InitProject_Deploy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		deploy     bool
		wantDeploy bool
	}{
		{name: "deploy flag adds deploy folder", deploy: true, wantDeploy: true},
		{name: "no deploy flag skips deploy folder", deploy: false, wantDeploy: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			acts := InitProject(project.TypeGo, false, false, tt.deploy)
			require.NotEmpty(t, acts)

			hasDeployFolder := false

			for _, a := range acts {
				if _, ok := a.(go_actions.PrepareDeployFolder); ok {
					hasDeployFolder = true
				}
			}

			require.Equal(t, tt.wantDeploy, hasDeployFolder)
		})
	}
}
