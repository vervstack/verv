package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.vervstack.ru/verv/cmd/environment"
	initCmd "go.vervstack.ru/verv/cmd/project"
	"go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/version"
)

func main() {
	newVersion, canUpdate := version.CanUpdate()
	if canUpdate {
		io.StdIO{}.Println(`
⚙️⚙️⚙️ Update is available ⚙️⚙️⚙️
Run this to install it:
	go install github.com/Red-Sock/rscli@` + newVersion + `
`)
	}

	root := &cobra.Command{
		Use: "rscli [command] [arguments] [flags]",

		Short: "RsCLI is a tool for handling developers environment",

		Version: version.GetVersion(),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRunE: config.InitConfig,
		SilenceErrors:     true,
		SilenceUsage:      true,
	}

	root.PersistentFlags().String(config.CustomPathToConfig, "", "path flag to custom config")

	root.AddCommand(initCmd.NewCmd())
	root.AddCommand(environment.NewCmd())

	err := root.Execute()
	if err != nil {
		io.StdIO{}.Error(colors.TerminalColor(colors.ColorRed) + fmt.Sprintf("%+v\n", err))
	}
}
