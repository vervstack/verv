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

// promptPassword asks for the password with masked echo. For an existing
// user an empty answer is valid and means "keep the current password". aborted
// is true if the user aborted (Ctrl+C/Esc) — a quiet abort rather than an error.
func promptPassword(
	printer io.IO, userName string, isExistingUser bool,
) (password string, aborted bool, err error) {
	title := "Password for " + userName

	input := huh.NewInput().
		Title(title).
		EchoMode(huh.EchoModePassword).
		Value(&password)

	if isExistingUser {
		title = "New password for " + userName
		input.Title(title).Description("Leave empty to keep the current password")
	} else {
		input.Validate(validatePassword)
	}

	form := huh.NewForm(huh.NewGroup(input)).WithShowHelp(false)

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", true, nil
		}

		return "", false, rerrors.Wrap(err, "error running password prompt")
	}

	if password == "" {
		printer.PrintlnColored(colors.ColorGreen, fmt.Sprintf("%s %s: keeping current", stepDone, title))

		return "", false, nil
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
