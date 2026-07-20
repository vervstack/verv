package tidy

import (
	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/actions"
)

type projectTidy struct {
	io     io.IO
	config *config.VervConfig

	path string
}

func NewCommand(basicProc processor.Processor) *cobra.Command {
	proc := projectTidy{
		io:     basicProc.IO,
		config: basicProc.VervConfig,
		path:   basicProc.WD,
	}
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Cleans project",
		Long:  "Can be used clean project",

		RunE: proc.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(
		processor.PathFlag,
		processor.PathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *projectTidy) run(_ *cobra.Command, _ []string) error {
	proj, err := project.LoadProject(p.path, p.config)
	if err != nil {
		return rerrors.Wrap(err, "error fetching project for tidy")
	}

	ap := actions.NewActionPerformer(p.io)

	err = ap.Tidy(proj)
	if err != nil {
		return rerrors.Wrap(err, "error performing tidy")
	}

	return nil
}
