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

const (
	stepDone              = "✅"
	sysboxSuggestionCount = 3
)

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

// confirmSysboxInstall asks whether to install sysbox at all, defaulting to yes. aborted is true if the user
// aborted (Ctrl+C/Esc).
func confirmSysboxInstall(printer io.IO) (isConfirmed bool, aborted bool, err error) {
	isConfirmed = true

	confirm := huh.NewConfirm().
		Title("Install sysbox?").
		Affirmative("Yes").
		Negative("No").
		Value(&isConfirmed)

	form := huh.NewForm(huh.NewGroup(confirm)).WithShowHelp(false)

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, true, nil
		}

		return false, false, rerrors.Wrap(err, "error running sysbox install confirmation")
	}

	if !isConfirmed {
		printer.PrintlnColored(colors.ColorGreen, stepDone+" Install sysbox: no")
	}

	return isConfirmed, false, nil
}

// promptSysboxVersion asks for a sysbox version with the newest ones as Tab suggestions; an empty answer
// means the latest. versions is newest first and must not be empty.
func promptSysboxVersion(printer io.IO, versions []string) (version string, aborted bool, err error) {
	suggestions := versions[:min(sysboxSuggestionCount, len(versions))]

	input := huh.NewInput().
		Title("Sysbox version").
		Placeholder("empty = latest (" + versions[0] + ")").
		Description("Recent: " + strings.Join(suggestions, ", ") + " (Tab to complete)").
		Suggestions(suggestions).
		Value(&version)

	form := huh.NewForm(huh.NewGroup(input)).WithShowHelp(false)

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", true, nil
		}

		return "", false, rerrors.Wrap(err, "error running sysbox version prompt")
	}

	version = strings.TrimSpace(version)

	shown := version
	if shown == "" {
		shown = "latest"
	}

	printer.PrintlnColored(colors.ColorGreen, fmt.Sprintf("%s Sysbox version: %s", stepDone, shown))

	return version, false, nil
}

// confirmDockerRestart asks whether stopping the running containers by a Docker restart is acceptable,
// defaulting to no. aborted is true if the user aborted (Ctrl+C/Esc).
func confirmDockerRestart(containers int) (isConfirmed bool, aborted bool, err error) {
	title := fmt.Sprintf("%d container(s) are running. Installing sysbox restarts Docker and stops them. Continue?",
		containers)

	confirm := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(&isConfirmed)

	form := huh.NewForm(huh.NewGroup(confirm)).WithShowHelp(false)

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, true, nil
		}

		return false, false, rerrors.Wrap(err, "error running docker restart confirmation")
	}

	return isConfirmed, false, nil
}
