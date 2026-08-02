# Graph Report - verv  (2026-08-02)

## Corpus Check
- 135 files · ~31,215 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1685 nodes · 2652 edges · 224 communities (112 shown, 112 thin omitted)
- Extraction: 85% EXTRACTED · 15% INFERRED · 0% AMBIGUOUS · INFERRED: 400 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `03aaf9a3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Project Action Pipeline
- gRPC Package Discovery
- Dockerfile Generation & Build Tests
- App File Generator
- Environment File Fetching
- Compiled Pattern Scaffold Files
- Config Folder Generation
- Env Tidy & Terminal Loader UI
- Project Structure Preparation Actions
- gRPC Client/Gateway Interfaces
- IO Mock (minimock)
- Docker Compose Assembly
- CLI Command Wiring (environment)
- IProject Mock (minimock)
- Postgres Client & Config
- Generated Proto Validation (Ping)
- gRPC Server Transport
- TS gRPC-Gateway Fetch Client
- Project Init Command
- Terminal Color Parser & IO Stub
- Generated Protobuf Message Methods
- Dependency AppendToProject (Sqlite/Env)
- IProjectMock GetFolder Expectations
- Telegram Dependency Wiring
- Sqlite Dependency Client
- IOMock GetInput Expectations
- gRPC Server Dependency Wiring
- Environment Init & File IO
- gRPC Server Transport (pattern)
- Virtual Folder Tree
- Go Fmt & Makefile Gen Actions
- HTTP Server Transport
- RsCli Config Loading
- HTTP Server Transport (pattern)
- Project Name Collection Prompt
- Postgres Dependency Client
- IProjectMock GetType Expectations
- Multiplexed Server Manager
- Telegram Version Handler
- Telegram Transport Listener
- Environment Config Struct
- Folder Loader Options
- Git Status Diff
- Build Project Action
- ProjEnv Config Access
- Migration Tool Interface
- RW File Locking Utility
- IOMock Error Expectations
- IProjectMock GetConfig Expectations
- Env Config/Variables Fetching
- Generated Project Config Loader
- Redis Dependency Client
- SQL Connection (pattern)
- Git Init Action
- RW Read/Execute
- IOMock Print Expectations
- IOMock Println Expectations
- Generator Tests & Folder Comparison
- SQL Connection (pattern, dup)
- Environment Aggregate Structs
- Port Manager
- IOMock PrintColored Expectations
- Env Variables Manager
- Env Install Command (unwired)
- SQL Connection Pattern Wiring
- Git Commit Action
- Project Loader
- gRPC Dependency Test Fixtures
- gRPC Implementation Registration
- Server Start Lifecycle
- RW Struct Core Methods
- File Server Generator
- Proto API Generator
- Make Binary Installer
- IOMock GetInput Builder
- IProjectMock GetName Builder
- IProjectMock GetProjectPath Builder
- IProjectMock GetShortName Builder
- GlobalEnvironment Init
- GitHub Workflow Tidy Action
- Generator Arg Structs
- CI / CLAUDE.md Dev Commands
- Environment Plugin Overview
- Go Mod Init Action
- Main Entrypoint (external app)
- gRPC Client Connect Helper
- Server Stop Lifecycle
- Action Interface (go_actions)
- Progress Loader Interface
- YAML Deep Copy Utility
- Project Interface Duplication
- Telegram Bot Connection
- Pattern Compile CI Check
- Name Collector Struct
- OS Path Collection
- Generated Environment Config
- Generated Environment Config (dup)
- Dependency Project Interface
- GitHub Version Struct
- Project/RsCliConfig Pair
- Migration Tool Registry
- ProjEnv Tidy Server APIs
- ProjEnv Tidy Migration Dirs
- ProjEnv Tidy Service
- Telegram Folder Name
- Embedded Template Registration
- Embedded Template Registration
- Embedded Template Registration
- Pattern GitHub Workflows
- Generated README Reference
- Goose Migration Tool
- Architecture Command Flow Note
- Env Tidy Struct
- Env Install Struct
- macOS Env Install
- Env Install Run
- gRPC Server Option
- SQL DB Interface
- Goose Logger Adapter
- Unexpected Ping Error
- Init Proc Struct
- Project Init Struct
- Project Link Struct
- Project Tidy Struct
- Release Tag Job
- gRPC Request Struct
- Compose AppendService
- Compose Struct
- Compose Marshal
- Pattern Manager
- Percent Loader Progress
- Compose Rule Struct
- Port Manager GetPort
- Processor Struct
- RW Struct
- rscli Module Path
- Project Name Placeholder
- GlobalEnvironment Existence Check
- Env Resources Type
- gRPC Client Struct
- gRPC Client Folder Name
- gRPC Discovery Struct
- gRPC Server Folder Name
- SQL Connection Struct
- SQL Connection Folder Name
- GitHub Workflow Tidy Struct
- GitHub Workflow Action Name
- HTTP API Template
- GitIgnore Template
- Linter Config Template
- Makefile Template
- README Template
- RsCli Makefile Template
- Errors Package Import Name
- App File Gen Args
- App Starter Struct
- Init Func Call Struct
- Server Listener Args
- Server Listeners Args
- Internal Config Struct
- Internal Config Generator
- Load Config File Gen Args
- App Data Name Constant
- Data Volume Name Constant
- Service Dockerfile Args
- Service Proto API Args
- File Server Gen Args
- Generic Gen Args Struct
- Postgres Conn Template
- Redis Conn Template
- Sqlite Conn Template
- Keys File Template
- Main File Template
- gRPC Server Manager Template
- HTTP Server Manager Template
- Server Manager Template
- Telegram Handler Template
- Telegram Server Template
- EasyP Config Template
- gRPC Impl Template
- gRPC Impl Folder Constant
- gRPC Server Makefile Template
- Proto Contract Template
- Golangci Lint Config
- Environment Config Struct (dup)
- Telegram Server Start
- Telegram Server Stop
- Version Command Handler
- Project Name Substitution
- Mock Project WriteFile
- proj_name
- Proc
- pre-commit
- pre-commit
- go.vervstack.ru/verv
- AppConfig

## God Nodes (most connected - your core abstractions)
1. `IOMock` - 74 edges
2. `IProjectMock` - 65 edges
3. `IProject` - 43 edges
4. `Folder` - 35 edges
5. `RsCliConfig` - 33 edges
6. `IO` - 29 edges
7. `Execute()` - 25 edges
8. `Test_AddDependency()` - 24 edges
9. `GetMockProject()` - 22 edges
10. `New()` - 21 edges

## Surprising Connections (you probably didn't know these)
- `rscli Architecture Overview (Cobra CLI, two top-level commands)` --references--> `NewCmd()`  [INFERRED]
  CLAUDE.md → cmd/project/cmd.go
- `InitConfig()` --semantically_similar_to--> `Load function`  [INFERRED] [semantically similar]
  internal/config/cfg.go → cmd/project/init_new/expected/load_config/load_config_file.go
- `main()` --references--> `go.mod module manifest (github.com/Red-Sock/rscli)`  [INFERRED]
  main.go → go.mod
- `Config` --shares_data_with--> `expected grpc config_template.yaml fixture`  [INFERRED]
  plugins/project/config/config.go → cmd/project/add/expected/grpc/config/config_template.yaml
- `Config` --shares_data_with--> `expected grpc dev.yaml fixture`  [INFERRED]
  plugins/project/config/config.go → cmd/project/add/expected/grpc/config/dev.yaml

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Cobra CLI command tree wiring (root -> project/env -> subcommands)** — main_main, cmd_project_cmd_newcmd, cmd_environment_cmd_newcmd, cmd_project_add_processor_newcommand, cmd_project_link_newlinkcmd, cmd_project_tidy_newtidycmd [INFERRED 0.85]
- **grpc dependency test fixture set (config + generated transport)** — config_config, cmd_project_add_expected_grpc_config_config_template, cmd_project_add_expected_grpc_config_dev, cmd_project_add_expected_grpc_internal_transport_grpc_grpcserver, cmd_project_add_processor_test_test_adddependency [INFERRED 0.80]
- **Test case matrix for project add dependencies** — cmd_project_add_processor_test_expectedgrpc, cmd_project_add_processor_test_expectedredis, cmd_project_add_processor_test_expectedpostgres, cmd_project_add_processor_test_expectedtelegram, cmd_project_add_processor_test_expectedsqlite, cmd_project_add_processor_test_expectedenv [INFERRED 0.85]
- **Dependency client 'New' constructor pattern (redis/telegram/sqldb)** — cmd_project_add_expected_redis_internal_clients_redis_conn_new, cmd_project_add_expected_telegram_internal_clients_telegram_conn_new, cmd_project_add_expected_postgres_internal_clients_sqldb_conn_new [INFERRED 0.75]
- **rscli project init command flow** — cmd_project_init_new_processor_newcommand, cmd_project_init_new_processor_run, cmd_project_init_new_collect_name_collect, cmd_project_init_new_collect_os_path_collectospath, cmd_project_init_new_project_assembler_createproject [INFERRED 0.85]
- **Shared app_info/data_sources/environment config.yaml schema across dependencies** — cmd_project_add_expected_postgres_config_config_postgres, cmd_project_add_expected_redis_config_config_redis, cmd_project_add_expected_sqlite_config_config_sqlite, cmd_project_add_expected_telegram_config_config_telegram [INFERRED 0.75]
- **RsCliConfig layered load: built-in -> env -> file** — internal_config_cfg_initconfig, internal_config_cfg_getconfigfromenvironment, internal_config_cfg_getconfigfromfile, internal_config_cfg_mergeconfigs, internal_config_cfg_builtinconfig [INFERRED 0.85]
- **Progress interface and its renderers** — internal_io_loader_multi_loader_progress, internal_io_loader_infinite_progress_infiniteloader, internal_io_loader_percent_progress_percentloader, internal_io_loader_multi_loader_runmultiloader, internal_io_loader_seq_loader_runseqloader [INFERRED 0.85]
- **go:embed-backed folder.Folder template constants** — internal_envpatterns_embedded_envfile, internal_envpatterns_embedded_dockercomposefile, internal_envpatterns_embedded_makefile, internal_envpatterns_embedded_buildincomposeexamples, internal_io_folder_folder_folder [INFERRED 0.85]
- **ProjEnv.Tidy orchestration pipeline (resources, server APIs, flush, migrations)** — plugins_environment_project_project_tidy, plugins_environment_project_tidy_resources_tidyresources, plugins_environment_project_tidy_api_tidyserverapis, plugins_environment_project_project_flush, plugins_environment_project_tidy_migrations_tidymigrationdirs [EXTRACTED 1.00]
- **GlobalEnvironment.fetchFiles data-gathering pipeline** — plugins_environment_env_fetchfiles, plugins_environment_env_fetchsrcprojectdirs, plugins_environment_env_fetchcompose, plugins_environment_env_fetchdotenv, plugins_environment_env_fetchmakefile [EXTRACTED 1.00]
- **GlobalEnvironment.Init scaffolding pipeline** — plugins_environment_init_init, plugins_environment_init_initbasis, plugins_environment_init_initprojectsdirs, plugins_environment_init_initprojectdir [EXTRACTED 1.00]
- **Dependency interface implementers (Sqlite, EnvVariable, GrpcServer, Telegram)** — plugins_project_actions_go_actions_dependencies_dependencies_dependency, plugins_project_actions_go_actions_dependencies_client_sqlite_sqlite, plugins_project_actions_go_actions_dependencies_env_variable_envvariable, plugins_project_actions_go_actions_dependencies_server_grpc_grpcserver, plugins_project_actions_go_actions_dependencies_telegram_telegram [EXTRACTED 1.00]
- **Duplicated gRPC package-discovery logic between link_service/client_grpc.go and grpc_discovery package** — plugins_project_actions_go_actions_dependencies_link_service_client_grpc_grpcclient, plugins_project_actions_go_actions_dependencies_link_service_client_grpc_readgrpcpackage, plugins_project_actions_go_actions_dependencies_link_service_client_grpc_grpcclient_filterpackagename, plugins_project_actions_go_actions_dependencies_link_service_grpc_discovery_grpc_grpcdiscovery, plugins_project_actions_go_actions_dependencies_link_service_grpc_discovery_grpc_readgrpcpackagefrompackageclientpath, plugins_project_actions_go_actions_dependencies_link_service_grpc_discovery_name_filter_filterpackagename [INFERRED 0.85]
- **Flow for adding a sqlite dependency to a project** — plugins_project_actions_go_actions_dependencies_client_sqlite_sqlite_appendtoproject, plugins_project_actions_go_actions_dependencies_sql_sqlconn_applysqlconnectionfile, plugins_project_actions_go_actions_dependencies_sql_sqlconn_applysqldriver, plugins_project_go_project_patterns_go_clients_sqlconnfile [EXTRACTED 1.00]
- **Config Folder Generation Pipeline** — plugins_project_go_project_patterns_generators_config_generators_gen_config_struct_generateconfigfolder, plugins_project_go_project_patterns_generators_config_generators_gen_servers_newgenerateserverconfigstruct, plugins_project_go_project_patterns_generators_config_generators_gen_data_sources_newgeneratedatasourcesconfigstruct, plugins_project_go_project_patterns_generators_config_generators_gen_env_newgenerateenvironmentconfigstruct [EXTRACTED 1.00]
- **App Files Generation Pipeline** — plugins_project_go_project_patterns_generators_app_struct_generators_app_generateappfiles, plugins_project_go_project_patterns_generators_app_struct_generators_data_sources_generatedatasourceinitfileandargs, plugins_project_go_project_patterns_generators_app_struct_generators_server_generateserverinitfileandargs [EXTRACTED 1.00]
- **Data Source Resource Type Dispatch** — plugins_project_go_project_patterns_generators_app_struct_generators_data_sources_sqlinitfunc, plugins_project_go_project_patterns_generators_app_struct_generators_data_sources_redisinitfunc, plugins_project_go_project_patterns_generators_app_struct_generators_data_sources_telegraminitfunc, plugins_project_go_project_patterns_generators_app_struct_generators_data_sources_grpcinitfunc [EXTRACTED 1.00]
- **GrpcImpl / GrpcWithGateway Implementation Pattern** — cmd_project_add_expected_grpc_internal_transport_grpc_grpcimpl, cmd_project_add_expected_grpc_internal_transport_grpc_grpcwithgateway, cmd_project_add_expected_grpc_internal_transport_grpc_addimplementation, plugins_project_go_project_patterns_pattern_internal_transport_grpc_example_api_impl_impl_impl [INFERRED 0.85]
- **Start/Stop Server Lifecycle Pattern** — cmd_project_add_expected_grpc_internal_transport_manager_serversmanager, cmd_project_add_expected_grpc_internal_transport_grpc_grpcserver, cmd_project_add_expected_grpc_internal_transport_http_httpserver, plugins_project_go_project_patterns_pattern_internal_transport_telegram_listener_server [INFERRED 0.80]
- **Pluggable SQL Dialect Driver Registration** — cmd_project_add_expected_postgres_internal_clients_sqldb_conn_new, plugins_project_go_project_patterns_pattern_internal_clients_sqldb_postgres_pq_driver, plugins_project_go_project_patterns_pattern_internal_clients_sqldb_sqlite_sqlite_driver, plugins_project_go_project_patterns_pattern_internal_config_keys_resourcepostgres [INFERRED 0.75]
- **Generated gRPC Version RPC Stack** — plugins_project_go_project_patterns_pattern_pkg_example_api_api_pb_pingrequest, plugins_project_go_project_patterns_pattern_pkg_example_api_api_pb_pingresponse, plugins_project_go_project_patterns_pattern_pkg_example_api_api_grpc_pb_projnameapiclient, plugins_project_go_project_patterns_pattern_pkg_example_api_api_grpc_pb_projnameapiserver, plugins_project_go_project_patterns_pattern_pkg_example_api_api_pb_gw_registerprojnameapihandlerserver [INFERRED 0.85]
- **PingRequest Message Across Codegen Targets** — plugins_project_go_project_patterns_pattern_c_easyp, plugins_project_go_project_patterns_pattern_pkg_example_api_api_pb_pingrequest, plugins_project_go_project_patterns_pattern_pkg_web_grpc_api_pb_pingrequest, plugins_project_go_project_patterns_pattern_pkg_docs_grpc_api_swagger_pingrequest [INFERRED 0.85]
- **Goose Migration Tool Component** — plugins_tools_migrations_goose_tool_tool, plugins_tools_migrations_migrators_migrationtool, plugins_tools_shared_ghversion_version_githubversion [INFERRED 0.75]
- **Test support chain for Dockerfile ENV var feature** — tasks_task_002_dockerfile_env_vars, tests_project_mock_config_getallenvvariables, tests_project_mock_opts_withenvironmentvariables [INFERRED 0.75]

## Communities (224 total, 112 thin omitted)

### Community 0 - "Project Action Pipeline"
Cohesion: 0.07
Nodes (35): testCase, Dependency Names (grpc, redis, postgres, telegram, sqlite, env), endMsg constant, preparingMsg constant, startingMsg constant, Proc.run() method, expectedEnv(), expectedGrpc() (+27 more)

### Community 1 - "gRPC Package Discovery"
Cohesion: 0.06
Nodes (36): File, FuncDecl, GenDecl, GrpcDiscovery, GrpcPackage, SnakeToPascal(), ToPascal(), slices.Contains (+28 more)

### Community 2 - "Dockerfile Generation & Build Tests"
Cohesion: 0.06
Nodes (45): DirEntry, serviceDockerfileArgs, BuildProjectAction, BuildProjectSuite, BuildProjectSuite.Test_BuildProject, T, Test_BuildProject(), fullConfigGoFile (+37 more)

### Community 3 - "App File Generator"
Cohesion: 0.16
Nodes (14): InternalConfig, internalConfigGenerator, loadConfigFileGenArgs, Environment, getTypeName(), GenerateConfigFolder(), DataSources, newGenerateDataSourcesConfigStruct() (+6 more)

### Community 4 - "Environment File Fetching"
Cohesion: 0.05
Nodes (37): MewEmptyMakefile, NewMakeFile, parseRule, parseVariable, ReadMakeFile, NewPortManager, PortManager.SaveIfNotExist, Makefile (+29 more)

### Community 5 - "Compiled Pattern Scaffold Files"
Cohesion: 0.06
Nodes (44): easyp.yaml proto codegen config (pattern_c), branch-push GitHub Actions workflow, master-actions RELEASE GitHub Actions workflow, .golangci.yaml lint config (pattern_c), README.md template (pattern_c), RedSock CLI (referenced tool), proj_name_apiPingRequest swagger definition, proj_name_apiPingResponse swagger definition (+36 more)

### Community 6 - "Config Folder Generation"
Cohesion: 0.19
Nodes (9): ActionPerformerMock, ActionPerformerMockTidyExpectation, ActionPerformerMockTidyParams, ActionPerformerMockTidyResults, mActionPerformerMockTidy, IProject, RWMutex, Tester (+1 more)

### Community 7 - "Env Tidy & Terminal Loader UI"
Cohesion: 0.06
Nodes (23): go.mod module manifest (github.com/Red-Sock/rscli), Color type, TerminalColor(), InfiniteLoader.Done method, NewInfiniteLoader constructor, RunMultiLoader function, percentLoader.GetLoaderSymb method, NewPercentLoader constructor (+15 more)

### Community 8 - "Project Structure Preparation Actions"
Cohesion: 0.12
Nodes (10): PrepareMakefile, PrepareProjectStructure, PrepareServer, renamer.ReplaceProjectName, ReplaceProjectNameFull, ReplaceProjectNameShort(), Dependency.AppendToProject, addMissingImplFolders() (+2 more)

### Community 9 - "gRPC Client/Gateway Interfaces"
Cohesion: 0.09
Nodes (20): Request, RW, Impl, ProjNameAPIClient, ProjNameAPIServer, UnimplementedProjNameAPIServer, UnsafeProjNameAPIServer, Duration (+12 more)

### Community 10 - "IO Mock (minimock)"
Cohesion: 0.07
Nodes (4): IOMock, Duration, Tester, NewIOMock()

### Community 11 - "Docker Compose Assembly"
Cohesion: 0.05
Nodes (33): Compose, ContainerSettings, Pattern, PatternManager, compose.examples.yaml template, docker-compose.yaml template, ContainerSettings struct, AddEnvironmentBrackets function (+25 more)

### Community 12 - "CLI Command Wiring (environment)"
Cohesion: 0.22
Nodes (7): rscli Architecture Overview (Cobra CLI, two top-level commands), NewCmd() (env command constructor), Command, NewCmd(), newLinkCmd(), projectLink.run() method, projectLink

### Community 13 - "IProject Mock (minimock)"
Cohesion: 0.09
Nodes (4): IProjectMock, Duration, Tester, NewIProjectMock()

### Community 14 - "Postgres Client & Config"
Cohesion: 0.05
Nodes (31): httpServer struct, ServersManager struct, Postgres master config.yaml, Postgres dev.yaml config, DB, sqlLogger, sqldb.New (SQL connection + goose migration) function, postgres.go blank import of lib/pq driver (+23 more)

### Community 15 - "Generated Proto Validation (Ping)"
Cohesion: 0.08
Nodes (6): PingRequestMultiError, PingRequestValidationError, PingResponseMultiError, PingResponseValidationError, PingRequest, PingResponse

### Community 16 - "gRPC Server Transport"
Cohesion: 0.05
Nodes (24): GrpcImpl, grpcServer, GrpcWithGateway, grpcServer struct, newGrpcServer(), httpServer, newHttpServer function, setUpCors function (+16 more)

### Community 17 - "TS gRPC-Gateway Fetch Client"
Cohesion: 0.10
Nodes (20): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+12 more)

### Community 18 - "Project Init Command"
Cohesion: 0.09
Nodes (20): projectInit.buildProject() method, Proc, newInitCmd(), projectInit.obtainFolderPathFromUser() method, projectInit.obtainNameFromUser() method, Proc, Project, projectInit.run() method (+12 more)

### Community 19 - "Terminal Color Parser & IO Stub"
Cohesion: 0.17
Nodes (5): Color, UIColor function, IOMockPrintColoredExpectation, IOMockPrintColoredParams, mIOMockPrintColored

### Community 20 - "Generated Protobuf Message Methods"
Cohesion: 0.13
Nodes (5): file_grpc_api_proto_init(), file_grpc_api_proto_rawDescGZIP(), PingRequest, PingResponse, init()

### Community 21 - "Dependency AppendToProject (Sqlite/Env)"
Cohesion: 0.14
Nodes (19): Command, Dependency, dependencyBase, EnvVariable, redisClient(), sqlite(), sqlite() constructor, Dependency interface (+11 more)

### Community 22 - "IProjectMock GetFolder Expectations"
Cohesion: 0.14
Nodes (10): IProjectMockGetFolderExpectation, IProjectMockGetFolderResults, IProjectMockGetNameExpectation, IProjectMockGetNameResults, IProjectMockGetProjectPathExpectation, IProjectMockGetProjectPathResults, IProjectMockGetShortNameExpectation, IProjectMockGetShortNameResults (+2 more)

### Community 23 - "Telegram Dependency Wiring"
Cohesion: 0.42
Nodes (4): Telegram, containsDependencyFolder(), Project, New()

### Community 24 - "Sqlite Dependency Client"
Cohesion: 0.17
Nodes (12): Sqlite.AppendToProject, EnvVariable.AppendToProject, Project interface (dependencies pkg), sqlConn.applySqlConnectionFile, sqlConn.applySqlDriver, Telegram.AppendToProject, Telegram.applyClient, Telegram.applyConfig (+4 more)

### Community 25 - "IOMock GetInput Expectations"
Cohesion: 0.24
Nodes (4): IOMockGetInputOneOfExpectation, IOMockGetInputOneOfParams, IOMockGetInputOneOfResults, mIOMockGetInputOneOf

### Community 26 - "gRPC Server Dependency Wiring"
Cohesion: 0.27
Nodes (6): GrpcServer, GrpcServer.addGrpcServerToConfig, GrpcServer.AppendToProject, GrpcServer.applyApiFolder, initServerManagerFiles func, prepareServerConfig func

### Community 27 - "Environment Init & File IO"
Cohesion: 0.18
Nodes (22): execute(), executeWith(), GenerateGRPCConn(), GeneratePostgresConn(), GeneratePostgresDriver(), GeneratePostgresInstances(), GenerateRedisConn(), GenerateSQLConn() (+14 more)

### Community 28 - "gRPC Server Transport (pattern)"
Cohesion: 0.17
Nodes (9): AppFileGenArgs, AppStarter, AppFileGenArgs.addAppContent method, T, Test_GenerateAppFiles_CustomFileAlreadyExists(), GenerateAppFiles(), DataSources, Servers (+1 more)

### Community 29 - "Virtual Folder Tree"
Cohesion: 0.14
Nodes (3): Folder, OverrideFile(), mIProjectMockGetFolder

### Community 30 - "Go Fmt & Makefile Gen Actions"
Cohesion: 0.13
Nodes (8): GoFmt, InitGoMod, InitGoProjectApp, PrepareDockerfile, RunGoTidyAction, UpdateAllPackages, Action interface (actions pkg), Action interface (go_actions pkg)

### Community 31 - "HTTP Server Transport"
Cohesion: 0.23
Nodes (14): generateTransportFiles(), execute(), GenerateGrpcServer(), GenerateHttpServer(), GenerateServerManager(), GenerateTelegramListener(), GenerateTelegramVersionHandler(), Template (+6 more)

### Community 32 - "RsCli Config Loading"
Cohesion: 0.22
Nodes (16): Project, RsCliConfig, VervConfig, RunMakeGenAction, builtInConfig embedded var, GetConfig(), getConfigFromEnvironment(), getConfigFromFile() (+8 more)

### Community 33 - "HTTP Server Transport (pattern)"
Cohesion: 0.19
Nodes (7): newNameCollector(), Command, Proc, NewCommand(), nameCollector, ErrInvalidNameErr sentinel error, ValidateProjectNameStr()

### Community 34 - "Project Name Collection Prompt"
Cohesion: 0.07
Nodes (26): nameCollector.askUserForName method, nameCollector.collect method, newNameCollector function, nameCollector.preAppendHost method, nameCollector.removeHttpProtoc method, Test_collectName test function, Proc.collectOsPath method, Proc (+18 more)

### Community 35 - "Postgres Dependency Client"
Cohesion: 0.24
Nodes (7): PostgresInstance, PostgresInstancesArgs, Postgres, DataSources, Project, postgresInstances(), ReplaceProjectName()

### Community 36 - "IProjectMock GetType Expectations"
Cohesion: 0.15
Nodes (5): IProjectMockGetTypeExpectation, IProjectMockGetTypeResults, mIProjectMockGetType, Project, Type

### Community 37 - "Multiplexed Server Manager"
Cohesion: 0.26
Nodes (9): PrepareConfigFolder, Node, appendToConfig(), AppConfig, isEnvVarName(), marshalEnvExample(), PrepareConfigFolder.generateConfigYamlFile, PrepareConfigFolder.generateEnvExampleFile (+1 more)

### Community 38 - "Telegram Version Handler"
Cohesion: 0.24
Nodes (9): enumGenArg, structGenArgs, InternalConfig, appendEnvField(), NewGenerateEnvironmentConfigStruct(), newStructGenArgs(), T, Test_GenerateEnvConfig() (+1 more)

### Community 39 - "Telegram Transport Listener"
Cohesion: 0.35
Nodes (10): buildCLI(), T, run(), Test_GeneratedProjectsCompile(), T, scaffoldComboProject(), stripGeneratedMarker(), Test_GeneratedProject_GoimportsClean() (+2 more)

### Community 40 - "Environment Config Struct"
Cohesion: 0.44
Nodes (10): AppContent, generateDataSourceInitFileAndArgs(), DataSources, Resource, grpcInitFunc function, redisInitFunc(), sqlInitFunc(), telegramInitFunc() (+2 more)

### Community 41 - "Folder Loader Options"
Cohesion: 0.29
Nodes (8): opt, opts, Load(), load internal helper function, matchesAnyPattern(), opts struct, WithIgnore function, WithIgnore()

### Community 42 - "Git Status Diff"
Cohesion: 0.05
Nodes (33): Context, Changes, CommitWithUntrackedAction, gitChangesType, InitGit, InstallHooksAction, StatusDiff, Tool (+25 more)

### Community 43 - "Build Project Action"
Cohesion: 0.36
Nodes (9): commonProjectTidyPostActions(), commonProjectTidyPreActions(), GetTidyActionsForProject(), Action, goProjectTidyActions(), T, Test_GetTidyActionsForProject_Fast(), Test_GetTidyActionsForProject_UnknownType() (+1 more)

### Community 44 - "ProjEnv Config Access"
Cohesion: 0.22
Nodes (6): NewCommand(), newTidyCmd(), projectTidy.run() method, ActionPerformer interface, NewActionPerformer(), projectTidy

### Community 45 - "Migration Tool Interface"
Cohesion: 0.38
Nodes (9): Postgres struct, postgresClient(), countOccurrences(), T, Test_Postgres_AppendToProject_Idempotent(), Test_Postgres_AppendToProject_MultipleInstances(), Test_Postgres_AppendToProject_Single(), testVervConfig() (+1 more)

### Community 46 - "RW File Locking Utility"
Cohesion: 0.25
Nodes (3): Mutex, Reader, RW

### Community 47 - "IOMock Error Expectations"
Cohesion: 0.32
Nodes (3): IOMockErrorExpectation, IOMockErrorParams, mIOMockError

### Community 48 - "IProjectMock GetConfig Expectations"
Cohesion: 0.16
Nodes (6): expected grpc config_template.yaml fixture, expected grpc dev.yaml fixture, Config, IProjectMockGetConfigExpectation, IProjectMockGetConfigResults, mIProjectMockGetConfig

### Community 49 - "Env Config/Variables Fetching"
Cohesion: 0.36
Nodes (7): Action, ActionPerformer, IActionPerformer, PipelineLabels, Proc, RunPipeline(), stepEmoji()

### Community 51 - "Redis Dependency Client"
Cohesion: 0.46
Nodes (3): DataSourcesConfig, Redis, Project

### Community 52 - "SQL Connection (pattern)"
Cohesion: 0.25
Nodes (6): Adding a new `add` dependency, Architecture, Command flow, Commands, Key packages, Test helpers

### Community 53 - "Git Init Action"
Cohesion: 0.29
Nodes (3): newTidyEnvCmd(), envTidy, GetWd function

### Community 54 - "RW Read/Execute"
Cohesion: 0.29
Nodes (5): Sqlite, Project, containsDependency(), DataSources, Resource

### Community 55 - "IOMock Print Expectations"
Cohesion: 0.28
Nodes (4): IOMockPrintExpectation, IOMockPrintParams, mIOMockPrint, RWMutex

### Community 56 - "IOMock Println Expectations"
Cohesion: 0.32
Nodes (3): IOMockPrintlnExpectation, IOMockPrintlnParams, mIOMockPrintln

### Community 57 - "Generator Tests & Folder Comparison"
Cohesion: 0.43
Nodes (6): T, TestGenServer(), AssertFolderInFs(), AssertVirtualFolder(), CompareLongStrings(), T

### Community 58 - "SQL Connection (pattern, dup)"
Cohesion: 0.38
Nodes (6): PrepareClients, Postgres.AppendToProject, Redis.AppendToProject, Redis.applyClientFolder, Redis.applyConfig, Redis.GetFolderName

### Community 59 - "Environment Aggregate Structs"
Cohesion: 0.29
Nodes (7): Makefile struct, PortManager struct, GlobalEnvironment struct, envConfig struct, envMakefile struct, envVariables struct, ProjEnv struct

### Community 60 - "Port Manager"
Cohesion: 0.48
Nodes (6): actionPerformer.Tidy, commonProjectTidyPostActions, commonProjectTidyPreActions, GetTidyActionsForProject, goProjectTidyActions, unknownProjectActions

### Community 61 - "IOMock PrintColored Expectations"
Cohesion: 0.29
Nodes (7): Acceptance Criteria, Context, Do NOT change, Files to Create / Modify, Goal, Notes, Task 001 — Pattern system: embed and copy static template files directly

### Community 62 - "Env Variables Manager"
Cohesion: 0.67
Nodes (4): InitDepFuncGenArgs, InitFuncCall, InitServerListenerArgs, InitServerListenersArgs

### Community 63 - "Env Install Command (unwired)"
Cohesion: 0.22
Nodes (6): newEnvInstallCmd(), envInstall, Mutex, NewSpinner(), IO, Spinner

### Community 65 - "Git Commit Action"
Cohesion: 0.53
Nodes (4): EnumGenArg, generalGenArgs, KeyValue, newConfigStructGenArgs()

### Community 66 - "Project Loader"
Cohesion: 0.27
Nodes (7): Project, goProjectLoader(), LoadProject(), LoadProjectConfig(), readIgnoredFiles(), unknownProjectLoader(), envConfig

### Community 67 - "gRPC Dependency Test Fixtures"
Cohesion: 0.33
Nodes (3): GenerateMain(), T, Test_GenerateMain()

### Community 68 - "gRPC Implementation Registration"
Cohesion: 0.50
Nodes (5): grpcServer.AddImplementation() method, GrpcImpl interface, GrpcWithGateway interface, Api_grpc key, Impl struct (example gRPC service implementation)

### Community 69 - "Server Start Lifecycle"
Cohesion: 0.40
Nodes (5): grpcServer.start() method, httpServer.AddHttpHandler method, httpServer.buildHomePageHandler method, httpServer.start method, ServersManager.Start method

### Community 71 - "File Server Generator"
Cohesion: 0.25
Nodes (6): FS, fileServerTemplate (fs.go.pattern template), GenerateFileServer(), T, Test_GenerateFileServer(), fileServerGenArgs

### Community 72 - "Proto API Generator"
Cohesion: 0.25
Nodes (5): serviceProtoApiArgs, GenerateServiceApiProto(), T, Test_GenerateServiceApiProto(), basicApiProtoTemplate (api.proto.pattern template)

### Community 73 - "Make Binary Installer"
Cohesion: 0.50
Nodes (3): Command, NewCommand(), projectTidy

### Community 74 - "IOMock GetInput Builder"
Cohesion: 0.23
Nodes (6): IOMockGetInputExpectation, IOMockGetInputResults, IOMockPrintlnColoredExpectation, IOMockPrintlnColoredParams, mIOMockGetInput, mIOMockPrintlnColored

### Community 75 - "IProjectMock GetName Builder"
Cohesion: 0.40
Nodes (4): ErrServerMustHaveName error var, generateServerInitFileAndArgs(), Servers, initServerTemplate (init_server.go.pattern template)

### Community 78 - "GlobalEnvironment Init"
Cohesion: 0.40
Nodes (5): GlobalEnvironment.getSpirits, GlobalEnvironment.Init, GlobalEnvironment.initBasis, GlobalEnvironment.initProjectDir, GlobalEnvironment.initProjectsDirs

### Community 79 - "GitHub Workflow Tidy Action"
Cohesion: 0.40
Nodes (5): TidyGithubWorkflowAction.Do, GithubWorkflowGoBranchPush embedded template, GithubWorkflowRelease embedded template, GithubFolder const, WorkflowsFolder const

### Community 80 - "Generator Arg Structs"
Cohesion: 0.40
Nodes (5): AppContent struct, InitDepFuncGenArgs struct, EnumGenArg struct, generalGenArgs struct, KeyValue struct

### Community 81 - "CI / CLAUDE.md Dev Commands"
Cohesion: 0.50
Nodes (4): CLAUDE.md Development Commands Section, .claude/settings.local.json permissions config, build-and-test job (branch-push workflow), create-pr job (branch-push workflow)

### Community 82 - "Environment Plugin Overview"
Cohesion: 0.50
Nodes (4): plugins/environment package description, envTidy.getEnvDirPath() method, envTidy.RunTidy() method, Environment Setup Feature (scans sibling projects, generates docker-compose)

### Community 84 - "Main Entrypoint (external app)"
Cohesion: 0.50
Nodes (3): app.New (inferred external constructor), app.Start (inferred external method), main() entrypoint

### Community 86 - "Server Stop Lifecycle"
Cohesion: 0.67
Nodes (3): grpcServer.stop() method, httpServer.stop method, ServersManager.Stop method

### Community 88 - "Progress Loader Interface"
Cohesion: 0.67
Nodes (3): InfiniteLoader struct, Progress interface, percentLoader struct

### Community 90 - "Project Interface Duplication"
Cohesion: 0.67
Nodes (3): Project interface (go_actions pkg), IProject interface, Project struct

### Community 105 - "Embedded Template Registration"
Cohesion: 0.67
Nodes (3): T, Test_generateDataSourceInitFileAndArgs_PostgresExcluded(), Test_generateDataSourceInitFileAndArgs_PostgresOnly()

### Community 209 - "Project Name Substitution"
Cohesion: 0.15
Nodes (12): Config example can be found in internal/config/verv.yaml, Configuration, Configuration structure (all fields are optional), Features:, Installation, SPECIAL VARS, THE LAST, NOT THE LEAST on project creation, TODO (+4 more)

## Ambiguous Edges - Review These
- `main() entrypoint` → `app.New (inferred external constructor)`  [AMBIGUOUS]
  plugins/project/go_project/patterns/pattern/cmd/service/main.go · relation: calls
- `main() entrypoint` → `app.Start (inferred external method)`  [AMBIGUOUS]
  plugins/project/go_project/patterns/pattern/cmd/service/main.go · relation: calls
- `Progress interface` → `percentLoader struct`  [AMBIGUOUS]
  internal/io/loader/percent_progress.go · relation: implements

## Knowledge Gaps
- **286 isolated node(s):** `Proc`, `Proc`, `go.vervstack.ru/verv`, `Project`, `opt` (+281 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **112 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `main() entrypoint` and `app.New (inferred external constructor)`?**
  _Edge tagged AMBIGUOUS (relation: calls) - confidence is low._
- **What is the exact relationship between `main() entrypoint` and `app.Start (inferred external method)`?**
  _Edge tagged AMBIGUOUS (relation: calls) - confidence is low._
- **What is the exact relationship between `Progress interface` and `percentLoader struct`?**
  _Edge tagged AMBIGUOUS (relation: implements) - confidence is low._
- **Why does `IO` connect `Env Install Command (unwired)` to `RsCli Config Loading`, `HTTP Server Transport (pattern)`, `Project Name Collection Prompt`, `Project Action Pipeline`, `gRPC Package Discovery`, `Environment File Fetching`, `Env Tidy & Terminal Loader UI`, `Make Binary Installer`, `IO Mock (minimock)`, `ProjEnv Config Access`, `CLI Command Wiring (environment)`, `Env Config/Variables Fetching`, `Generated Project Config Loader`, `Project Init Command`, `Git Init Action`, `SQL Connection (pattern, dup)`?**
  _High betweenness centrality (0.152) - this node is a cross-community bridge._
- **Why does `IProject` connect `Config Folder Generation` to `gRPC Package Discovery`, `Dockerfile Generation & Build Tests`, `Project Structure Preparation Actions`, `IProject Mock (minimock)`, `Project Init Command`, `gRPC Server Transport (pattern)`, `Go Fmt & Makefile Gen Actions`, `Project Name Collection Prompt`, `IProjectMock GetType Expectations`, `Multiplexed Server Manager`, `Git Status Diff`, `IProjectMock GetConfig Expectations`, `Env Config/Variables Fetching`, `Generator Tests & Folder Comparison`, `gRPC Dependency Test Fixtures`, `Proto API Generator`, `IProjectMock GetProjectPath Builder`, `Go Mod Init Action`, `Telegram Bot Connection`, `Embedded Template Registration`?**
  _High betweenness centrality (0.142) - this node is a cross-community bridge._
- **Why does `RsCliConfig` connect `RsCli Config Loading` to `Project Action Pipeline`, `SQL Connection Pattern Wiring`, `gRPC Package Discovery`, `Project Name Collection Prompt`, `Environment File Fetching`, `Project Loader`, `Dockerfile Generation & Build Tests`, `CLI Command Wiring (environment)`, `ProjEnv Config Access`, `Project Init Command`, `Dependency AppendToProject (Sqlite/Env)`, `Git Init Action`, `SQL Connection (pattern, dup)`, `Env Install Command (unwired)`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **What connects `Proc`, `Proc`, `go.vervstack.ru/verv` to the rest of the system?**
  _287 weakly-connected nodes found - possible documentation gaps or missing edges._