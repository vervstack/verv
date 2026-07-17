package grpc_api_generator

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed templates/api.proto.pattern
	basicProto            string
	basicApiProtoTemplate *template.Template
)

//nolint:gochecknoinits // one-time compile of embedded template into a package-level *template.Template value
func init() {
	basicApiProtoTemplate = template.Must(
		template.New("basic_proto").
			Parse(basicProto))
}
