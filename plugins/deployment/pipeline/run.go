package pipeline

import (
	"context"
	"fmt"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/plugins/deployment/imagebuild"
	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
)

func Run(ctx context.Context, deps Deps, params Params) error {
	if params.Environment == "" {
		return rerrors.Wrap(ErrEnvironmentRequired)
	}

	platform, err := deps.Node.Version(ctx)
	if err != nil {
		return rerrors.Wrap(err, "error reading node version")
	}

	if platform.Arch == "" {
		return rerrors.Wrap(ErrNodePlatformUnknown)
	}

	imagePlatform := "linux/" + platform.Arch
	printStep(deps, "Node platform", imagePlatform)

	names, err := deps.Node.ListRegistryNames(ctx)
	if err != nil {
		return rerrors.Wrap(err, "error listing registries")
	}

	if len(names) == 0 {
		return rerrors.Wrap(ErrNoRegistry)
	}

	creds, err := deps.Node.RegistryCredentials(ctx, names[0])
	if err != nil {
		return rerrors.Wrap(err, "error getting registry credentials")
	}

	printStep(deps, "Registry", creds.RegistryUrl)

	serviceName, err := deps.Project.ServiceName(params.ProjectDir)
	if err != nil {
		return rerrors.Wrap(err, "error resolving service name")
	}

	tag := params.Tag
	if tag == "" {
		tag, err = deps.Project.Tag(ctx, params.ProjectDir)
		if err != nil {
			return rerrors.Wrap(err, "error resolving tag")
		}
	}

	imageRef := imagebuild.ImageRef(creds.RegistryUrl, serviceName, tag)

	buildRequest := imagebuild.Request{
		ProjectDir: params.ProjectDir,
		Dockerfile: params.Dockerfile,
		Platform:   imagePlatform,
		ImageRef:   imageRef,
		Credentials: imagebuild.Credentials{
			Username:    creds.Username,
			Password:    creds.Password,
			RegistryUrl: creds.RegistryUrl,
		},
	}

	err = deps.Builder.Build(ctx, buildRequest)
	if err != nil {
		return rerrors.Wrap(err, "error building image")
	}

	printStep(deps, "Image", imageRef)

	deployRequest := nodeclient.DeployRequest{
		ServiceName: serviceName,
		Environment: params.Environment,
		Image:       imageRef,
	}

	err = deps.Node.CreateDeploy(ctx, deployRequest)
	if err != nil {
		return rerrors.Wrap(err, "error creating deploy")
	}

	deps.Printer.PrintlnColored(colors.ColorGreen,
		fmt.Sprintf("✅ Deployed %s to %s", serviceName, params.Environment))

	return nil
}

func printStep(deps Deps, what string, value string) {
	deps.Printer.PrintlnColored(colors.ColorGreen, fmt.Sprintf("✅ %s: %s", what, value))
}
