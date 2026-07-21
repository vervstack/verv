// Package e2e scaffolds real projects with the compiled verv CLI and checks
// that the generated source actually compiles. Unlike the generator unit
// tests, this exercises the full `project init` / `project add` pipeline
// end-to-end (real go.mod, real `go mod tidy`, real `go build`), which is the
// only way to catch generators that produce code referencing something that
// doesn't exist (missing imports, renamed fields, wrong function arity...).
package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
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

		cmd := exec.Command("go", "build", "-o", out, ".")
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

	cmd := exec.Command(name, args...)
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

// TestGeneratedProjectsCompile scaffolds a fresh project for each dependency
// set below via the real CLI, then runs `go mod tidy` and `go build ./...`
// against the result. A case failing here means `verv project init`/`add`
// produced code that doesn't compile.
func TestGeneratedProjectsCompile(t *testing.T) {
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
		{name: "redis", deps: []string{"redis"}},
		{name: "postgres", deps: []string{"postgres"}},
		{name: "sqlite", deps: []string{"sqlite"}},
		{name: "env", deps: []string{"env"}},
		{name: "telegram", deps: []string{"telegram"}},
		{name: "grpc", deps: []string{"grpc"}},
		{name: "combo", deps: []string{"redis", "postgres", "sqlite", "env", "telegram", "grpc"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			workDir := t.TempDir()
			projName := "sanity_" + tc.name

			run(t, workDir, bin, "project", "init", projName)

			projDir := filepath.Join(workDir, projName)

			if len(tc.deps) > 0 {
				run(t, projDir, bin, append([]string{"project", "add"}, tc.deps...)...)
			}

			run(t, projDir, "go", "mod", "tidy")
			run(t, projDir, "go", "build", "./...")
		})
	}
}
