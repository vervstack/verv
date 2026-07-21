package project_mock

import (
	_ "embed"
)

var (
	//go:embed basic_config.yaml
	basicConfigFile []byte
)

func BasicConfig() []byte {
	n := make([]byte, len(basicConfigFile))
	copy(n, basicConfigFile)

	return n
}
