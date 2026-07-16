package config_generators

import (
	"reflect"

	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators"
)

type generalGenArgs struct {
	GenTag string

	StructName string
	Imports    map[string]string // path to alias
	Fields     []generators.KeyValue
	Enums      []EnumGenArg
}

type EnumGenArg struct {
	Name   string
	Values []generators.KeyValue
}

func newConfigStructGenArgs(structName string) generalGenArgs {
	return generalGenArgs{
		GenTag:     defaultGenTag,
		StructName: structName,
		Imports:    map[string]string{},
		Fields:     nil,
	}
}

func getTypeName(in any) string {
	return reflect.ValueOf(in).Type().Name()
}
