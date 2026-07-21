package version

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed version.yaml
	versionConfig []byte
	version       string
)

//nolint:gochecknoinits // one-time parse of embedded version.yaml into the package-level version string
func init() {
	versionsMap := map[string]map[string]string{}

	err := yaml.Unmarshal(versionConfig, versionsMap)
	if err != nil {
		panic("error parsing version config" + err.Error())
	}

	version = versionsMap["app_info"]["version"]
}

func GetVersion() string {
	return version
}
