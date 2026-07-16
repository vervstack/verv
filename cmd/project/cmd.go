package project

import (
	"github.com/spf13/cobra"

	"go.vervstack.ru/verv/cmd/project/add"
	"go.vervstack.ru/verv/cmd/project/init_new"
	"go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/processor"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Handles project",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	stdIO := io.StdIO{}
	wd := io.GetWd()
	cfg := config.GetConfig()

	basicProc := processor.New()

	cmd.AddCommand(init_new.NewCommand(basicProc))
	cmd.AddCommand(add.NewCommand(basicProc))

	cmd.AddCommand(newLinkCmd(projectLink{
		io:     stdIO,
		path:   wd,
		config: cfg,
	}))

	cmd.AddCommand(newTidyCmd(projectTidy{
		io:     stdIO,
		path:   wd,
		config: cfg,
	}))

	return cmd
}
