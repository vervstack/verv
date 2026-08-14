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

	TransportFolder  = "transport"
	MiddlewareFolder = "middleware"

	ServerManagerFileName          = "manager.go"
	GrpcServerFileName             = "grpc.go"
	HttpServerFileName             = "http.go"
	GatewayMuxFileName             = "gateway_mux.go"
	TelegramListenerFileName       = "listener.go"
	TelegramVersionHandlerFileName = "handler.go"

	CookieNamesFileName            = "cookie_names.go"
	CookieAnnotatorFileName        = "cookie_annotator.go"
	CookieResponseFileName         = "cookie_response.go"
	CSRFInterceptorFileName        = "csrf_interceptor.go"
	RequestSchemeAnnotatorFileName = "request_scheme_annotator.go"

	HandlersFolderName = "handlers"
	VersionFolderName  = "version"

	ConfigsFolder      = "config"
	ConfigTemplateYaml = "config_template.yaml"

	ConfigLoadFileName        = "load.go"
	ConfigDataSourcesFileName = "data_sources.go"
	ConfigEnvironmentFileName = "environment.go"
	ConfigServersFileName     = "servers.go"
	ConfigEnvExampleFile      = ".env.example"

	ConfigSkeletonGoFileName   = "skeleton.go"
	ConfigSkeletonYamlFileName = "skeleton.yaml"

	GoMod   = "go.mod"
	EnvFile = ".env"
)
