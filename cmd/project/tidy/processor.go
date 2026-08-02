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

		Annotations: map[string]string{"verv:emoji": "🧹", "verv:requiresProject": "true"},

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(
		processor.PathFlag,
		processor.PathFlag[:1], "", `path to folder with project`)

	c.Flags().BoolP(
		processor.FastFlag, "f", false,
		`skip git init/hooks/commit steps`)

	c.Flags().BoolP(
		processor.DirtyFlag, "d", false,
		`skip git commit after codegen (repo is still initialized/hooks installed)`)

	return c
}

func (p *projectTidy) run(cmd *cobra.Command, _ []string) error {
	proj, err := project.LoadProject(p.path, p.config)
	if err != nil {
		return rerrors.Wrap(err, "error fetching project for tidy")
	}

	fast, err := cmd.Flags().GetBool(processor.FastFlag)
	if err != nil {
		return rerrors.Wrap(err, "error reading fast flag")
	}

	dirty, err := cmd.Flags().GetBool(processor.DirtyFlag)
	if err != nil {
		return rerrors.Wrap(err, "error reading dirty flag")
	}

	ap := actions.NewActionPerformer(p.io)

	err = ap.Tidy(proj, fast, dirty)
	if err != nil {
		return rerrors.Wrap(err, "error performing tidy")
	}

	return nil
}
