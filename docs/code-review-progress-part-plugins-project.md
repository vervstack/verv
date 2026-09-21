# Dev-profile conformance audit — `plugins/project/`

Scope: `plugins/project/` in full — `IProject`/`Project` core, `actions/` (the
`ActionPerformer` pipeline: prepare structure → clients → server → Dockerfile →
build → init → fmt → git commit, plus `actions/go_actions/dependencies/`),
`config/`, `go_project/` (including `go_project/patterns/generators/` — the
codegen that produces config structs, app structs, gRPC transport, and
Dockerfiles for every downstream project), and `validators/`.

Methodology matches the sibling audit of `~/verv/Velez`
(`docs/dev-profile-coverage.md` / `docs/code-review-progress.md` there, used
only as a style/depth/tagging reference — not part of this repo). Rule tags:
**L** linter-enforced, **C** custom-lint candidate, **F** language-forced,
**P** deliberate preference. Applicable packs: `00-core.md`, `10-workflow.md`,
`30-go.md`, `32-go-testing.md`, `20-architecture.md`. `31-go-layering.md`
(transport/service/repository/domain) mostly doesn't apply — this is a CLI
tool built around a codegen/action-pipeline shape, not a backend service.

This package's generators are what produced the matreshka-migration
regenerated code recognizable from Velez's `internal/config/` and
`internal/app/` — bugs here have wide blast radius across every downstream
project scaffolded or `tidy`'d with this CLI.

## Fixed (no-brainer, applied this pass)

| # | File | Issue | Fix |
|---|------|-------|-----|
| 1 | `actions/go_actions/init.go` (`InitGoMod.Do`) | Opened `go.mod` in append mode after `go mod init`, never wrote to it, and closed it inside a `defer` that reassigned the *local, non-named* `err` — so any close failure was silently discarded and the function always returned the literal `nil`. The whole open/close block had no observable effect either way. | Removed the dead open/close block entirely; `InitGoMod.Do` now just runs `go mod init` and returns its error. Confirmed via `git log --follow -p` this bug (swallowed close error via a non-named return) has existed unchanged since the earliest `RSI-*`-era commits, pre-`rscli`→`verv` rename, and never manifested because `Close()` on a freshly-opened file essentially never fails. |
| 2 | `actions/go_actions/preparation.go` (`addMissingImplFolders`) | `transportFolder.Add(implFolders...)` added the **entire** `implFolders` slice on every iteration where the current folder was found missing, instead of just the one missing folder (`implF`). Verified via `internal/io/folder/folder.go`'s `Add()` (replaces-by-name at the top level, doesn't merge) that this was mostly inert in practice, because `impl_gen.GenerateImpl` already only returns genuinely-missing packages — but the mechanical bug is real and a no-brainer to fix regardless. | `transportFolder.Add(implF)`. |
| 3 | `actions/go_actions/dependencies/telegram.go` (`Telegram.applyFolder`) | Returned the raw `err` from `containsDependencyFolder` unwrapped, inconsistent with the sibling `applyClient`, which wraps the identical call with `rerrors.Wrap(err, "error finding Dependency path")`. Style/consistency gap against `30-go.md`'s "never construct/return errors ad hoc, wrap consistently" posture. | Wrapped identically to `applyClient`. |
| 4 | `go_project/patterns/generators/app_struct_generators/app.go` (`mergeContentImports`) | Error message started capitalized, "Fatal error: app already imported package … . But dependency requires …" — violates `30-go.md`'s lowercase, non-sentence-case error-message convention (and injects an editorial "Fatal error:" prefix a wrapped error shouldn't carry). | Reworded to lowercase, single-sentence: `"app already imported package X with alias Y, but dependency requires this package to be imported as Z"`. |
| 5 | `go_project/patterns/generators/app_struct_generators/server.go` (`generateServerInitFileAndArgs`) | `server.Name == ""` was exactly the triggering condition for the wrap, so the message always rendered as `server "" doesn't exist in config` — actively misleading (the real problem is a missing required name when multiple servers are configured, not a lookup failure). The wrap also added nothing over the sentinel: `ErrServerMustHaveName` already carries the accurate message ("application contains more than two servers. Names required to be specified"). | Return `ErrServerMustHaveName` directly, no wrap. Confirmed no test or caller depends on the old wrapped text (`grep` for both strings found only the definition site). |
| 6 | `go_project/patterns/generators/server_generators/templates/fs.go.pattern` + `file_server_test.go` | **Compile-breaking generator bug.** `fs.Sub(distFS, {{ .DistPath }})` rendered `{{ .DistPath }}` as a bare, unquoted Go identifier instead of a string literal — `fs.Sub(distFS, dist)` for a simple path (undefined identifier `dist`), or outright invalid syntax for a nested path (`fs.Sub(distFS, web/build)`). `file_server_test.go`'s own `expected` fixtures asserted this broken output as correct, so the bug was fully test-covered as "working." Also fixed the template's `if err != nil { return ... }` body, which wasn't tab-indented (cosmetic; `gofmt`/`goimports` runs over generated files during `init`/`tidy` so this alone wasn't load-bearing, but it's free to fix while touching the line). | Quoted the template value: `fs.Sub(distFS, "{{ .DistPath }}")`; updated both `expected` fixtures in `file_server_test.go` to the corrected, quoted output; fixed the indentation in both the template and the test fixtures. **Scoping note:** `GenerateFileServer` is currently called only from its own test — `grep -rn "GenerateFileServer"` found no production caller anywhere in the codebase (confirmed `preparation.go`'s `PrepareServer.Do` wires GRPC transport only, no FS/file-server generation path). So this bug, while genuine and worth fixing, has not been shipping broken output into any real generated project today — it would have on first use. |

All fixes verified with `go build ./...` (clean) and `go test ./plugins/project/...` (all packages pass). `golangci-lint run ./plugins/project/...` shows 4 pre-existing issues, none touched by or introduced by these edits (1 `lll` line-length hit in `preparation.go`'s untouched `generateMiddlewareFiles`, 3 `nolintlint` hits in `actions/git/git_test.go` for stale `//nolint:mnd` directives) — left alone as out of scope.

## Flagged, not fixed — needs a decision

These are real defects or gaps, but fixing them requires a behavioral/semantic
judgment call beyond "no-brainer" (per `00-core.md`: confidence <95% on a
change → ask with options, don't guess).

### 1. `plugins/project/load.go` — dead sort comparator, unknown config-precedence semantics

```go
var configOrder = map[string]int{
    prodConfigFileName:     1,
    templateConfigFileName: 2,
}
...
configsPaths = append(configsPaths, path.Join(c.ConfigDir, d.Name())) // full path stored
...
sort.Slice(configsPaths, func(i, j int) bool {
    return configOrder[configsPaths[i]] > configOrder[configsPaths[j]] // map keyed by bare filename, looked up by full path — always 0 == 0
})
```

`configOrder`'s keys are bare filenames (`config.yaml`, `config_template.yaml`),
but the comparator looks them up by the **full joined path** — the lookup
always misses on both sides, so the comparator is a no-op. Effective ordering
falls back to whatever `os.ReadDir` returns (alphabetical), which happens to
coincidentally put `config.yaml` before `config_template.yaml` today.

**Why not fixed:** I could not determine from this codebase alone whether
`matreshka.ReadConfig`'s merge semantics mean "later path in the list wins" or
"loses" — flipping a currently-broken-but-accidentally-correct comparator to a
working one could silently invert real config-precedence behavior for every
project built with this CLI. Needs either a design decision or a
`matreshka.ReadConfig` read to confirm merge direction before touching.

### 2. `actions/go_actions/dependencies/link_service/grpc_discovery/discovery.go` — brittle `GOPATH` resolution

```go
var modFolderPath = os.Getenv("GOPATH") + "/pkg/mod/"
```

If `GOPATH` isn't explicitly exported (common on modern Go toolchains, which
compute a default internally without setting the env var), this silently
resolves to `/pkg/mod/` from filesystem root instead of the real module cache
— a latent portability bug that would only surface as a confusing "module not
found"-type failure in some environments.

**Why not fixed:** a correct fix means either shelling out to `go env GOPATH`
or replicating Go's default-GOPATH computation — a design choice (extra
subprocess call vs. duplicated logic) beyond a mechanical edit.

### 3. `go_project/patterns/generators/config_generators/env_config_generator/env_config_generator.go` — `[]int`-valued enums silently dropped

```go
switch vals := enumVal.(type) {
case []string:
    ... // generates enum code
case []int:
    // empty — no-op, nothing generated, no error
default:
    return rerrors.New("error generating enums for config value. Unsupported enum type %T. Expected String slice", enumVal)
}
```

An `[]int` enum value produces neither codegen nor an error — silently
different from a genuinely unsupported type (which correctly errors via
`default`). Looks like either a forgotten case or a deliberately deferred
feature; can't tell which from the code alone.

**Why not fixed:** needs a decision — is int-enum codegen intentionally out of
scope (in which case it should error like any other unsupported type), or is
it meant to work and was never implemented?

### 4. `actions/go_actions/preparation.go` — `addMissingImplFolders`'s existence check is folder-level, not file-level

Beyond the mechanical `implFolders...` → `implF` fix applied above, the
function's existence check only looks at whether a folder **name** already
exists under `transportFolder` — not whether the specific file
(`impl.go`) inside it exists. If a package folder already exists but is
missing only `impl.go` (partial state), `impl_gen.GenerateImpl` would
correctly generate a fresh `impl.go` for it, but `addMissingImplFolders` would
see the folder name already present and skip attaching it — silently
dropping the freshly-generated file. Combined with `folder.Add()`'s
replace-by-name-not-merge semantics, this is also a latent data-loss path if
the (now-fixed) add-all bug were ever to fire against such a partial folder.

**Why not fixed:** correctly resolving this needs a decision on intended
semantics — should `addMissingImplFolders` merge at the file level, or is
whole-folder skip actually the intended behavior and the partial-folder case
is simply not expected to occur? Flagging only.

## Read, clean, no issues

Everything else in scope was read line-by-line with no findings:

- `i_project.go`, `project.go`, `marker.go`/`marker_test.go`, `create.go`,
  `config/config.go`, `validators/name_validation.go`.
- `actions/actions.go`, `actions/init.go`, `actions/tidy.go` and their tests —
  thorough table-driven `t.Parallel()` suites using `require`.
- `actions/git/*.go` (all files) and `git_test.go`.
- `actions/ci_cd/github_workflow.go`.
- `actions/go_actions/actions.go`, `build.go`/`build_test.go`,
  `go_mod.go` (`GoFmt`, `RunGoTidyAction`, `RunMakeGenAction`,
  `UpdateAllPackages`), `marker.go`/`marker_test.go`, `git_hooks.go`,
  `renamer/renaming.go`.
- `actions/go_actions/dependencies/*.go` (`dependencies.go`, `interfaces.go`,
  `client_pg.go`+test, `client_redis.go`, `client_sqlite.go`,
  `env_variable.go`, `sql.go`) — idempotent-by-design, well tested.
- `actions/go_actions/dependencies/link_service/client_grpc.go`,
  `grpc_discovery/grpc.go`, `fetcher.go`, `name_filter.go` (the
  module-path-escaping scheme in `FilterPackageName` is correct, verified by
  manual trace).
- `go_project/patterns/*.go` (`ci_cd.go`, `common_files.go`,
  `common_imports.go`, `conn_names.go`, `go_files.go`, `grpc.go`, `hooks.go`)
  — static `//go:embed` constants, nothing to break.
- `go_project/patterns/generators/config.go`, `internal_config.go`.
- `go_project/patterns/generators/config_generators/*` (`gen_data_sources.go`,
  `gen_servers.go` + test, `gen_skeleton.go` + test, `gen_config_struct.go` +
  test, `common.go`, `templates.go`) — confirmed `gen_servers.go`'s
  `ErrorMessage: "error parsing servers to config"` is already lowercase from
  an earlier pass this session; left untouched.
- `go_project/patterns/generators/app_struct_generators/common.go`,
  `data_sources.go`+test, `templates.go`, `app_files_test.go` — golden-file
  snapshot tests, cover Postgres-exclusion and custom.go-preservation
  behavior.
- `go_project/patterns/generators/clients_generators/*`,
  `dockerfile_generator/*`, `grpc_api_generator/*`, `main_generators/*`,
  `middleware_generators/*` — all read, clean, straightforward
  template-execution wrappers.
- `go_project/patterns/generators/transport_generators/*`
  (`transport_generator.go`, `templates.go`, `transport_generator_test.go`) —
  clean; simple `text/template`/`//go:embed` execution wrappers with direct
  equality and `Contains` assertions.
- `internal/io/folder/folder.go` (read for context on `.Add()` semantics,
  needed to assess flagged item 4 above — not itself in audit scope).

Notable: `preparation_test.go` documents (via
`Test_PrepareServer_AttachesTransportFolder`) a previously-fixed real bug
where `PrepareServer.Do` used to build but never attach the transport folder
— evidence this codebase is actively maintained and self-correcting, not
simply undertested.

## Rule-profile notes surfaced by this audit

No profile-rule text changes are proposed. One observation worth recording in
the ledger if useful going forward: this package is a strong example of
`00-core.md`'s "never edit generated files, change the source and
regenerate" principle in action from the *producer* side — `file_server_test.go`
baking in a broken template's literal output as its `expected` fixture is the
exact failure mode that rule exists to prevent one level up (a generator's
own test asserting broken output as ground truth). Worth keeping in mind for
any future codegen work in this codebase: a generator test's golden output is
only as trustworthy as a manual read of what it actually asserts.
