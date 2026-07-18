package project

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

	proj *project.Project
	path string
}

func newTidyCmd(pl projectTidy) *cobra.Command {
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Cleans project",
		Long:  "Can be used clean project",

		RunE: pl.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(
		processor.PathFlag,
		processor.PathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *projectTidy) run(_ *cobra.Command, _ []string) (err error) {
	if p.proj == nil {
		p.proj, err = project.LoadProject(p.path, p.config)
		if err != nil {
			return rerrors.Wrap(err, "error fetching project for tidy")
		}
	}

	ap := actions.NewActionPerformer(p.io)

	err = ap.Tidy(p.proj)
	if err != nil {
		return rerrors.Wrap(err, "error performing tidy")
	}

	return nil
}
