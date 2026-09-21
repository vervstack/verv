package deploy

import (
	"errors"

	"github.com/charmbracelet/huh"
	"go.redsock.ru/rerrors"
)

// selectVersion renders an arrow-key-navigable picker over the given Velez
// image tags and returns the chosen one. aborted is true if the user
// aborted (Ctrl+C/Esc) — a quiet abort rather than an error.
func selectVersion(tags []string) (version string, aborted bool, err error) {
	options := make([]huh.Option[string], 0, len(tags))
	for _, t := range tags {
		options = append(options, huh.NewOption(t, t))
	}

	var chosen string

	sel := huh.NewSelect[string]().
		Title("Which Velez version would you like to deploy?").
		Options(options...).
		Value(&chosen)

	form := huh.NewForm(huh.NewGroup(sel)).WithShowHelp(false)

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", true, nil
		}

		return "", false, rerrors.Wrap(err, "error running version picker")
	}

	return chosen, false, nil
}

// promptWithDefault shows a pre-filled, editable text prompt for title —
// the user accepts def by pressing enter, or edits it first. aborted is
// true if the user aborted (Ctrl+C/Esc) — a quiet abort rather than an
// error.
func promptWithDefault(title, def string) (value string, aborted bool, err error) {
	value = def

	input := huh.NewInput().
		Title(title).
		Value(&value)

	form := huh.NewForm(huh.NewGroup(input)).WithShowHelp(false)

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", true, nil
		}

		return "", false, rerrors.Wrap(err, "error running input prompt")
	}

	return value, false, nil
}
