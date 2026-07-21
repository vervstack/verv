// TODO redo onto the evon parsing

package config

import (
	"database/sql"
	_ "embed"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"
)

const (
	CustomPathToConfig = "cfg"

	configFilename          = "verv.yaml"
	environmentPathToConfig = "VERV_CONFIG_PATH"

	envPathToConfig = "VERV_PATH_TO_CONFIG"
	envPathToMain   = "VERV_PATH_TO_MAIN"

	envPathToProtoClients    = "VERV_PATH_TO_PROTO_CLIENTS"
	envPathToCompiledClients = "VERV_PATH_TO_COMPILED_CLIENTS"
	envPathToClients         = "VERV_PATH_TO_CLIENTS"

	envPathToServers          = "VERV_PATH_TO_SERVERS"
	envPathToServerDefinition = "VERV_PATH_TO_SERVER_DEFINITION"

	envPathToMigrations      = "VERV_PATH_TO_MIGRATIONS"
	envDefaultProjectGitPath = "VERV_DEFAULT_PROJECT_GIT_PATH"
)

var (
	//go:embed verv.yaml
	builtInConfig []byte

	vervConfig *VervConfig
)

type VervConfig struct {
	Env                   Project `yaml:"env"`
	DefaultProjectGitPath string  `yaml:"default_project_git_path"`
}

type Project struct {
	PathToMain   string `yaml:"path_to_main"`
	PathToConfig string `yaml:"path_to_config"`

	PathsToCompiledClients []string `yaml:"paths_to_compiled_clients"`
	PathsToProtoClients    []string `yaml:"paths_to_proto_clients"`
	PathsToClients         []string `yaml:"paths_to_clients"`

	PathToServers          []string `yaml:"path_to_servers"`
	PathToServerDefinition string   `yaml:"path_to_server_definition"`

	PathToMigrations string `yaml:"path_to_migrations"`
}

func GetConfig() *VervConfig {
	if vervConfig == nil {
		err := InitConfig(nil, nil)
		if err != nil {
			panic(err)
		}
	}

	return vervConfig
}

func InitConfig(cmd *cobra.Command, _ []string) error {
	if vervConfig == nil {
		vervConfig = &VervConfig{}
	}

	err := yaml.Unmarshal(builtInConfig, vervConfig)
	if err != nil {
		panic(
			rerrors.Wrap(err, "error parsing built in config file"))
	}

	*vervConfig = mergeConfigs(getConfigFromEnvironment(), *vervConfig)

	configFromFile, err := getConfigFromFile(cmd)
	if err != nil {
		return rerrors.Wrap(err, "error obtaining config from custom file")
	}

	if configFromFile.Valid {
		*vervConfig = mergeConfigs(configFromFile.V, *vervConfig)
	}

	return nil
}

func getConfigFromEnvironment() (r VervConfig) {
	r.Env.PathToMain = os.Getenv(envPathToMain)
	r.Env.PathToConfig = os.Getenv(envPathToConfig)
	r.Env.PathToMigrations = os.Getenv(envPathToMigrations)
	r.Env.PathToServerDefinition = os.Getenv(envPathToServerDefinition)

	r.DefaultProjectGitPath = os.Getenv(envDefaultProjectGitPath)

	if pathToClients := strings.Split(os.Getenv(envPathToClients), ","); pathToClients[0] != "" {
		r.Env.PathsToClients = pathToClients
	}

	if pathToServers := strings.Split(os.Getenv(envPathToServers), ","); pathToServers[0] != "" {
		r.Env.PathToServers = pathToServers
	}

	if pathToProtoClients := strings.Split(os.Getenv(envPathToProtoClients), ","); pathToProtoClients[0] != "" {
		r.Env.PathsToProtoClients = pathToProtoClients
	}

	if pathToCompiledClients := strings.Split(os.Getenv(envPathToCompiledClients), ","); pathToCompiledClients[0] != "" {
		r.Env.PathsToCompiledClients = pathToCompiledClients
	}

	return
}

func getConfigFromFile(cmd *cobra.Command) (sql.Null[VervConfig], error) {
	if cmd == nil {
		return sql.Null[VervConfig]{}, nil
	}

	cfgFilePath := cmd.Flag(CustomPathToConfig).Value.String()

	if cfgFilePath == "" {
		cfgFilePath = os.Getenv(environmentPathToConfig)
	}

	if cfgFilePath == "" {
		exePath, _ := os.Executable()
		cfgFilePath = path.Join(path.Dir(exePath), configFilename)
	}

	file, err := os.ReadFile(cfgFilePath)
	if err != nil && !rerrors.Is(err, os.ErrNotExist) {
		return sql.Null[VervConfig]{}, rerrors.Wrap(err, "error reading file from FS")
	}

	if len(file) == 0 {
		return sql.Null[VervConfig]{}, nil
	}

	var externalConf VervConfig
	err = yaml.Unmarshal(file, &externalConf)
	if err != nil {
		return sql.Null[VervConfig]{}, rerrors.Wrap(err, "error unmarshalling config from: "+cfgFilePath)
	}

	return sql.Null[VervConfig]{
		V:     externalConf,
		Valid: true,
	}, nil
}

func mergeConfigs(master, slave VervConfig) VervConfig {
	if master.Env.PathToMain == "" {
		master.Env.PathToMain = slave.Env.PathToMain
	}

	if master.Env.PathToConfig == "" {
		master.Env.PathToConfig = slave.Env.PathToConfig
	}

	if master.DefaultProjectGitPath == "" {
		master.DefaultProjectGitPath = slave.DefaultProjectGitPath
	}

	if len(master.Env.PathsToClients) == 0 {
		master.Env.PathsToClients = slave.Env.PathsToClients
	}

	if len(master.Env.PathToServers) == 0 {
		master.Env.PathToServers = slave.Env.PathToServers
	}

	if master.Env.PathToMigrations == "" {
		master.Env.PathToMigrations = slave.Env.PathToMigrations
	}

	if master.Env.PathToServerDefinition == "" {
		master.Env.PathToServerDefinition = slave.Env.PathToServerDefinition
	}

	if len(master.Env.PathsToProtoClients) == 0 {
		master.Env.PathsToProtoClients = slave.Env.PathsToProtoClients
	}

	if len(master.Env.PathsToCompiledClients) == 0 {
		master.Env.PathsToCompiledClients = slave.Env.PathsToCompiledClients
	}

	return master
}
