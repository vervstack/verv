package patterns

import (
	_ "embed"

	"go.vervstack.ru/verv/internal/io/folder"
)

// GitHub Workflows
var (
	//go:embed pattern_c/.github/workflows/release.yaml
	githubWorkflowRelease []byte
	GithubWorkflowRelease = &folder.Folder{
		Name:    "release.yaml",
		Content: githubWorkflowRelease,
	}

	//go:embed pattern_c/.github/workflows/go-branch-push.yml
	githubWorkflowGoBranchPush []byte
	GithubWorkflowGoBranchPush = &folder.Folder{
		Name:    "branch-push.yaml",
		Content: githubWorkflowGoBranchPush,
	}
)
