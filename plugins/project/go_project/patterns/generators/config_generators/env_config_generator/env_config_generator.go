package env_config_generator

import (
	"reflect"

	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/environment"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/internal/rw"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators"
)

type structGenArgs struct {
	GenTag string

	StructName string
	Imports    map[string]string // path to alias
	Fields     []generators.KeyValue
	Enums      []enumGenArg
}

type enumGenArg struct {
	Name   string
	Values []generators.KeyValue
}

func newStructGenArgs(structName string) structGenArgs {
	return structGenArgs{
		GenTag:     genTag,
		StructName: structName,
		Imports:    map[string]string{},
	}
}

func NewGenerateEnvironmentConfigStruct(env matreshka.Environment,
) func() (generators.InternalConfig, *folder.Folder, error) {
	return func() (generators.InternalConfig, *folder.Folder, error) {
		ic := generators.InternalConfig{
			FieldName:    "Environment",
			StructName:   "EnvironmentConfig",
			From:         reflect.TypeOf(matreshka.Environment{}).Name(),
			ErrorMessage: "error parsing environment config",
		}

		genArgs := newStructGenArgs(ic.StructName)

		for _, e := range env {
			err := appendEnvField(e, &genArgs)
			if err != nil {
				return generators.InternalConfig{}, nil, err
			}
		}

		buf := &rw.RW{}
		err := structTemplate.Execute(buf, genArgs)
		if err != nil {
			return generators.InternalConfig{}, nil, rerrors.Wrap(err, "error executing config struct template")
		}

		f := &folder.Folder{
			Name:    patterns.ConfigEnvironmentFileName,
			Content: buf.Bytes(),
		}

		return ic, f, nil
	}
}

func appendEnvField(env *environment.Variable, genArgs *structGenArgs) error {
	var fieldKV generators.KeyValue
	fieldKV.Key = generators.NormalizeResourceName(env.Name)

	v := env.Value.Value()

	if v != nil {
		refVal := reflect.ValueOf(v)
		tp := refVal.Type()
		fieldKV.Value = tp.String()

		if tp.PkgPath() != "" {
			genArgs.Imports[tp.PkgPath()] = ""
		}
	} else {
		fieldKV.Value = ""
	}

	enumVal := env.Enum.Value()
	if !env.Enum.IsZero() {
		switch vals := enumVal.(type) {
		case []string:
			enumToGen := enumGenArg{
				Name: generators.NormalizeResourceName(env.Name),
			}

			enumToGen.Values = make([]generators.KeyValue, 0, len(vals))
			for _, old := range vals {
				enumToGen.Values = append(enumToGen.Values,
					generators.KeyValue{
						Key:   generators.NormalizeResourceName(old),
						Value: "\"" + old + "\"",
					})
			}
			genArgs.Enums = append(genArgs.Enums, enumToGen)
		case []int:
		default:
			return rerrors.New("error generating enums for config value. Unsupported enum type %T. Expected String slice",
				enumVal)
		}
	}

	genArgs.Fields = append(genArgs.Fields, fieldKV)

	return nil
}
