package actions

//go:generate minimock -i ActionPerformer -o ./../../../tests/mocks -g -s "_mock.go"

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/plugins/project"
)

type Action interface {
	Do(p project.IProject) error
	NameInAction() string
}

type IActionPerformer interface {
	Tidy(proj project.IProject) error
}

type ActionPerformer struct {
	printer io.IO
}

func NewActionPerformer(printer io.IO) *ActionPerformer {
	return &ActionPerformer{
		printer: printer,
	}
}

func (a *ActionPerformer) Tidy(proj project.IProject) error {
	acts := GetTidyActionsForProject(proj.GetType())

	for _, ac := range acts {
		a.printer.Println(ac.NameInAction())

		err := ac.Do(proj)
		if err != nil {
			return rerrors.Wrap(err)
		}
	}

	return nil
}
