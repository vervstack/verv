package app_struct_generators

import (
	"go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"go.vervstack.ru/verv/internal/rw"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns"
	"go.vervstack.ru/verv/plugins/project/go_project/patterns/generators"
)

func generateDataSourceInitFileAndArgs(dataSources matreshka.DataSources) (*AppContent, []byte, error) {
	const initDataSourcesFunctionName = "InitDataSources"
	initDsFileArgs := InitDepFuncGenArgs{
		InitFunctionName: initDataSourcesFunctionName,
		Imports:          make(map[string]string),
		Functions:        make([]InitFuncCall, 0, len(dataSources)),
	}

	appContent := &AppContent{
		Comment:              "/* Data source connection */",
		Fields:               nil,
		InitFunc:             initDataSourcesFunctionName,
		InitFuncErrorMessage: "error during data sources initialization",
		Imports:              make(map[string]string),
	}

	for _, ds := range dataSources {
		var fc InitFuncCall

		switch ds.GetType() {
		case resources.PostgresResourceName, resources.SqliteResourceName:
			fc = sqlInitFunc(ds, appContent)
		case resources.RedisResourceName:
			fc = redisInitFunc(ds, appContent)
		case resources.TelegramResourceName:
			fc = telegramInitFunc(ds, appContent)
		default:
			return nil, nil, rerrors.New("unknown resource " + ds.GetType())
		}
		for importPath, alias := range fc.Import {
			initDsFileArgs.Imports[importPath] = alias
		}

		initDsFileArgs.Functions = append(initDsFileArgs.Functions, fc)

		appContent.Fields = append(appContent.Fields,
			generators.KeyValue{
				Key:   fc.ResultName,
				Value: fc.ResultType,
			})
	}

	initDsFileArgs.Imports[patterns.ImportNameErrorsPackage] = ""

	file := &rw.RW{}
	err := initAppStructFuncTemplate.Execute(file, initDsFileArgs)
	if err != nil {
		return nil, nil, rerrors.Wrap(err, "error generating server init file")
	}

	return appContent, file.Bytes(), nil
}

func sqlInitFunc(res resources.Resource, appContent *AppContent) (fc InitFuncCall) {
	fc.ResultName = generators.NormalizeResourceName(res.GetName())
	fc.FuncName = "sqldb.New"
	fc.ResultType = "*sql.DB"
	fc.Args = "a.Cfg.DataSources." + fc.ResultName
	fc.ErrorMessage = "error during sql connection initialization"
	fc.Import = map[string]string{
		"proj_name/internal/clients/sqldb": "",
	}

	appContent.Imports["database/sql"] = ""

	return fc
}

func redisInitFunc(res resources.Resource, appContent *AppContent) (fc InitFuncCall) {
	fc.ResultName = generators.NormalizeResourceName(res.GetName())
	fc.ResultType = "*redis.Client"

	fc.FuncName = "redis.New"
	fc.Args = "a.Cfg.DataSources." + fc.ResultName
	fc.ErrorMessage = "error during redis connection initialization"

	fc.Import = map[string]string{
		"proj_name/internal/clients/redis": "",
	}

	appContent.Imports["github.com/go-redis/redis"] = ""

	return
}

func telegramInitFunc(res resources.Resource, appContent *AppContent) (fc InitFuncCall) {
	fc.ResultName = generators.NormalizeResourceName(res.GetName())
	fc.ResultType = "*go_tg.Bot"

	fc.FuncName = "telegram.New"
	fc.Args = "a.Cfg.DataSources." + fc.ResultName
	fc.ErrorMessage = "error during telegram bot initialization"

	fc.Import = map[string]string{
		"proj_name/internal/clients/telegram": "",
	}

	appContent.Imports["github.com/Red-Sock/go_tg"] = ""

	return
}
