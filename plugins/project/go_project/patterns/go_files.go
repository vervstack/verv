package patterns

import (
	_ "embed"
)

// Constants naming: Purpose+Type (File)

const (
	GithubFolder    = ".github"
	WorkflowsFolder = "workflows"

	CmdFolder     = "cmd"
	ServiceFolder = "service"
	MainFileName  = "main.go"

	InternalFolder             = "internal"
	AppFolder                  = "app"
	AppFileName                = "app.go"
	AppInitServerFileName      = "server.go"
	AppInitDataSourcesFileName = "data_sources.go"
	AppConfigFileName          = "config.go"
	AppCustomFileName          = "custom.go"

	ConnFileName = "conn.go"

	TransportFolder = "transport"

	ServerManagerFileName          = "manager.go"
	GrpcServerFileName             = "grpc.go"
	HttpServerFileName             = "http.go"
	TelegramListenerFileName       = "listener.go"
	TelegramVersionHandlerFileName = "handler.go"

	HandlersFolderName = "handlers"
	VersionFolderName  = "version"

	ConfigsFolder      = "config"
	ConfigTemplateYaml = "config_template.yaml"

	ConfigLoadFileName        = "load.go"
	ConfigDataSourcesFileName = "data_sources.go"
	ConfigEnvironmentFileName = "environment.go"
	ConfigServersFileName     = "servers.go"
	ConfigEnvExampleFile      = ".env.example"

	GoMod = "go.mod"
)
