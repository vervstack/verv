package project

import (
	"github.com/spf13/cobra"

	"go.vervstack.ru/verv/cmd/project/add"
	projinit "go.vervstack.ru/verv/cmd/project/init"
	"go.vervstack.ru/verv/cmd/project/tidy"
	"go.vervstack.ru/verv/internal/processor"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Handles project",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	basicProc := processor.New()

	cmd.AddCommand(projinit.NewCommand(basicProc))
	cmd.AddCommand(add.NewCommand(basicProc))
	cmd.AddCommand(tidy.NewCommand(basicProc))

	return cmd
}
