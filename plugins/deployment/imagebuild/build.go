package imagebuild

import (
	"context"
	"strings"
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
)

const (
	defaultDockerfile = "Dockerfile"
	buildTimeout      = 30 * time.Minute
	failureTailLines  = 20
)

type Credentials struct {
	Username    string
	Password    string
	RegistryUrl string
}

type Request struct {
	ProjectDir  string
	Dockerfile  string
	Platform    string
	ImageRef    string
	Credentials Credentials
}

func Build(ctx context.Context, printer io.IO, executor Executor, req Request) error {
	dockerfile := req.Dockerfile
	if dockerfile == "" {
		dockerfile = defaultDockerfile
	}

	configDir, err := prepareDockerConfig(req.Credentials)
	if err != nil {
		return rerrors.Wrap(err, "error preparing docker config")
	}
	defer removeDir(configDir)

	command := Command{
		Tool: "docker",
		Args: []string{
			"buildx", "build",
			"--platform", req.Platform,
			"-f", dockerfile,
			"-t", req.ImageRef,
			"--push",
			".",
		},
		WorkDir: req.ProjectDir,
		Env:     []string{dockerConfigEnv + "=" + configDir},
		Timeout: buildTimeout,
	}

	spinner := io.NewSpinner(printer)
	spinner.Start("Building and pushing " + req.ImageRef)

	output, err := executor.Run(ctx, command)
	if err != nil {
		spinner.Stop(false, "Failed to build and push "+req.ImageRef)
		printer.Error(lastLines(output, failureTailLines) + "\n")

		return rerrors.Wrap(errBuildFailed)
	}

	spinner.Stop(true, "Built and pushed "+req.ImageRef)

	return nil
}

func lastLines(output string, count int) string {
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) > count {
		lines = lines[len(lines)-count:]
	}

	return strings.Join(lines, "\n")
}
