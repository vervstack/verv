package deploy

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/io/colors"
)

const (
	// stepDone is the mark huh's own field UI collapses into once a step
	// completes — bubbletea's inline renderer erases the field's whole render
	// region on quit (see standard_renderer.go's EraseScreenBelow on a
	// shrinking frame), so without this the completed step would leave no
	// trace at all. Printing it after Run() is what makes each step stack
	// downward instead — the same one-line-per-step ledger `npx create-*`
	// installers leave behind.
	stepDone = "✅"
)

// selectVersion renders an arrow-key-navigable picker over the given Velez
// image tags and returns the chosen one. aborted is true if the user
// aborted (Ctrl+C/Esc) — a quiet abort rather than an error.
func selectVersion(printer io.IO, tags []string) (version string, aborted bool, err error) {
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

	printer.PrintlnColored(colors.ColorGreen, fmt.Sprintf("%s Velez version: %s", stepDone, chosen))

	return chosen, false, nil
}

// promptWithDefault shows a pre-filled, editable text prompt for title —
// the user accepts def by pressing enter, or edits it first. aborted is
// true if the user aborted (Ctrl+C/Esc) — a quiet abort rather than an
// error.
func promptWithDefault(printer io.IO, title, def string) (value string, aborted bool, err error) {
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

	printer.PrintlnColored(colors.ColorGreen, fmt.Sprintf("%s %s: %s", stepDone, title, value))

	return value, false, nil
}
