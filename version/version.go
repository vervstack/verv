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
	m := map[string]map[string]string{}

	err := yaml.Unmarshal(versionConfig, m)
	if err != nil {
		panic("error parsing version config" + err.Error())
	}
	version = m["app_info"]["version"]
}

func GetVersion() string {
	return version
}
