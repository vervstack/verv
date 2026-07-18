package main

import (
	"fmt"

	"github.com/spf13/cobra"

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
	go install go.vervstack.ru/verv@` + newVersion + `
`)
	}

	root := &cobra.Command{
		Use: "verv [command] [arguments] [flags]",

		Short: "Verv CLI is a tool for handling verv projects",

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

	err := root.Execute()
	if err != nil {
		io.StdIO{}.Error(colors.TerminalColor(colors.ColorRed) + fmt.Sprintf("%+v\n", err))
	}
}
