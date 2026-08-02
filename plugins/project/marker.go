package project

import (
	"os"
	"path"

	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"
)

const (
	VervMarkerDir    = ".verv"
	VervonomiconFile = "vervonomicon.yaml"
)

type Vervonomicon struct {
	Name string `yaml:"name"`
}

// IsVervProject reports whether projectPath contains the .verv/vervonomicon.yaml
// marker written by init/tidy. This is the only signal used to decide whether
// a directory is a verv-generated project — no other heuristics are applied.
func IsVervProject(projectPath string) bool {
	_, err := os.Stat(path.Join(projectPath, VervMarkerDir, VervonomiconFile))

	return err == nil
}

func ReadVervonomicon(projectPath string) (*Vervonomicon, error) {
	data, err := os.ReadFile(path.Join(projectPath, VervMarkerDir, VervonomiconFile))
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading vervonomicon marker")
	}

	v := &Vervonomicon{}

	err = yaml.Unmarshal(data, v)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vervonomicon marker")
	}

	return v, nil
}
