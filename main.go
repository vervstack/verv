package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.redsock.ru/toolbox/closer"

	addProject "go.vervstack.ru/verv/cmd/project/add"
	initProject "go.vervstack.ru/verv/cmd/project/init"
	tidyProject "go.vervstack.ru/verv/cmd/project/tidy"
	"go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/version"
)

func main() {
	defer func() {
		err := closer.Close()
		if err != nil {
			log.Err(err).Msg("errors on closing background closers")
		}
	}()

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

	basicProc := processor.New()

	root.AddCommand(initProject.NewCommand(basicProc))
	root.AddCommand(tidyProject.NewCommand(basicProc))
	root.AddCommand(addProject.NewCommand(basicProc))

	err := root.Execute()
	if err != nil {
		io.StdIO{}.Error(colors.TerminalColor(colors.ColorRed) + fmt.Sprintf("%+v\n", err))
	}
}
