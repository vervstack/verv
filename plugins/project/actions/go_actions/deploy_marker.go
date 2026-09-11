package go_actions

import (
	"path"

	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
)

type PrepareDeployMarker struct{}

func (a PrepareDeployMarker) Do(p project.IProject) error {
	content, err := yaml.Marshal(project.Vervonomicon{Name: p.GetName()})
	if err != nil {
		return rerrors.Wrap(err, "error marshalling deploy vervonomicon marker")
	}

	p.GetFolder().Add(&folder.Folder{
		Name:    path.Join(project.VervMarkerDir, project.VervDeployDir, project.VervonomiconFile),
		Content: content,
	})

	return nil
}

func (a PrepareDeployMarker) NameInAction() string {
	return "Writing deploy marker"
}
