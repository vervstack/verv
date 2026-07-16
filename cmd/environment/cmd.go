package environment

import (
	"github.com/spf13/cobra"

	"go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Handles environment",

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	stdIO := io.StdIO{}

	cmd.AddCommand(newTidyEnvCmd(stdIO, config.GetConfig()))

	return cmd
}
