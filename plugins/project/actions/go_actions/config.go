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

	// AllowedOriginsEvonName and CookieSecureEvonName are only registered when the project
	// has at least one server (len(cfg.Servers) != 0) — they're meaningless without an
	// HTTP/gateway server. They're one *available* source for the allowedOrigins/cookieSecure
	// values a project author can wire into their hand-written custom.go call to
	// NewServerManager/transport.NewGatewayMux — nothing generated reads them automatically
	// (see transport_generators/templates/http.go.pattern's CORS change: allowedOrigins is a
	// plain function parameter, no baked-in default or auto-wired source).
	AllowedOriginsEvonName = "allowed_origins"
	CookieSecureEvonName   = "cookie_secure"
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

	if len(cfg.Servers) != 0 {
		if envVars[AllowedOriginsEvonName] == nil {
			allowedOrigins := environment.MustNewVariable(AllowedOriginsEvonName, "")

			allowedOrigins.Comment = "Comma-separated list of allowed CORS origins for the HTTP gateway. " +
				"No default — wire it into your custom.go call to NewServerManager."

			cfg.Environment = append(cfg.Environment, allowedOrigins)
		}

		if envVars[CookieSecureEvonName] == nil {
			cookieSecure := environment.MustNewVariable(CookieSecureEvonName, true)

			cookieSecure.Comment = "Secure attribute on auth cookies set by transport.NewGatewayMux. " +
				"Set false only for plain-HTTP local dev."

			cfg.Environment = append(cfg.Environment, cookieSecure)
		}
	}

	return nil
}
func (a GenerateProjectConfig) NameInAction() string {
	return "Generating project config"
}

type PrepareConfigFolder struct{}

func (a PrepareConfigFolder) Do(p project.IProject) (err error) {
	// generateConfigYamlFile must run before GenerateConfigFolder: it's what
	// (re)computes config.yaml's content for *this* run into the in-memory
	// folder tree (p.GetFolder()) — the real filesystem isn't touched until
	// the whole action pipeline flushes at the end. GenerateConfigFolder
	// embeds a byte-for-byte copy of that content as the config skeleton
	// (internal/config/skeleton.yaml), so it needs the up-to-date bytes, not
	// whatever was on disk from the previous generation run.
	err = a.generateConfigYamlFile(p)
	if err != nil {
		return rerrors.Wrap(err, "error generating config yaml-files")
	}

	hasDotEnv := p.GetFolder().GetByPath(patterns.EnvFile) != nil

	cfgFolder, err := config_generators.GenerateConfigFolder(p.GetConfig(), a.resolveConfigYamlBytes(p), hasDotEnv)
	if err != nil {
		return rerrors.Wrap(err, "error generating config folder")
	}

	cfgFolder.Name = path.Join(patterns.InternalFolder, patterns.ConfigsFolder)

	p.GetFolder().Add(cfgFolder)

	err = a.generateEnvExampleFile(p)
	if err != nil {
		return rerrors.Wrap(err, "error generating .env.example file")
	}

	return nil
}
func (a PrepareConfigFolder) NameInAction() string {
	return "Preparing config folder"
}

// resolveConfigYamlBytes returns the just-generated config.yaml content from
// the in-memory folder tree (populated by generateConfigYamlFile, which must
// run before this is called). Returns nil if it's genuinely absent — callers
// treat that as "no config content yet" and fall back to a stub.
func (a PrepareConfigFolder) resolveConfigYamlBytes(p project.IProject) []byte {
	configFolder := p.GetFolder().GetByPath(patterns.ConfigsFolder)
	if configFolder == nil {
		return nil
	}

	configFile := configFolder.GetByPath(patterns.ConfigMasterYamlFile)
	if configFile == nil {
		return nil
	}

	return configFile.Content
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
