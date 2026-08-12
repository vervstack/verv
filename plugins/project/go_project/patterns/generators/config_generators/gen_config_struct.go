package config_generators

import (
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io/folder"
	"go.vervstack.ru/verv/internal/rw"
	"go.vervstack.ru/verv/plugins/project/config"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators/config_generators/env_config_generator"
)

type loadConfigFileGenArgs struct {
	Configs   []generators.InternalConfig
	HasDotEnv bool
}

type internalConfigGenerator func() (generators.InternalConfig, *folder.Folder, error)

func GenerateConfigFolder(cfg *config.Config, configYamlBytes []byte, hasDotEnv bool) (*folder.Folder, error) {
	args := loadConfigFileGenArgs{}
	args.HasDotEnv = hasDotEnv

	configFolder := &folder.Folder{}

	configGenerators := make([]internalConfigGenerator, 0)

	if len(cfg.Servers) != 0 {
		configGenerators = append(configGenerators, newGenerateServerConfigStruct(cfg.Servers))
	}

	// Data sources
	if len(cfg.DataSources) != 0 {
		configGenerators = append(configGenerators, newGenerateDataSourcesConfigStruct(cfg.DataSources))
	}

	// Environment
	if len(cfg.Environment) != 0 {
		configGenerators = append(configGenerators, env_config_generator.NewGenerateEnvironmentConfigStruct(cfg.Environment))
	}

	for _, g := range configGenerators {
		ic, f, err := g()
		if err != nil {
			return nil, rerrors.Wrap(err)
		}

		configFolder.Add(f)

		args.Configs = append(args.Configs, ic)
	}

	skeletonGoFile, skeletonYamlFile, err := newGenerateConfigSkeleton(configYamlBytes)()
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating config skeleton")
	}

	configFolder.Add(skeletonGoFile, skeletonYamlFile)

	autoLoadFile := &rw.RW{}

	err = configAutoLoadTemplate.Execute(autoLoadFile, args)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating load-config file ")
	}

	configFolder.Add(&folder.Folder{
		Name:    patterns.ConfigLoadFileName,
		Content: autoLoadFile.Bytes(),
	})

	return configFolder, nil
}
