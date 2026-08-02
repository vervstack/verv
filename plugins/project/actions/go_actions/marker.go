package go_actions

import (
	"path"

	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
)

type PrepareVervMarker struct{}

func (a PrepareVervMarker) Do(p project.IProject) error {
	content, err := yaml.Marshal(project.Vervonomicon{Name: p.GetName()})
	if err != nil {
		return rerrors.Wrap(err, "error marshalling vervonomicon marker")
	}

	p.GetFolder().Add(&folder.Folder{
		Name:    path.Join(project.VervMarkerDir, project.VervonomiconFile),
		Content: content,
	})

	return nil
}

func (a PrepareVervMarker) NameInAction() string {
	return "Writing verv marker"
}
