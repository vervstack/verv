package deploy

import (
	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/plugins/deployment/imagebuild"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
	"go.vervstack.ru/verv/plugins/deployment/projectinfo"
)

const (
	EnvironmentFlag = "environment"
	TagFlag         = "tag"
	DockerfileFlag  = "dockerfile"

	defaultDockerfile = "Dockerfile"
)

type projectDeploy struct {
	io io.IO
	wd string
}

func NewCommand(basicProc processor.Processor) *cobra.Command {
	proc := projectDeploy{
		io: basicProc.IO,
		wd: basicProc.WD,
	}

	c := &cobra.Command{
		Use:   "deploy",
		Short: "Builds this project and deploys it to a Velez node",
		Long: "Builds the project image for the node's platform, pushes it to the node's registry " +
			"and asks the node to deploy it. Runs without prompts when CI is set or stdin is not a terminal; " +
			"then VERV_NODE_URL, VERV_NODE_SECRET and --environment are required",

		RunE: proc.run,

		Annotations: map[string]string{
			"verv:emoji":          "🚀",
			"verv:group":          menu.GroupDeployment,
			"verv:requiresDocker": "true",
		},

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(EnvironmentFlag, "e", "", "environment to deploy (.verv/<env>/ or config/<env>.yaml)")
	c.Flags().String(TagFlag, "", "image tag (default: git short sha)")
	c.Flags().String(DockerfileFlag, defaultDockerfile, "path to the Dockerfile, relative to the project")

	return c
}

func (p *projectDeploy) run(cmd *cobra.Command, _ []string) error {
	params, err := p.readParams(cmd)
	if err != nil {
		return rerrors.Wrap(err, "error reading flags")
	}

	if isCiMode() {
		err = p.runCi(cmd.Context(), params)
		if err != nil {
			return rerrors.Wrap(err, "error deploying in ci mode")
		}

		return nil
	}

	err = p.runInteractive(cmd.Context(), params)
	if err != nil {
		return rerrors.Wrap(err, "error deploying interactively")
	}

	return nil
}

func (p *projectDeploy) readParams(cmd *cobra.Command) (pipeline.Params, error) {
	environment, err := cmd.Flags().GetString(EnvironmentFlag)
	if err != nil {
		return pipeline.Params{}, rerrors.Wrap(err, "error reading environment flag")
	}

	tag, err := cmd.Flags().GetString(TagFlag)
	if err != nil {
		return pipeline.Params{}, rerrors.Wrap(err, "error reading tag flag")
	}

	dockerfile, err := cmd.Flags().GetString(DockerfileFlag)
	if err != nil {
		return pipeline.Params{}, rerrors.Wrap(err, "error reading dockerfile flag")
	}

	params := pipeline.Params{
		ProjectDir:  p.wd,
		Dockerfile:  dockerfile,
		Environment: environment,
		Tag:         tag,
	}

	return params, nil
}

func (p *projectDeploy) newDeps(node pipeline.Node) pipeline.Deps {
	return pipeline.Deps{
		Node:    node,
		Project: pipeline.NewProject(projectinfo.NewExecutor()),
		Builder: pipeline.NewBuilder(p.io, imagebuild.NewExecutor()),
		Printer: p.io,
	}
}
