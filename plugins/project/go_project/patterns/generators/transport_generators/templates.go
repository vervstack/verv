package transport_generators

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed templates/manager.go.pattern
	serverManagerPattern  string
	serverManagerTemplate *template.Template

	//go:embed templates/grpc.go.pattern
	grpcServerFile []byte

	// httpServerPattern is embedded content, not a text/template: it contains a
	// runtime HTML template (`{{ range .Routes }}`) as a Go string literal, which
	// would be misparsed as template actions if run through text/template.
	//go:embed templates/http.go.pattern
	httpServerPattern string

	//go:embed templates/telegram/listener.go.pattern
	telegramListenerPattern  string
	telegramListenerTemplate *template.Template

	//go:embed templates/telegram/version/handler.go.pattern
	telegramVersionHandlerPattern  string
	telegramVersionHandlerTemplate *template.Template

	//go:embed templates/gateway_mux.go.pattern
	gatewayMuxPattern  string
	gatewayMuxTemplate *template.Template
)

//nolint:gochecknoinits // one-time compile of embedded templates into package-level *template.Template values
func init() {
	serverManagerTemplate = template.Must(
		template.New("server_manager").
			Parse(serverManagerPattern))

	telegramListenerTemplate = template.Must(
		template.New("telegram_listener").
			Parse(telegramListenerPattern))

	telegramVersionHandlerTemplate = template.Must(
		template.New("telegram_version_handler").
			Parse(telegramVersionHandlerPattern))

	gatewayMuxTemplate = template.Must(
		template.New("gateway_mux").
			Parse(gatewayMuxPattern))
}
