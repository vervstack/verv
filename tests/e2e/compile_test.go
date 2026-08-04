// Package e2e scaffolds real projects with the compiled verv CLI and checks
// that the generated source actually compiles. Unlike the generator unit
// tests, this exercises the full `init` / `add` pipeline
// end-to-end (real go.mod, real `go mod tidy`, real `go build`), which is the
// only way to catch generators that produce code referencing something that
// doesn't exist (missing imports, renamed fields, wrong function arity...).
package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"go.vervstack.ru/verv/plugins/project/actions/go_actions/dependencies"
)

var (
	buildOnce sync.Once
	binPath   string
	errBuild  error
)

// buildCLI compiles the verv CLI once and shares the binary across subtests.
func buildCLI(t *testing.T) string {
	t.Helper()

	buildOnce.Do(func() {
		repoRoot, err := filepath.Abs("../..")
		if err != nil {
			errBuild = fmt.Errorf("resolving repo root: %w", err)

			return
		}

		out := filepath.Join(t.TempDir(), "verv-e2e")
		if runtime.GOOS == "windows" {
			out += ".exe"
		}

		cmd := exec.CommandContext(t.Context(), "go", "build", "-o", out, ".")

		cmd.Dir = repoRoot

		var output []byte

		output, errBuild = cmd.CombinedOutput()
		if errBuild != nil {
			errBuild = fmt.Errorf("building verv CLI: %w\n%s", errBuild, output)

			return
		}

		binPath = out
	})

	if errBuild != nil {
		t.Fatal(errBuild)
	}

	return binPath
}

// run executes name with args in dir, failing the test with the combined
// output if it exits non-zero.
func run(t *testing.T, dir, name string, args ...string) {
	t.Helper()

	cmd := exec.CommandContext(t.Context(), name, args...)

	cmd.Dir = dir

	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=verv-e2e",
		"GIT_AUTHOR_EMAIL=verv-e2e@localhost",
		"GIT_COMMITTER_NAME=verv-e2e",
		"GIT_COMMITTER_EMAIL=verv-e2e@localhost",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s\n(dir=%s)\nerror: %v\noutput:\n%s",
			name, strings.Join(args, " "), dir, err, out)
	}
}

// addLocalMatreshkaReplace points a scaffolded project's go.mod at the local
// Matreshka checkout instead of its latest published tag.
//
// TEMPORARY: this exists only because go.vervstack.ru/matreshka's
// ReadConfig/WithConfigPaths/WithConfigBytes API (used by the generated
// internal/config/load.go and skeleton.go) is not released yet — it only
// exists in the sibling Matreshka working tree. Every scaffolded project
// resolves its own go.mod against the module proxy, so without this replace
// every case in this file fails to build with "undefined: matreshka.ReadConfig"
// etc. Once that Matreshka change is tagged and released, delete this
// function and its call sites; scaffolded projects will resolve the real
// published version on their own.
func addLocalMatreshkaReplace(t *testing.T, projDir string) {
	t.Helper()

	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("resolving repo root: %v", err)
	}

	matreshkaDir, err := filepath.Abs(filepath.Join(repoRoot, "..", "Matreshka"))
	if err != nil {
		t.Fatalf("resolving matreshka dir: %v", err)
	}

	run(t, projDir, "go", "mod", "edit", "-replace", "go.vervstack.ru/matreshka="+matreshkaDir)
	run(t, projDir, "go", "mod", "tidy")
}

// Test_GeneratedProjectsCompile scaffolds a fresh project for each dependency
// set below via the real CLI, then runs `go build ./...` against the result.
// A case failing here means `verv project init`/`add` produced code that
// doesn't compile. `go mod tidy` is not run here — `init`/`add` already leave
// a tidy go.mod/go.sum as the last step of their own action pipeline.
func Test_GeneratedProjectsCompile(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("scaffolds and compiles real projects; skipped in -short mode")
	}

	bin := buildCLI(t)

	cases := []struct {
		name string
		deps []string
	}{
		{name: "base"},
		{
			name: dependencies.DependencyNameRedis,
			deps: []string{dependencies.DependencyNameRedis},
		},
		{
			name: dependencies.DependencyNamePostgres,
			deps: []string{dependencies.DependencyNamePostgres},
		},
		{
			name: dependencies.DependencyNameSqlite,
			deps: []string{dependencies.DependencyNameSqlite},
		},
		{
			name: dependencies.DependencyEnvVariable,
			deps: []string{dependencies.DependencyEnvVariable},
		},
		{
			name: dependencies.DependencyNameTelegram,
			deps: []string{dependencies.DependencyNameTelegram},
		},
		{
			name: "combo",
			deps: []string{
				dependencies.DependencyNameRedis,
				dependencies.DependencyNamePostgres,
				dependencies.DependencyNameSqlite,
				dependencies.DependencyEnvVariable,
				dependencies.DependencyNameTelegram,
			}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			workDir := t.TempDir()
			projName := "sanity_" + tc.name

			run(t, workDir, bin, "init", projName, "--fast")

			projDir := filepath.Join(workDir, projName)

			if len(tc.deps) > 0 {
				run(t, projDir, bin, append([]string{"add", "--fast"}, tc.deps...)...)
			}

			addLocalMatreshkaReplace(t, projDir)

			run(t, projDir, "go", "build", "./...")
		})
	}
}

// Test_GeneratedProject_RunsWithoutConfigFile scaffolds a dependency-free
// project, deletes its generated config/config.yaml and
// config/config_template.yaml, and runs the compiled service binary with
// only the env vars from its generated config/.env.example set — nothing
// else on disk. This is the scenario the embedded config skeleton exists
// for: internal/config/skeleton.go (//go:embed skeleton.yaml) feeds
// matreshka.ReadConfig via WithConfigBytes as a fallback source, so even
// with zero config files present at runtime, env vars alone must be enough
// for the app to load its config and start up cleanly.
func Test_GeneratedProject_RunsWithoutConfigFile(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("scaffolds, compiles and runs a real project; skipped in -short mode")
	}

	bin := buildCLI(t)

	workDir := t.TempDir()
	projName := "sanity_no_config_file"

	run(t, workDir, bin, "init", projName, "--fast")

	projDir := filepath.Join(workDir, projName)

	addLocalMatreshkaReplace(t, projDir)

	envExampleContent, err := os.ReadFile(filepath.Join(projDir, "config", ".env.example"))
	if err != nil {
		t.Fatalf("reading generated .env.example: %v", err)
	}

	envVars := parseEnvExample(envExampleContent)

	err = os.Remove(filepath.Join(projDir, "config", "config.yaml"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("removing config.yaml: %v", err)
	}

	err = os.Remove(filepath.Join(projDir, "config", "config_template.yaml"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("removing config_template.yaml: %v", err)
	}

	serviceBin := filepath.Join(workDir, "service")
	if runtime.GOOS == "windows" {
		serviceBin += ".exe"
	}

	run(t, projDir, "go", "build", "-o", serviceBin, "./cmd/...")

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, serviceBin)

	cmd.Dir = projDir

	cmd.Env = append([]string{"PATH=" + os.Getenv("PATH")}, envVars...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("service exited with error running with no config file on disk, env-vars-only config:\n%v\noutput:\n%s",
			err, out)
	}
}

// parseEnvExample extracts KEY=VALUE env entries from a generated
// .env.example file's content, skipping blank lines and comments.
func parseEnvExample(content []byte) []string {
	var vars []string

	lines := strings.SplitSeq(string(content), "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		vars = append(vars, line)
	}

	return vars
}
