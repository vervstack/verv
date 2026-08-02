# Graph Report - verv  (2026-08-02)

## Corpus Check
- 135 files · ~31,215 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1034 nodes · 1452 edges · 228 communities (53 shown, 175 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 194 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3f3bb4a5`
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
- Postgres Driver Import
- Unexpected Ping Error
- Sqlite Driver Import
- Init Proc Struct
- Collect OS Path
- Project Assembler
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
- Environment Init Entry
- Env Resources Type
- Tidy API
- Tidy Migrations
- Tidy Resources
- Tidy Service
- Environment Tidy Entry
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
- Common Imports Vars
- Errors Package Import Name
- Connection Name Constants
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
- EasyP Config Template
- gRPC Impl Template
- gRPC Impl Folder Constant
- gRPC Server Makefile Template
- Proto Contract Template
- Golangci Lint Config
- Postgres Driver Import (dup)
- Sqlite Driver Import (dup)
- Environment Config Struct (dup)
- Config Keys File
- Telegram Server Start
- Telegram Server Stop
- Version Command Handler
- Tools List Entry
- Project Name Substitution
- Mock Project WriteFile
- proj_name
- Proc
- pre-commit
- pre-commit
- go.vervstack.ru/verv
- AppConfig
- Environment Setup Feature (scans sibling projects, generates docker-compose)

## God Nodes (most connected - your core abstractions)
1. `IOMock` - 71 edges
2. `IProjectMock` - 63 edges
3. `IProject` - 36 edges
4. `Folder` - 32 edges
5. `VervConfig` - 20 edges
6. `Execute()` - 19 edges
7. `Color` - 17 edges
8. `ActionPerformerMock` - 17 edges
9. `New()` - 16 edges
10. `IO` - 14 edges

## Surprising Connections (you probably didn't know these)
- `Config` --shares_data_with--> `expected grpc config_template.yaml fixture`  [INFERRED]
  plugins/project/config/config.go → cmd/project/add/expected/grpc/config/config_template.yaml
- `Config` --shares_data_with--> `expected grpc dev.yaml fixture`  [INFERRED]
  plugins/project/config/config.go → cmd/project/add/expected/grpc/config/dev.yaml
- `NormalizeResourceName()` --calls--> `SnakeToPascal()`  [INFERRED]
  plugins/project/go_project/patterns/generators/config.go → internal/utils/cases/cases.go
- `Test_InitGit_Do_SkipCommit()` --calls--> `Execute()`  [INFERRED]
  plugins/project/actions/git/git_test.go → internal/cmd/executer.go
- `FetchPackage()` --calls--> `Execute()`  [INFERRED]
  plugins/project/actions/go_actions/dependencies/link_service/grpc_discovery/fetcher.go → internal/cmd/executer.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared app_info/data_sources/environment config.yaml schema across dependencies** — cmd_project_add_expected_postgres_config_config_postgres, cmd_project_add_expected_redis_config_config_redis, cmd_project_add_expected_sqlite_config_config_sqlite, cmd_project_add_expected_telegram_config_config_telegram [INFERRED 0.75]

## Communities (228 total, 175 thin omitted)

### Community 0 - "Project Action Pipeline"
Cohesion: 0.08
Nodes (15): Command, NewCmd(), New(), main(), run(), init(), init(), init() (+7 more)

### Community 1 - "gRPC Package Discovery"
Cohesion: 0.12
Nodes (16): File, FuncDecl, GenDecl, GrpcDiscovery, GrpcPackage, GrpcClient, GetPathToGlobalModule(), DiscoverPackage() (+8 more)

### Community 2 - "Dockerfile Generation & Build Tests"
Cohesion: 0.08
Nodes (31): DirEntry, BuildProjectSuite, OverrideFile(), T, Test_InitGit_Do_SkipCommit(), T, postgresClient(), countOccurrences() (+23 more)

### Community 4 - "Environment File Fetching"
Cohesion: 0.25
Nodes (7): Acceptance Criteria, Context, Do NOT change, Files to Create / Modify, Goal, Notes, Task 002 — Dockerfile: emit ENV directives from project config on tidy

### Community 6 - "Config Folder Generation"
Cohesion: 0.26
Nodes (5): ActionPerformerMockTidyExpectation, ActionPerformerMockTidyParams, ActionPerformerMockTidyResults, mActionPerformerMockTidy, RWMutex

### Community 10 - "IO Mock (minimock)"
Cohesion: 0.07
Nodes (4): IOMock, Duration, Tester, NewIOMock()

### Community 14 - "Postgres Client & Config"
Cohesion: 0.67
Nodes (3): Redis master config.yaml, Redis dev.yaml config, Redis config_template.yaml

### Community 18 - "Project Init Command"
Cohesion: 0.19
Nodes (11): Proc, Project, Action, InitProject(), initVirtualGoProject(), T, Test_InitProject_Fast(), Test_InitProject_UnknownType() (+3 more)

### Community 19 - "Terminal Color Parser & IO Stub"
Cohesion: 0.15
Nodes (6): Color, TerminalColor(), StdIO, IOMockPrintColoredExpectation, IOMockPrintColoredParams, mIOMockPrintColored

### Community 21 - "Dependency AppendToProject (Sqlite/Env)"
Cohesion: 0.08
Nodes (25): Dependency, dependencyBase, EnvVariable, Postgres, Redis, Sqlite, Telegram, DataSources (+17 more)

### Community 22 - "IProjectMock GetFolder Expectations"
Cohesion: 0.17
Nodes (11): IProjectMockGetFolderExpectation, IProjectMockGetFolderResults, IProjectMockGetNameExpectation, IProjectMockGetNameResults, IProjectMockGetProjectPathExpectation, IProjectMockGetProjectPathResults, IProjectMockGetShortNameExpectation, IProjectMockGetShortNameResults (+3 more)

### Community 25 - "IOMock GetInput Expectations"
Cohesion: 0.15
Nodes (7): IOMockGetInputExpectation, IOMockGetInputOneOfExpectation, IOMockGetInputOneOfParams, IOMockGetInputOneOfResults, IOMockGetInputResults, mIOMockGetInput, mIOMockGetInputOneOf

### Community 27 - "Environment Init & File IO"
Cohesion: 0.17
Nodes (24): PostgresInstance, PostgresInstancesArgs, execute(), executeWith(), GenerateGRPCConn(), GeneratePostgresConn(), GeneratePostgresDriver(), GeneratePostgresInstances() (+16 more)

### Community 28 - "gRPC Server Transport (pattern)"
Cohesion: 0.10
Nodes (25): AppContent, AppFileGenArgs, AppStarter, InitDepFuncGenArgs, InitFuncCall, InitServerListenerArgs, InitServerListenersArgs, T (+17 more)

### Community 29 - "Virtual Folder Tree"
Cohesion: 0.05
Nodes (20): Folder, FS, serviceProtoApiArgs, genArgs, SnakeToPascal(), ToPascal(), GenerateServiceApiProto(), T (+12 more)

### Community 30 - "Go Fmt & Makefile Gen Actions"
Cohesion: 0.17
Nodes (5): GoFmt, RunGoTidyAction, RunMakeGenAction, UpdateAllPackages, IProject

### Community 31 - "HTTP Server Transport"
Cohesion: 0.08
Nodes (22): PrepareClients, PrepareDockerfile, PrepareProjectStructure, PrepareServer, addMissingImplFolders(), generateTransportFiles(), GenerateMain(), T (+14 more)

### Community 32 - "RsCli Config Loading"
Cohesion: 0.11
Nodes (23): Project, VervConfig, sqlConn, opt, opts, GetConfig(), getConfigFromEnvironment(), getConfigFromFile() (+15 more)

### Community 33 - "HTTP Server Transport (pattern)"
Cohesion: 0.06
Nodes (31): Action, ActionPerformer, IActionPerformer, PipelineLabels, Proc, Command, NewCommand(), newNameCollector() (+23 more)

### Community 36 - "IProjectMock GetType Expectations"
Cohesion: 0.28
Nodes (4): IProjectMockGetTypeExpectation, IProjectMockGetTypeResults, mIProjectMockGetType, Type

### Community 37 - "Multiplexed Server Manager"
Cohesion: 0.29
Nodes (7): PrepareConfigFolder, Node, appendToConfig(), AppConfig, isEnvVarName(), marshalEnvExample(), sortEnv()

### Community 38 - "Telegram Version Handler"
Cohesion: 0.08
Nodes (26): EnumGenArg, generalGenArgs, internalConfigGenerator, loadConfigFileGenArgs, enumGenArg, structGenArgs, Environment, InternalConfig (+18 more)

### Community 39 - "Telegram Transport Listener"
Cohesion: 0.35
Nodes (10): buildCLI(), T, run(), Test_GeneratedProjectsCompile(), T, scaffoldComboProject(), stripGeneratedMarker(), Test_GeneratedProject_GoimportsClean() (+2 more)

### Community 42 - "Git Status Diff"
Cohesion: 0.06
Nodes (29): Request, RW, Context, Changes, CommitWithUntrackedAction, gitChangesType, InitGit, InstallHooksAction (+21 more)

### Community 43 - "Build Project Action"
Cohesion: 0.36
Nodes (9): commonProjectTidyPostActions(), commonProjectTidyPreActions(), GetTidyActionsForProject(), Action, goProjectTidyActions(), T, Test_GetTidyActionsForProject_Fast(), Test_GetTidyActionsForProject_UnknownType() (+1 more)

### Community 46 - "RW File Locking Utility"
Cohesion: 0.25
Nodes (3): Mutex, Reader, RW

### Community 47 - "IOMock Error Expectations"
Cohesion: 0.28
Nodes (4): IOMockErrorExpectation, IOMockErrorParams, mIOMockError, RWMutex

### Community 48 - "IProjectMock GetConfig Expectations"
Cohesion: 0.16
Nodes (8): expected grpc config_template.yaml fixture, expected grpc dev.yaml fixture, Config, serviceDockerfileArgs, IProjectMockGetConfigExpectation, IProjectMockGetConfigResults, mIProjectMockGetConfig, GenerateDockerfile()

### Community 52 - "SQL Connection (pattern)"
Cohesion: 0.25
Nodes (6): Adding a new `add` dependency, Architecture, Command flow, Commands, Key packages, Test helpers

### Community 55 - "IOMock Print Expectations"
Cohesion: 0.32
Nodes (3): IOMockPrintExpectation, IOMockPrintParams, mIOMockPrint

### Community 56 - "IOMock Println Expectations"
Cohesion: 0.32
Nodes (3): IOMockPrintlnExpectation, IOMockPrintlnParams, mIOMockPrintln

### Community 61 - "IOMock PrintColored Expectations"
Cohesion: 0.25
Nodes (7): Acceptance Criteria, Context, Do NOT change, Files to Create / Modify, Goal, Notes, Task 001 — Pattern system: embed and copy static template files directly

### Community 70 - "RW Struct Core Methods"
Cohesion: 0.23
Nodes (4): ActionPerformerMock, Duration, Tester, NewActionPerformerMock()

### Community 74 - "IOMock GetInput Builder"
Cohesion: 0.32
Nodes (3): IOMockPrintlnColoredExpectation, IOMockPrintlnColoredParams, mIOMockPrintlnColored

### Community 81 - "CI / CLAUDE.md Dev Commands"
Cohesion: 0.67
Nodes (3): CLAUDE.md Development Commands Section, build-and-test job (branch-push workflow), create-pr job (branch-push workflow)

### Community 209 - "Project Name Substitution"
Cohesion: 0.15
Nodes (12): Config example can be found in internal/config/verv.yaml, Configuration, Configuration structure (all fields are optional), Features:, Installation, SPECIAL VARS, THE LAST, NOT THE LEAST on project creation, TODO (+4 more)

## Knowledge Gaps
- **203 isolated node(s):** `Proc`, `Proc`, `go.vervstack.ru/verv`, `Project`, `opt` (+198 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **175 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Folder` connect `Virtual Folder Tree` to `RsCli Config Loading`, `Multiplexed Server Manager`, `Telegram Version Handler`, `IProject Mock (minimock)`, `Dependency AppendToProject (Sqlite/Env)`, `IProjectMock GetFolder Expectations`, `gRPC Server Transport (pattern)`, `HTTP Server Transport`?**
  _High betweenness centrality (0.168) - this node is a cross-community bridge._
- **Why does `IProject` connect `Go Fmt & Makefile Gen Actions` to `HTTP Server Transport (pattern)`, `gRPC Package Discovery`, `Multiplexed Server Manager`, `RW Struct Core Methods`, `Config Folder Generation`, `gRPC Client/Gateway Interfaces`, `Git Status Diff`, `Embedded Template Registration`, `IProjectMock GetProjectPath Builder`, `IProjectMock GetConfig Expectations`, `TS gRPC-Gateway Fetch Client`, `Go Mod Init Action`, `gRPC Server Transport (pattern)`, `Virtual Folder Tree`, `HTTP Server Transport`?**
  _High betweenness centrality (0.148) - this node is a cross-community bridge._
- **Why does `Color` connect `Terminal Color Parser & IO Stub` to `Generated Project Config Loader`, `IO Mock (minimock)`, `IOMock GetInput Builder`?**
  _High betweenness centrality (0.141) - this node is a cross-community bridge._
- **What connects `Proc`, `Proc`, `go.vervstack.ru/verv` to the rest of the system?**
  _204 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Project Action Pipeline` be split into smaller, more focused modules?**
  _Cohesion score 0.07692307692307693 - nodes in this community are weakly interconnected._
- **Should `gRPC Package Discovery` be split into smaller, more focused modules?**
  _Cohesion score 0.12169312169312169 - nodes in this community are weakly interconnected._
- **Should `Dockerfile Generation & Build Tests` be split into smaller, more focused modules?**
  _Cohesion score 0.08013937282229965 - nodes in this community are weakly interconnected._