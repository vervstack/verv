package go_actions

import (
	"path"

	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
)

// PrepareDeployFolder scaffolds .verv/deploy — the folder future deploy
// tooling (e.g. verv rp update) reads and writes its own state in. For now it
// only seeds a basic vervonomicon file with the service name.
type PrepareDeployFolder struct{}

func (a PrepareDeployFolder) Do(p project.IProject) error {
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

func (a PrepareDeployFolder) NameInAction() string {
	return "Preparing deploy folder"
}
