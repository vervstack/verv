# Graph Report - .  (2026-07-14)

## Corpus Check
- Corpus is ~44,230 words - fits in a single context window. You may not need a graph.

## Summary
- 1533 nodes · 2338 edges · 211 communities (104 shown, 107 thin omitted)
- Extraction: 87% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 292 edges (avg confidence: 0.82)
- Token cost: 764,541 input · 0 output

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

## God Nodes (most connected - your core abstractions)
1. `IOMock` - 74 edges
2. `IProjectMock` - 65 edges
3. `IProject` - 39 edges
4. `RsCliConfig` - 33 edges
5. `Folder` - 26 edges
6. `Test_AddDependency()` - 25 edges
7. `IO` - 23 edges
8. `Execute function` - 21 edges
9. `Color` - 18 edges
10. `ActionPerformerMock` - 18 edges

## Surprising Connections (you probably didn't know these)
- `rscli Architecture Overview (Cobra CLI, two top-level commands)` --references--> `NewCmd() (env command constructor)`  [INFERRED]
  CLAUDE.md → cmd/environment/cmd.go
- `rscli Architecture Overview (Cobra CLI, two top-level commands)` --references--> `NewCmd() (project command constructor)`  [INFERRED]
  CLAUDE.md → cmd/project/cmd.go
- `InitConfig function` --semantically_similar_to--> `Load function`  [INFERRED] [semantically similar]
  internal/config/cfg.go → cmd/project/init_new/expected/load_config/load_config_file.go
- `main()` --references--> `go.mod module manifest (github.com/Red-Sock/rscli)`  [INFERRED]
  main.go → go.mod
- `pattern go.mod module dependencies` --references--> `setUpCors function`  [EXTRACTED]
  plugins/project/go_project/patterns/pattern/go.mod → cmd/project/add/expected/grpc/internal/transport/http.go

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
- **Structs implementing the actions.Action interface** — plugins_project_actions_actions_action, plugins_project_actions_go_actions_build_buildprojectaction, plugins_project_actions_go_actions_config_generateprojectconfig, plugins_project_actions_go_actions_config_prepareconfigfolder, plugins_project_actions_go_actions_go_mod_gofmt, plugins_project_actions_go_actions_preparation_prepareprojectstructure, plugins_project_actions_git_git_initgit, plugins_project_actions_git_commit_commitwithuntrackedaction [INFERRED 0.85]
- **Tidy action pipeline assembled by GetTidyActionsForProject** — plugins_project_actions_tidy_project_gettidyactionsforproject, plugins_project_actions_tidy_project_goprojecttidyactions, plugins_project_actions_tidy_project_commonprojecttidypreactions, plugins_project_actions_tidy_project_commonprojecttidypostactions, plugins_project_actions_go_actions_preparation_prepareclients, plugins_project_actions_go_actions_preparation_prepareserver, plugins_project_actions_git_commit_commitwithuntrackedaction [EXTRACTED 1.00]
- **Go project initialization pipeline assembled by initVirtualGoProject** — plugins_project_actions_init_project_initvirtualgoproject, plugins_project_actions_go_actions_preparation_prepareprojectstructure, plugins_project_actions_go_actions_init_initgoprojectapp, plugins_project_actions_go_actions_config_generateprojectconfig, plugins_project_actions_go_actions_config_prepareconfigfolder, plugins_project_actions_go_actions_preparation_preparemakefile, plugins_project_actions_go_actions_preparation_prepareclients, plugins_project_actions_go_actions_preparation_prepareserver, plugins_project_actions_go_actions_build_buildprojectaction, plugins_project_actions_go_actions_init_initgomod, plugins_project_actions_go_actions_go_mod_gofmt, plugins_project_actions_git_git_initgit [EXTRACTED 1.00]
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
- **Minimock-generated interface mocks (ActionPerformer, IProject, IO)** — tests_mocks_action_performer_mock_actionperformermock, tests_mocks_i_project_mock_iprojectmock, tests_mocks_io_mock_iomock [INFERRED 0.85]
- **MockProject functional-options construction pattern** — tests_project_mock_prepare_getmockproject, tests_project_mock_prepare_mockproject, tests_project_mock_opts_withfile, tests_project_mock_opts_withenvironmentvariables, tests_project_mock_opts_withfilesystem, tests_project_mock_opts_withbasicconfig, tests_project_mock_opts_withgit, tests_project_mock_opts_withsqlite, tests_project_mock_opts_withgrpcserver [INFERRED 0.85]
- **Test support chain for Dockerfile ENV var feature** — tasks_task_002_dockerfile_env_vars, tests_project_mock_config_getallenvvariables, tests_project_mock_opts_withenvironmentvariables [INFERRED 0.75]

## Communities (211 total, 107 thin omitted)

### Community 0 - "Project Action Pipeline"
Cohesion: 0.06
Nodes (51): Action, ActionPerformer, Proc, testCase, Dependency Names (grpc, redis, postgres, telegram, sqlite, env), endMsg constant, preparingMsg constant, startingMsg constant (+43 more)

### Community 1 - "gRPC Package Discovery"
Cohesion: 0.05
Nodes (43): GrpcDiscovery, GrpcPackage, genArgs, KebabToSnake, SnakeToPascal, ToPascal, slices.Contains, slices.Exclude (+35 more)

### Community 2 - "Dockerfile Generation & Build Tests"
Cohesion: 0.07
Nodes (37): serviceDockerfileArgs, BuildProjectSuite, InitGoProjectApp, PrepareDockerfile, DirEntry, T, fullConfigGoFile, T (+29 more)

### Community 3 - "App File Generator"
Cohesion: 0.07
Nodes (43): AppContent, AppFileGenArgs, AppStarter, InitDepFuncGenArgs, InitFuncCall, InitServerListenerArgs, InitServerListenersArgs, EnumGenArg (+35 more)

### Community 4 - "Environment File Fetching"
Cohesion: 0.07
Nodes (31): Container, Variable, MewEmptyMakefile, NewMakeFile, parseRule, parseVariable, ReadMakeFile, PortManager.SaveIfNotExist (+23 more)

### Community 5 - "Compiled Pattern Scaffold Files"
Cohesion: 0.06
Nodes (45): easyp.yaml proto codegen config (pattern_c), branch-push GitHub Actions workflow, master-actions RELEASE GitHub Actions workflow, .golangci.yaml lint config (pattern_c), README.md template (pattern_c), RedSock CLI (referenced tool), telegram version Handler struct, proj_name_apiPingRequest swagger definition (+37 more)

### Community 6 - "Config Folder Generation"
Cohesion: 0.08
Nodes (21): GenerateProjectConfig, PrepareConfigFolder, ActionPerformerMock, ActionPerformerMockTidyExpectation, ActionPerformerMockTidyParams, ActionPerformerMockTidyResults, mActionPerformerMockTidy, Node (+13 more)

### Community 7 - "Env Tidy & Terminal Loader UI"
Cohesion: 0.07
Nodes (23): Color type, TerminalColor function, InfiniteLoader.Done method, Mutex, NewInfiniteLoader constructor, Context, RunMultiLoader function, percentLoader.GetLoaderSymb method (+15 more)

### Community 8 - "Project Structure Preparation Actions"
Cohesion: 0.06
Nodes (27): PrepareClients, PrepareMakefile, PrepareProjectStructure, PrepareServer, PortManager.GetNextPort, renamer.ReplaceProjectName, ReplaceProjectNameFull, ReplaceProjectNameShort (+19 more)

### Community 9 - "gRPC Client/Gateway Interfaces"
Cohesion: 0.08
Nodes (34): ClientConnInterface, Request, DialOption, Impl, ProjNameAPIClient, ProjNameAPIServer, UnimplementedProjNameAPIServer, UnsafeProjNameAPIServer (+26 more)

### Community 10 - "IO Mock (minimock)"
Cohesion: 0.07
Nodes (4): IOMock, Duration, Tester, NewIOMock()

### Community 11 - "Docker Compose Assembly"
Cohesion: 0.08
Nodes (26): Compose, ContainerSettings, Pattern, PatternManager, compose.examples.yaml template, docker-compose.yaml template, ContainerSettings struct, AddEnvironmentBrackets function (+18 more)

### Community 12 - "CLI Command Wiring (environment)"
Cohesion: 0.08
Nodes (18): rscli Architecture Overview (Cobra CLI, two top-level commands), Command, NewCmd() (env command constructor), Command, newTidyEnvCmd(), Command, NewCmd() (project command constructor), Command (+10 more)

### Community 14 - "Postgres Client & Config"
Cohesion: 0.08
Nodes (24): Postgres master config.yaml, Postgres dev.yaml config, DB, DB, sqlLogger, SqlResource, sqldb.New (SQL connection + goose migration) function, postgres.go blank import of lib/pq driver (+16 more)

### Community 15 - "Generated Proto Validation (Ping)"
Cohesion: 0.08
Nodes (6): PingRequestMultiError, PingRequestValidationError, PingResponseMultiError, PingResponseValidationError, PingRequest, PingResponse

### Community 16 - "gRPC Server Transport"
Cohesion: 0.10
Nodes (19): Context, Listener, ServeMux, ServerOption, GrpcImpl, grpcServer, GrpcWithGateway, grpcServer struct (+11 more)

### Community 17 - "TS gRPC-Gateway Fetch Client"
Cohesion: 0.10
Nodes (20): b64, b64Encode(), fetchStreamingRequest(), FlattenedRequestPayload, flattenRequestPayload(), getNewLineDelimitedJSONDecodingStream(), getNotifyEntityArrivalSink(), InitReq (+12 more)

### Community 18 - "Project Init Command"
Cohesion: 0.12
Nodes (16): projectInit.buildProject() method, Command, Proc, newInitCmd(), projectInit.obtainFolderPathFromUser() method, projectInit.obtainNameFromUser() method, projectInit.run() method, Action (+8 more)

### Community 19 - "Terminal Color Parser & IO Stub"
Cohesion: 0.11
Nodes (7): Attribute, Color, UIColor function, IoDevNul, IOMockPrintlnColoredExpectation, IOMockPrintlnColoredParams, mIOMockPrintlnColored

### Community 20 - "Generated Protobuf Message Methods"
Cohesion: 0.12
Nodes (10): MessageState, file_grpc_api_proto_init(), file_grpc_api_proto_rawDescGZIP(), PingRequest, PingResponse, Message, init(), SizeCache (+2 more)

### Community 21 - "Dependency AppendToProject (Sqlite/Env)"
Cohesion: 0.15
Nodes (16): Dependency, EnvVariable, Sqlite struct, sqlite() constructor, Dependency interface, dependencyBase struct, GetDependencies func, HelpWithDependencyNames func (+8 more)

### Community 22 - "IProjectMock GetFolder Expectations"
Cohesion: 0.17
Nodes (11): IProjectMockGetFolderExpectation, IProjectMockGetFolderResults, IProjectMockGetNameExpectation, IProjectMockGetNameResults, IProjectMockGetProjectPathExpectation, IProjectMockGetProjectPathResults, IProjectMockGetShortNameExpectation, IProjectMockGetShortNameResults (+3 more)

### Community 23 - "Telegram Dependency Wiring"
Cohesion: 0.27
Nodes (9): Telegram, containsDependencyFolder func, Project, Telegram.AppendToProject, Telegram.applyClient, Telegram.applyConfig, Telegram.applyFolder, ReplaceProjectName func (+1 more)

### Community 24 - "Sqlite Dependency Client"
Cohesion: 0.14
Nodes (12): Sqlite, Project, Sqlite.AppendToProject, containsDependency func, DataSources, Resource, EnvVariable.AppendToProject, Project interface (dependencies pkg) (+4 more)

### Community 25 - "IOMock GetInput Expectations"
Cohesion: 0.22
Nodes (6): IOMockGetInputExpectation, IOMockGetInputOneOfExpectation, IOMockGetInputOneOfParams, IOMockGetInputOneOfResults, IOMockGetInputResults, mIOMockGetInputOneOf

### Community 26 - "gRPC Server Dependency Wiring"
Cohesion: 0.26
Nodes (8): GrpcServer, Project, Project, GrpcServer.addGrpcServerToConfig, GrpcServer.AppendToProject, GrpcServer.applyApiFolder, initServerManagerFiles func, prepareServerConfig func

### Community 27 - "Environment Init & File IO"
Cohesion: 0.24
Nodes (6): CreateFileIfNotExists function, CreateFolderIfNotExists function, OverrideFile function, DirEntry, GlobalEnvironment, T

### Community 28 - "gRPC Server Transport (pattern)"
Cohesion: 0.22
Nodes (8): Context, Listener, ServeMux, ServerOption, GrpcImpl, grpcServer, GrpcWithGateway, newGrpcServer()

### Community 29 - "Virtual Folder Tree"
Cohesion: 0.18
Nodes (3): Folder, Client, New()

### Community 30 - "Go Fmt & Makefile Gen Actions"
Cohesion: 0.23
Nodes (6): GoFmt, RunGoTidyAction, RunMakeGenAction, UpdateAllPackages, Action interface (actions pkg), Action interface (go_actions pkg)

### Community 31 - "HTTP Server Transport"
Cohesion: 0.29
Nodes (6): Cors, Listener, ServeMux, httpServer, newHttpServer function, setUpCors function

### Community 32 - "RsCli Config Loading"
Cohesion: 0.29
Nodes (10): Project, builtInConfig embedded var, GetConfig function, getConfigFromEnvironment function, getConfigFromFile function, Command, init(), InitConfig function (+2 more)

### Community 33 - "HTTP Server Transport (pattern)"
Cohesion: 0.29
Nodes (6): Cors, Listener, ServeMux, httpServer, newHttpServer(), setUpCors()

### Community 34 - "Project Name Collection Prompt"
Cohesion: 0.20
Nodes (9): nameCollector.askUserForName method, nameCollector.collect method, nameCollector.preAppendHost method, nameCollector.removeHttpProtoc method, T, Test_collectName test function, Proc.collectOsPath method, Proc.run method (+1 more)

### Community 35 - "Postgres Dependency Client"
Cohesion: 0.27
Nodes (7): dependencyBase, Postgres, Project, Postgres struct, postgresClient, Redis struct, redisClient

### Community 36 - "IProjectMock GetType Expectations"
Cohesion: 0.24
Nodes (4): IProjectMockGetTypeExpectation, IProjectMockGetTypeResults, mIProjectMockGetType, Type

### Community 37 - "Multiplexed Server Manager"
Cohesion: 0.24
Nodes (7): CMux, Context, grpcServer, httpServer, Listener, ServersManager, NewServerManager()

### Community 38 - "Telegram Version Handler"
Cohesion: 0.25
Nodes (5): Chat, MessageIn, Config, version.New constructor, Handler

### Community 39 - "Telegram Transport Listener"
Cohesion: 0.31
Nodes (5): Bot, Config, Context, NewServer (telegram transport), Server

### Community 40 - "Environment Config Struct"
Cohesion: 0.31
Nodes (6): RsCliConfig, Project, AppConfig, Config, LoadProjectConfig, envConfig

### Community 41 - "Folder Loader Options"
Cohesion: 0.33
Nodes (7): opt, opts, Load function (folder loader), load internal helper function, matchesAnyPattern function, opts struct, WithIgnore function

### Community 42 - "Git Status Diff"
Cohesion: 0.36
Nodes (5): Changes, gitChangesType, StatusDiff, Changes struct, Status

### Community 43 - "Build Project Action"
Cohesion: 0.22
Nodes (4): BuildProjectAction, Run (make bin), BuildProjectSuite.Test_BuildProject, Test_BuildProject

### Community 45 - "Migration Tool Interface"
Cohesion: 0.28
Nodes (4): Tool, ErrUnknownResourceToMigrate sentinel error, Resource, GithubVersion struct

### Community 46 - "RW File Locking Utility"
Cohesion: 0.25
Nodes (3): Mutex, Reader, RW

### Community 47 - "IOMock Error Expectations"
Cohesion: 0.28
Nodes (4): IOMockErrorExpectation, IOMockErrorParams, mIOMockError, RWMutex

### Community 48 - "IProjectMock GetConfig Expectations"
Cohesion: 0.28
Nodes (4): IProjectMockGetConfigExpectation, IProjectMockGetConfigResults, mIProjectMockGetConfig, Config

### Community 49 - "Env Config/Variables Fetching"
Cohesion: 0.28
Nodes (6): envConfig.fetch, envConfig.findEnvConfig, envVariables.fetch, Container, ProjEnv, LoadProjectEnvironment

### Community 50 - "Generated Project Config Loader"
Cohesion: 0.29
Nodes (7): AppInfo, EnvironmentConfig struct, Config struct, Load function, EnvironmentConfig, Config, ServiceDiscovery

### Community 51 - "Redis Dependency Client"
Cohesion: 0.46
Nodes (3): DataSourcesConfig, Redis, Project

### Community 52 - "SQL Connection (pattern)"
Cohesion: 0.25
Nodes (5): DB, DB, sqlLogger, SqlResource, New()

### Community 53 - "Git Init Action"
Cohesion: 0.32
Nodes (4): InitGit, git.Init, Init, SetOrigin

### Community 54 - "RW Read/Execute"
Cohesion: 0.29
Nodes (5): Execute function, RW struct, RW.Read method, RW.String method, FetchPackage func

### Community 55 - "IOMock Print Expectations"
Cohesion: 0.32
Nodes (3): IOMockPrintExpectation, IOMockPrintParams, mIOMockPrint

### Community 56 - "IOMock Println Expectations"
Cohesion: 0.32
Nodes (3): IOMockPrintlnExpectation, IOMockPrintlnParams, mIOMockPrintln

### Community 57 - "Generator Tests & Folder Comparison"
Cohesion: 0.43
Nodes (6): T, TestGenServer test, AssertFolderInFs(), AssertVirtualFolder(), CompareLongStrings(), T

### Community 58 - "SQL Connection (pattern, dup)"
Cohesion: 0.25
Nodes (5): DB, DB, sqlLogger, SqlResource, New()

### Community 59 - "Environment Aggregate Structs"
Cohesion: 0.29
Nodes (7): Makefile struct, PortManager struct, GlobalEnvironment struct, envConfig struct, envMakefile struct, envVariables struct, ProjEnv struct

### Community 60 - "Port Manager"
Cohesion: 0.33
Nodes (3): Mutex, NewPortManager, PortManager

### Community 61 - "IOMock PrintColored Expectations"
Cohesion: 0.38
Nodes (3): IOMockPrintColoredExpectation, IOMockPrintColoredParams, mIOMockPrintColored

### Community 62 - "Env Variables Manager"
Cohesion: 0.52
Nodes (4): Container, newEnvManager, envResources, envVariables

### Community 63 - "Env Install Command (unwired)"
Cohesion: 0.40
Nodes (3): Command, newEnvInstallCmd(), envInstall

### Community 65 - "Git Commit Action"
Cohesion: 0.53
Nodes (3): CommitWithUntrackedAction, Commit, CommitWithUntracked

### Community 66 - "Project Loader"
Cohesion: 0.53
Nodes (5): Project, goProjectLoader, LoadProject, readIgnoredFiles, unknownProjectLoader

### Community 67 - "gRPC Dependency Test Fixtures"
Cohesion: 0.50
Nodes (4): expected grpc config_template.yaml fixture, expected grpc dev.yaml fixture, Config, AppConfig

### Community 68 - "gRPC Implementation Registration"
Cohesion: 0.50
Nodes (5): grpcServer.AddImplementation() method, GrpcImpl interface, GrpcWithGateway interface, Api_grpc key, Impl struct (example gRPC service implementation)

### Community 69 - "Server Start Lifecycle"
Cohesion: 0.40
Nodes (5): grpcServer.start() method, httpServer.AddHttpHandler method, httpServer.buildHomePageHandler method, httpServer.start method, ServersManager.Start method

### Community 71 - "File Server Generator"
Cohesion: 0.40
Nodes (4): FS, fileServerTemplate (fs.go.pattern template), GenerateFileServer function, fileServerGenArgs

### Community 72 - "Proto API Generator"
Cohesion: 0.40
Nodes (3): serviceProtoApiArgs, GenerateServiceApiProto function, basicApiProtoTemplate (api.proto.pattern template)

### Community 73 - "Make Binary Installer"
Cohesion: 0.70
Nodes (4): Exists (make bin check), Install (make bin), installLinux, installMacOS

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

### Community 85 - "gRPC Client Connect Helper"
Cohesion: 0.50
Nodes (3): connect (grpc client dial helper), ClientConn, Resource_GRPC_Rscli_example key

### Community 86 - "Server Stop Lifecycle"
Cohesion: 0.67
Nodes (3): grpcServer.stop() method, httpServer.stop method, ServersManager.Stop method

### Community 88 - "Progress Loader Interface"
Cohesion: 0.67
Nodes (3): InfiniteLoader struct, Progress interface, percentLoader struct

### Community 90 - "Project Interface Duplication"
Cohesion: 0.67
Nodes (3): Project interface (go_actions pkg), IProject interface, Project struct

## Ambiguous Edges - Review These
- `main() entrypoint` → `app.New (inferred external constructor)`  [AMBIGUOUS]
  plugins/project/go_project/patterns/pattern/cmd/service/main.go · relation: calls
- `main() entrypoint` → `app.Start (inferred external method)`  [AMBIGUOUS]
  plugins/project/go_project/patterns/pattern/cmd/service/main.go · relation: calls
- `Progress interface` → `percentLoader struct`  [AMBIGUOUS]
  internal/io/loader/percent_progress.go · relation: implements

## Knowledge Gaps
- **261 isolated node(s):** `GrpcWithGateway`, `DB`, `DB`, `Proc`, `EnvironmentConfig` (+256 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **107 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `main() entrypoint` and `app.New (inferred external constructor)`?**
  _Edge tagged AMBIGUOUS (relation: calls) - confidence is low._
- **What is the exact relationship between `main() entrypoint` and `app.Start (inferred external method)`?**
  _Edge tagged AMBIGUOUS (relation: calls) - confidence is low._
- **What is the exact relationship between `Progress interface` and `percentLoader struct`?**
  _Edge tagged AMBIGUOUS (relation: implements) - confidence is low._
- **Why does `RsCliConfig` connect `Environment Config Struct` to `Project Action Pipeline`, `SQL Connection Pattern Wiring`, `gRPC Package Discovery`, `Postgres Dependency Client`, `RsCli Config Loading`, `Environment File Fetching`, `Project Loader`, `Dockerfile Generation & Build Tests`, `Project Structure Preparation Actions`, `CLI Command Wiring (environment)`, `Env Config/Variables Fetching`, `Project Init Command`, `Dependency AppendToProject (Sqlite/Env)`, `Go Fmt & Makefile Gen Actions`, `Env Install Command (unwired)`?**
  _High betweenness centrality (0.167) - this node is a cross-community bridge._
- **Why does `IO` connect `Project Action Pipeline` to `gRPC Package Discovery`, `Environment File Fetching`, `Env Tidy & Terminal Loader UI`, `Project Structure Preparation Actions`, `IO Mock (minimock)`, `CLI Command Wiring (environment)`, `Project Init Command`, `Terminal Color Parser & IO Stub`, `Go Fmt & Makefile Gen Actions`, `Env Install Command (unwired)`?**
  _High betweenness centrality (0.126) - this node is a cross-community bridge._
- **Why does `Folder` connect `Virtual Folder Tree` to `gRPC Package Discovery`, `App File Generator`, `Environment File Fetching`, `Config Folder Generation`, `File Server Generator`, `Proto API Generator`, `Folder Loader Options`, `ProjEnv Config Access`, `IProject Mock (minimock)`, `IProjectMock GetFolder Expectations`, `Telegram Dependency Wiring`, `Generator Tests & Folder Comparison`, `Environment Init & File IO`?**
  _High betweenness centrality (0.116) - this node is a cross-community bridge._
- **What connects `GrpcWithGateway`, `DB`, `DB` to the rest of the system?**
  _262 weakly-connected nodes found - possible documentation gaps or missing edges._