package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.redsock.ru/toolbox/closer"

	addProject "go.vervstack.ru/verv/cmd/project/add"
	initProject "go.vervstack.ru/verv/cmd/project/init"
	tidyProject "go.vervstack.ru/verv/cmd/project/tidy"
	"go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/internal/processor"
	"go.vervstack.ru/verv/version"
)

func main() {
	os.Exit(run())
}

// run returns the process exit code. It's a separate function from main so
// that main's deferred cleanup (closer.Close) still executes before the
// process exits — os.Exit does not run deferred calls.
func run() int {
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

	if len(os.Args) == 1 {
		code, exit := runCommandMenu(root, basicProc)
		if exit {
			return code
		}
	}

	err := root.Execute()
	if err != nil {
		io.StdIO{}.Error(colors.TerminalColor(colors.ColorRed) + fmt.Sprintf("%+v\n", err))

		return 1
	}

	return 0
}

// runCommandMenu drives the interactive command picker shown for a bare
// `verv` invocation. It returns the exit code to use and whether run should
// return immediately with it; when exit is false, root has been prepared
// (via SetArgs) for the caller's normal root.Execute() call.
func runCommandMenu(root *cobra.Command, basicProc processor.Processor) (code int, exit bool) {
	header := menu.BuildHeader(basicProc.WD, basicProc.VervConfig)

	entry, err := menu.Select(menu.BuildEntries(root.Commands()), header)
	if err != nil {
		io.StdIO{}.Error(colors.TerminalColor(colors.ColorRed) + fmt.Sprintf("%+v\n", err))

		return 1, true
	}

	if entry == nil {
		return 0, true
	}

	root.SetArgs([]string{entry.Name})

	return 0, false
}
