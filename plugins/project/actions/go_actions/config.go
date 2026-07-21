package go_actions

import (
	"bytes"
	"fmt"
	"path"
	"sort"
	"strings"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/environment"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/config_generators"
)

const (
	LogLevelEvonName = "log_level"
	LogLevelTrace    = "Trace"
	LogLevelDebug    = "Debug"
	LogLevelInfo     = "Info"
	LogLevelWarn     = "Warn"
	LogLevelError    = "Error"
	LogLevelFatal    = "Fatal"
	LogLevelPanic    = "Panic"

	LogFormatEvonName = "log_format"
	LogFormatJSON     = "JSON"
	LogFormatTEXT     = "TEXT"
)

type GenerateProjectConfig struct {
}

func (a GenerateProjectConfig) Do(p project.IProject) error {
	cfg := p.GetConfig()

	envVars := map[string]*environment.Variable{}

	for _, v := range cfg.Environment {
		envVars[v.Name] = v
	}

	if envVars[LogLevelEvonName] == nil {
		cfg.Environment = append(cfg.Environment,
			environment.MustNewVariable(LogLevelEvonName, LogLevelInfo,
				environment.WithEnum(
					LogLevelTrace,
					LogLevelDebug,
					LogLevelInfo,
					LogLevelWarn,
					LogLevelError,
					LogLevelFatal,
					LogLevelPanic,
				),
			),
		)
	}

	if envVars[LogFormatEvonName] == nil {
		cfg.Environment = append(cfg.Environment,
			environment.MustNewVariable(LogFormatEvonName, LogFormatTEXT,
				environment.WithEnum(
					LogFormatJSON,
					LogFormatTEXT,
				)))
	}

	return nil
}
func (a GenerateProjectConfig) NameInAction() string {
	return "Generating project config"
}

type PrepareConfigFolder struct{}

func (a PrepareConfigFolder) Do(p project.IProject) (err error) {
	cfgFolder, err := config_generators.GenerateConfigFolder(p.GetConfig())
	if err != nil {
		return rerrors.Wrap(err, "error generating config folder")
	}

	cfgFolder.Name = path.Join(patterns.InternalFolder, patterns.ConfigsFolder)

	p.GetFolder().Add(cfgFolder)

	err = a.generateConfigYamlFile(p)
	if err != nil {
		return rerrors.Wrap(err, "error generating config yaml-files")
	}

	err = a.generateEnvExampleFile(p)
	if err != nil {
		return rerrors.Wrap(err, "error generating .env.example file")
	}

	return nil
}
func (a PrepareConfigFolder) NameInAction() string {
	return "Preparing config folder"
}

func (a PrepareConfigFolder) generateConfigYamlFile(p project.IProject) (err error) {
	configFolder := p.GetFolder().GetByPath(patterns.ConfigsFolder)

	newConfig := p.GetConfig()

	sortEnv(newConfig.AppConfig)

	for _, cfgName := range []string{
		patterns.ConfigTemplateYaml,
		patterns.ConfigMasterYamlFile,
	} {
		newConfig.ServiceDiscovery = matreshka.ServiceDiscovery{}
		err := appendToConfig(newConfig.AppConfig, configFolder, cfgName)
		if err != nil {
			return rerrors.Wrap(err, "error appending changes to dev config")
		}
	}

	return nil
}

func appendToConfig(newConfig matreshka.AppConfig, configFolder *folder.Folder, path string) (err error) {
	currentConfig := matreshka.NewEmptyConfig()

	configFile := configFolder.GetByPath(path)
	if configFile == nil {
		configFile = &folder.Folder{
			Name: path,
		}
		configFolder.Add(configFile)
	}

	if len(configFile.Content) != 0 {
		err = currentConfig.Unmarshal(configFile.Content)
		if err != nil {
			return rerrors.Wrap(err, "error reading dev config file")
		}
	}

	currentConfig = matreshka.MergeConfigs(currentConfig, newConfig)
	sortEnv(currentConfig)
	configFile.Content, err = currentConfig.Marshal()
	if err != nil {
		return rerrors.Wrap(err, "error marshalling dev config to yaml")
	}

	return nil
}

func sortEnv(cfg matreshka.AppConfig) {
	sort.Slice(cfg.Environment, func(i, j int) bool {
		return cfg.Environment[i].Name < cfg.Environment[j].Name
	})
}

func (a PrepareConfigFolder) generateEnvExampleFile(p project.IProject) error {
	cfg := p.GetConfig()

	var allNodes []*evon.Node

	if len(cfg.Environment) > 0 {
		nodes, err := cfg.Environment.MarshalEnv("ENVIRONMENT")
		if err != nil {
			return rerrors.Wrap(err, "error marshalling environment to env")
		}
		allNodes = append(allNodes, nodes...)
	}

	if len(cfg.DataSources) > 0 {
		nodes, err := cfg.DataSources.MarshalEnv("DATA_SOURCES")
		if err != nil {
			return rerrors.Wrap(err, "error marshalling data sources to env")
		}
		allNodes = append(allNodes, nodes...)
	}

	if len(cfg.Servers) > 0 {
		// Servers.MarshalEnv mutates each *server.Server.Name as a side
		// effect (assigns it a canonical, upper-cased name). Marshal a copy
		// so later generators (app struct, config struct) still see the
		// original name and derive a consistent Go identifier from it.
		serversCopy := make(matreshka.Servers, len(cfg.Servers))
		for port, srv := range cfg.Servers {
			srvCopy := *srv
			serversCopy[port] = &srvCopy
		}

		nodes, err := serversCopy.MarshalEnv("SERVERS")
		if err != nil {
			return rerrors.Wrap(err, "error marshalling servers to env")
		}
		allNodes = append(allNodes, nodes...)
	}

	configFolder := p.GetFolder().GetByPath(patterns.ConfigsFolder)
	configFolder.Add(&folder.Folder{
		Name:    patterns.ConfigEnvExampleFile,
		Content: marshalEnvExample(allNodes),
	})

	return nil
}

// marshalEnvExample serializes evon nodes to .env format.
// Nodes with a Value are printed directly (env vars and plain struct fields);
// nodes without a Value are structural containers whose InnerNodes are visited.
// Dashes are replaced with underscores so names match the form apps read.
// Nodes whose names contain characters outside [A-Z0-9_] (e.g. route path
// entries like "/{GRPC}") are skipped — they cannot be valid env var names.
func marshalEnvExample(nodes []*evon.Node) []byte {
	b := &bytes.Buffer{}
	for _, node := range nodes {
		if node.Value != nil {
			name := strings.ReplaceAll(node.Name, "-", "_")
			if isEnvVarName(name) {
				fmt.Fprintf(b, "%s=%v\n", name, node.Value)
			}
		} else {
			b.Write(marshalEnvExample(node.InnerNodes))
		}
	}

	return b.Bytes()
}

func isEnvVarName(s string) bool {
	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}

	return true
}
