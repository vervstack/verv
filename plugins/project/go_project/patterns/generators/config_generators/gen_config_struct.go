package config_generators

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/internal/rw"
	"go.vervstack.ru/verv/plugins/project/config"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
)

type loadConfigFileGenArgs struct {
	Configs []InternalConfig
}

type InternalConfig struct {
	FieldName    string
	StructName   string
	From         string
	ErrorMessage string
}

type internalConfigGenerator func() (InternalConfig, *folder.Folder, error)

func GenerateConfigFolder(cfg *config.Config) (*folder.Folder, error) {
	args := loadConfigFileGenArgs{}

	configFolder := &folder.Folder{}

	generators := make([]internalConfigGenerator, 0, 3)

	if len(cfg.Servers) != 0 {
		generators = append(generators, newGenerateServerConfigStruct(cfg.Servers))
	}

	// Data sources
	if len(cfg.DataSources) != 0 {
		generators = append(generators, newGenerateDataSourcesConfigStruct(cfg.DataSources))
	}

	// Environment
	if len(cfg.Environment) != 0 {
		generators = append(generators, newGenerateEnvironmentConfigStruct(cfg.Environment))
	}

	for _, g := range generators {
		ic, f, err := g()
		if err != nil {
			return nil, rerrors.Wrap(err)
		}

		configFolder.Add(f)
		args.Configs = append(args.Configs, ic)
	}

	autoLoadFile := &rw.RW{}
	err := configAutoLoadTemplate.Execute(autoLoadFile, args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating load-config file ")
	}

	configFolder.Add(&folder.Folder{
		Name:    patterns.ConfigLoadFileName,
		Content: autoLoadFile.Bytes(),
	})
	return configFolder, nil
}
