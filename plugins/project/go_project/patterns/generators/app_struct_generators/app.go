package app_struct_generators

import (
	"maps"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/pkg/matreshka"

	"go.vervstack.ru/verv/internal/rw"
	"go.vervstack.ru/verv/plugins/project"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators"
)

type AppFileGenArgs struct {
	Imports    map[string]string
	AppContent []AppContent
	Starters   []AppStarter
}

type AppContent struct {
	Comment              string
	Fields               []generators.KeyValue
	InitFunc             string
	InitFuncErrorMessage string
	Imports              map[string]string
}

type AppStarter struct {
	FieldName string
	StartCall string
	StopCall  string
}

func GenerateAppFiles(p project.IProject) (map[string][]byte, error) {
	initAppArgs := AppFileGenArgs{
		Imports: make(map[string]string),
	}

	cfg := p.GetConfig()

	out := make(map[string][]byte)

	err := initAppArgs.addDataSources(cfg.DataSources, out)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	err = initAppArgs.addServers(cfg.Servers, out)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	err = initAppArgs.mergeContentImports()
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	mainAppFile := &rw.RW{}

	err = appTemplate.Execute(mainAppFile, initAppArgs)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating app file")
	}

	out[patterns.AppFileName] = mainAppFile.Bytes()
	out[patterns.AppConfigFileName] = appConfigPattern

	if p.GetFolder().GetByPath(patterns.InternalFolder, patterns.AppFolder, patterns.AppCustomFileName) == nil {
		out[patterns.AppCustomFileName] = customPattern
	}

	return out, nil
}

func (a *AppFileGenArgs) addDataSources(dataSources matreshka.DataSources, out map[string][]byte) error {
	if len(dataSources) == 0 {
		return nil
	}

	initDataSourcesArgs, initDataSourcesFile, err := generateDataSourceInitFileAndArgs(dataSources)
	if err != nil {
		return rerrors.Wrap(err, "error generation data source init file")
	}

	// Every data source may have opted out of app-struct wiring (e.g. a
	// Postgres-only project), in which case there's no file or app content to
	// add at all.
	if initDataSourcesArgs == nil {
		return nil
	}

	out[patterns.AppInitDataSourcesFileName] = initDataSourcesFile

	a.AppContent = append(a.AppContent, *initDataSourcesArgs)

	return nil
}

func (a *AppFileGenArgs) addServers(servers matreshka.Servers, out map[string][]byte) error {
	if len(servers) == 0 {
		return nil
	}

	initServerArgs, initServerFile, err := generateServerInitFileAndArgs(servers)
	if err != nil {
		return rerrors.Wrap(err, "error generating server init file")
	}

	out[patterns.AppInitServerFileName] = initServerFile

	depArgs := InitDepFuncGenArgs{
		InitFunctionName: "InitServers",
		Imports:          initServerArgs.Imports,
	}

	for _, sn := range initServerArgs.Servers {
		depArgs.Functions = append(depArgs.Functions,
			InitFuncCall{
				ResultName: sn.ServerName,
				ResultType: "net.Listener",
			})
	}

	a.addAppContent(
		"/* Servers network listeners */",
		"error during network listeners initialization",
		depArgs,
	)

	return nil
}

func (a *AppFileGenArgs) mergeContentImports() error {
	for _, ac := range a.AppContent {
		for importPath, dependencyAlias := range ac.Imports {
			appAlias, ok := a.Imports[importPath]
			if ok && appAlias != dependencyAlias {
				return rerrors.New("Fatal error: app already imported package " +
					importPath + " with alias " + appAlias +
					". But dependency requires this package to be imported as " +
					dependencyAlias)
			}

			a.Imports[importPath] = dependencyAlias
		}
	}

	return nil
}

func (a *AppFileGenArgs) addAppContent(comment, errMsg string, args InitDepFuncGenArgs) {
	serverAppContent := AppContent{
		Comment:              comment,
		Fields:               make([]generators.KeyValue, 0, len(args.Functions)),
		InitFunc:             args.InitFunctionName,
		InitFuncErrorMessage: errMsg,
	}
	for _, serverInitFunc := range args.Functions {
		serverAppContent.Fields = append(serverAppContent.Fields,
			generators.KeyValue{
				Key:   serverInitFunc.ResultName,
				Value: serverInitFunc.ResultType,
			})
	}

	maps.Copy(a.Imports, args.Imports)

	a.AppContent = append(a.AppContent, serverAppContent)
}
