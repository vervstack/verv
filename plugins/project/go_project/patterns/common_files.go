package patterns

import (
	_ "embed"

	"go.vervstack.ru/verv/internal/io/folder"
)

const (
	EnvConfigYamlFile    = "env.yaml"
	ConfigDevYamlFile    = "dev.yaml"
	ConfigMasterYamlFile = "config.yaml"
	MakefileFile         = "Makefile"
	RscliMakefileFile    = "rscli.mk"
	DockerfileFile       = "Dockerfile"

	GenCommand           = "gen"
	GenGrpcServerCommand = "gen-server-grpc"
)

// Build and deploy
var (
	//go:embed static/.gitignore
	gitIgnore []byte
	GitIgnore = &folder.Folder{
		Name:    ".gitignore",
		Content: gitIgnore,
	}

	//go:embed static/.golangci.yaml
	linter []byte
	Linter = &folder.Folder{
		Name:    ".golangci.yaml",
		Content: linter,
	}

	// Documentation
	//go:embed static/README.md
	readme []byte
	Readme = &folder.Folder{
		Name:    "README.md",
		Content: readme,
	}

	// Scripts
	//go:embed static/rscli.mk
	rscliMK []byte
	RscliMK = &folder.Folder{
		Name:    "rscli.mk",
		Content: rscliMK,
	}

	//go:embed static/Makefile
	makefile []byte
	Makefile = &folder.Folder{
		Name:    MakefileFile,
		Content: makefile,
	}
)
