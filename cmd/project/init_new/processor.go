package init_new

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/plugins/project"
)

const newProjectInitMessage = `New project with name %s initialized at %s`

type Proc struct {
	processor.Processor

	nameCollector nameCollector
}

func NewCommand(basicProc processor.Processor) *cobra.Command {
	proc := &Proc{
		Processor:     basicProc,
		nameCollector: newNameCollector(basicProc.IO, basicProc.RscliConfig.DefaultProjectGitPath),
	}

	c := &cobra.Command{
		Use:   "init",
		Short: "Initializes project",
		Long:  `Can be used to init a project via configuration file, constructor or global config`,
		RunE:  proc.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	return c
}

func (p *Proc) run(_ *cobra.Command, cmdArgs []string) (err error) {
	cArgs := project.CreateArgs{
		CfgPath: p.RscliConfig.Env.PathToConfig,
	}

	// step 1: obtain name
	cArgs.Name, err = p.nameCollector.collect(cmdArgs)
	if err != nil {
		return rerrors.Wrap(err, "can't obtain name")
	}

	// step 2: obtain path to project folder
	cArgs.ProjectPath = p.collectOsPath(cArgs.Name, cmdArgs)

	proj, err := p.createProject(cArgs)
	if err != nil {
		return rerrors.Wrap(err, "error building project")
	}

	p.IO.PrintlnColored(colors.ColorGreen,
		fmt.Sprintf(newProjectInitMessage, proj.GetName(), proj.GetProjectPath()))

	return nil
}
