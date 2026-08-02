package add

import (
	"strings"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/plugins/project/actions"
	"go.vervstack.ru/verv/plugins/project/actions/git"
	"go.vervstack.ru/verv/plugins/project/actions/go_actions/dependencies"
)

type Proc struct {
	processor.Processor

	ActionPerformer actions.IActionPerformer
}

func NewCommand(basicProc processor.Processor) *cobra.Command {
	proc := &Proc{
		Processor:       basicProc,
		ActionPerformer: actions.NewActionPerformer(basicProc.IO),
	}
	c := &cobra.Command{
		Use:   "add",
		Short: "Adds resource dependency to project",
		Long:  `Can be used to add a datasource or external API dependency to project`,

		RunE: proc.run,

		Annotations: map[string]string{"verv:emoji": "➕", "verv:requiresProject": "true"},

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(
		processor.PathFlag,
		processor.PathFlag[:1],
		proc.WD,
		`path to folder with project`)

	c.Flags().BoolP(
		processor.FastFlag, "f", false,
		`skip git init/hooks/commit steps`)

	c.Flags().BoolP(
		processor.DirtyFlag, "d", false,
		`skip git commit after codegen (repo is still initialized/hooks installed)`)

	return c
}

func (p *Proc) run(cmd *cobra.Command, args []string) error {
	project, err := p.LoadProject(cmd)
	if err != nil {
		return rerrors.Wrap(err, "error loading project for add action")
	}

	p.IO.Println(preparingMsg)

	deps := dependencies.GetDependencies(p.VervConfig, args)
	if len(deps) == 0 {
		missingDeps := dependencies.HelpWithDependencyNames()
		p.IO.Println("Unknown dependencies. Try use these: [" + strings.Join(missingDeps, ",") + "]")

		return nil
	}

	for _, d := range deps {
		err = d.AppendToProject(project)
		if err != nil {
			return rerrors.Wrap(err, "error adding dependency to project")
		}
	}

	p.IO.Println(startingMsg)

	fast, err := cmd.Flags().GetBool(processor.FastFlag)
	if err != nil {
		return rerrors.Wrap(err, "error reading fast flag")
	}

	dirty, err := cmd.Flags().GetBool(processor.DirtyFlag)
	if err != nil {
		return rerrors.Wrap(err, "error reading dirty flag")
	}

	err = p.ActionPerformer.Tidy(project, fast, dirty)
	if err != nil {
		return rerrors.Wrap(err, "error tidying project")
	}

	p.IO.Println(endMsg)

	if !fast && !dirty {
		err = git.CommitWithUntracked(project.GetProjectPath(), "added "+strings.Join(args, "; "))
		if err != nil {
			return rerrors.Wrap(err, "error performing git commit")
		}
	}

	return nil
}
