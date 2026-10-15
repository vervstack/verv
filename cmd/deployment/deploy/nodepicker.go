package deploy

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io/colors"
	"go.vervstack.ru/verv/internal/menu"
	"go.vervstack.ru/verv/plugins/deployment/nodes"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

const (
	stepDone = "✅"

	addNodeIndex = -1
	addNodeLabel = "＋ Add node"
)

type nodeChoiceValidator struct {
	rejectReasons []string
}

func (v nodeChoiceValidator) validate(index int) error {
	if index < 0 || index >= len(v.rejectReasons) {
		return nil
	}

	reason := v.rejectReasons[index]
	if reason == "" {
		return nil
	}

	return rejectedNodeError(reason)
}

func (p *projectDeploy) pickNode(
	entries []nodes.Entry,
	statuses []pipeline.Status,
) (index int, aborted bool, err error) {
	options := make([]huh.Option[int], 0, len(entries)+1)
	validator := nodeChoiceValidator{rejectReasons: make([]string, len(entries))}

	for idx, entry := range entries {
		label, reason := describeNode(entry, statuses[idx])

		validator.rejectReasons[idx] = reason
		options = append(options, huh.NewOption(label, idx))
	}

	options = append(options, huh.NewOption(addNodeLabel, addNodeIndex))

	var chosen int

	sel := huh.NewSelect[int]().
		Title("Which node would you like to deploy to?").
		Options(options...).
		Validate(validator.validate).
		Value(&chosen)

	form := huh.NewForm(huh.NewGroup(sel)).WithShowHelp(false).WithTheme(menu.Theme())

	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return 0, true, nil
		}

		return 0, false, rerrors.Wrap(err, "error running node picker")
	}

	if chosen != addNodeIndex {
		p.io.PrintlnColored(colors.ColorGreen, fmt.Sprintf("%s Node: %s", stepDone, entries[chosen].Name))
	}

	return chosen, false, nil
}
