package main_generators

import (
	_ "embed"
	"text/template"
)

var (
	//go:embed templates/main.go.pattern
	mainPattern  string
	mainTemplate *template.Template
)

//nolint:gochecknoinits // one-time compile of embedded template into a package-level *template.Template value
func init() {
	mainTemplate = template.Must(
		template.New("main").
			Parse(mainPattern))
}
