package deploy

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/plugins/deployment/projectinfo"
)

func (p *projectDeploy) resolveEnvironment(flagValue string) (environment string, aborted bool, err error) {
	if flagValue != "" {
		p.printEnvironment(flagValue)

		return flagValue, false, nil
	}

	environments, err := projectinfo.Environments(p.wd)
	if err != nil {
		return "", false, rerrors.Wrap(err, "error listing environments")
	}

	if len(environments) == 0 {
		return "", false, rerrors.Wrap(errNoEnvironments)
	}

	options := make([]huh.Option[string], 0, len(environments))
	for _, env := range environments {
		options = append(options, huh.NewOption(env, env))
	}

	sel := huh.NewSelect[string]().
		Title("Which environment would you like to deploy?").
		Options(options...).
		Value(&environment)

	form := huh.NewForm(huh.NewGroup(sel)).WithShowHelp(false).WithTheme(menu.Theme())

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", true, nil
		}

		return "", false, rerrors.Wrap(err, "error running environment picker")
	}

	p.printEnvironment(environment)

	return environment, false, nil
}

func (p *projectDeploy) printEnvironment(environment string) {
	p.io.PrintlnColored(colors.ColorGreen, fmt.Sprintf("%s Environment: %s", stepDone, environment))
}
