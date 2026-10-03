package setup

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
)

const stepDone = "✅"

var errEmptyPassword = rerrors.New("password must not be empty")

// promptPassword asks for the password with masked echo. aborted is true if
// the user aborted (Ctrl+C/Esc) — a quiet abort rather than an error.
func promptPassword(printer io.IO, userName string) (password string, aborted bool, err error) {
	title := "Password for " + userName

	input := huh.NewInput().
		Title(title).
		EchoMode(huh.EchoModePassword).
		Validate(validatePassword).
		Value(&password)

	form := huh.NewForm(huh.NewGroup(input)).WithShowHelp(false)

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", true, nil
		}

		return "", false, rerrors.Wrap(err, "error running password prompt")
	}

	printer.PrintlnColored(colors.ColorGreen, fmt.Sprintf("%s %s: ****", stepDone, title))

	return password, false, nil
}

func validatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return rerrors.Wrap(errEmptyPassword)
	}

	return nil
}

// confirmDeploy asks whether to deploy a Velez node right away, defaulting to
// no. aborted is true if the user aborted (Ctrl+C/Esc).
func confirmDeploy() (isConfirmed bool, aborted bool, err error) {
	confirm := huh.NewConfirm().
		Title("Deploy a Velez node now?").
		Affirmative("Yes").
		Negative("No").
		Value(&isConfirmed)

	form := huh.NewForm(huh.NewGroup(confirm)).WithShowHelp(false)

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, true, nil
		}

		return false, false, rerrors.Wrap(err, "error running deploy confirmation")
	}

	return isConfirmed, false, nil
}
