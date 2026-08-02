package e2e

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"golang.org/x/tools/imports"

	"go.vervstack.ru/verv/plugins/project/actions/go_actions/dependencies"
)

// scaffoldComboProject builds the CLI, scaffolds a fresh project with every
// `add` dependency wired in (mirroring compile_test.go's "combo" case), and
// returns its directory.
func scaffoldComboProject(t *testing.T) string {
	t.Helper()

	bin := buildCLI(t)

	workDir := t.TempDir()
	projName := "lint_combo"

	run(t, workDir, bin, "init", projName, "--fast")

	projDir := filepath.Join(workDir, projName)

	run(t, projDir, bin, "add", "--fast",
		dependencies.DependencyNameRedis,
		dependencies.DependencyNamePostgres,
		dependencies.DependencyNameSqlite,
		dependencies.DependencyEnvVariable,
		dependencies.DependencyNameTelegram,
	)

	return projDir
}

// walkGoFiles calls fn for every *.go file under root, skipping .git/vendor.
func walkGoFiles(t *testing.T, root string, fn func(path string) error) {
	t.Helper()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "vendor" {
				return filepath.SkipDir
			}

			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		return fn(path)
	})
	if err != nil {
		t.Fatalf("walking %s: %v", root, err)
	}
}

// Test_GeneratedProject_GoimportsClean scaffolds a combo project via the real
// CLI, then runs every generated *.go file through the same
// golang.org/x/tools/imports formatter the production GoFmt action uses,
// failing if any file's on-disk content differs from the formatted output.
// This is deliberately independent of golangci-lint (and its generated-file
// skip) - it's the check that would have caught the original bug where the
// tidy pipeline only ran `go fmt`, never goimports, and left a broken import
// group behind on disk.
func Test_GeneratedProject_GoimportsClean(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("scaffolds a real project; skipped in -short mode")
	}

	projDir := scaffoldComboProject(t)

	var dirty []string

	walkGoFiles(t, projDir, func(path string) error {
		src, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		formatted, err := imports.Process(path, src, nil)
		if err != nil {
			return fmt.Errorf("running goimports on %s: %w", path, err)
		}

		if !bytes.Equal(src, formatted) {
			rel, relErr := filepath.Rel(projDir, path)
			if relErr != nil {
				rel = path
			}

			dirty = append(dirty, rel)
		}

		return nil
	})

	if len(dirty) > 0 {
		t.Fatalf("goimports found unformatted files in generated project:\n%s", dirty)
	}
}

// generatedFileMarker matches the "Code generated ... DO NOT EDIT." header
// golangci-lint's generated-file detection looks for. Stripping it lets us
// lint the underlying code for real instead of trusting it's clean because
// the marker hid it.
var generatedFileMarker = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

// Test_GeneratedProject_GolangciLintClean copies a freshly scaffolded combo
// project into a temp dir with the "Code generated ... DO NOT EDIT." marker
// stripped from every *.go file, then runs golangci-lint against that copy.
// The marker is legitimate, deliberate Go convention that this codebase
// itself relies on elsewhere (see custom.go's different, intentionally-kept
// marker) - this test just makes sure that whatever it's hiding is
// genuinely clean, so future generator regressions can't coast on the skip.
func Test_GeneratedProject_GolangciLintClean(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("scaffolds a real project and shells out to golangci-lint; skipped in -short mode")
	}

	lintBin, err := exec.LookPath("golangci-lint")
	if err != nil {
		t.Skip("golangci-lint not found on PATH; skipping")
	}

	projDir := scaffoldComboProject(t)

	// t.TempDir() on macOS lives under /var/folders/..., itself a symlink to
	// /private/var/folders/.... golangci-lint resolves its own working
	// directory via the physical (already-resolved) path, then computes
	// issue-printing/generated-file-filter paths relative to the string we
	// hand it in cmd.Dir; if that string still has the unresolved /var
	// symlink, the two disagree and every relative path golangci-lint builds
	// from them escapes into nonexistent directories - manifesting as bogus
	// "no such file" warnings and, worse, nolint directives that appear
	// "unused" simply because the filter reading them via a broken path
	// silently no-ops. Resolving symlinks up front keeps cmd.Dir consistent
	// with what golangci-lint sees internally.
	copyDir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("resolving symlinks for copy dir: %v", err)
	}

	err = filepath.WalkDir(projDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(projDir, path)
		if err != nil {
			return err
		}

		dst := filepath.Join(copyDir, rel)

		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}

			return os.MkdirAll(dst, 0o755)
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		if filepath.Ext(path) == ".go" {
			src = stripGeneratedMarker(src)

			if _, err := parser.ParseFile(token.NewFileSet(), path, src, parser.AllErrors); err != nil {
				return fmt.Errorf("stripped file %s no longer parses: %w", rel, err)
			}
		}

		return os.WriteFile(dst, src, 0o644)
	})
	if err != nil {
		t.Fatalf("copying project to %s: %v", copyDir, err)
	}

	cmd := exec.CommandContext(t.Context(), lintBin, "run")
	cmd.Dir = copyDir
	// Every run scaffolds byte-identical generated source into a fresh temp
	// dir, so golangci-lint's own result cache (~/Library/Caches/golangci-lint
	// on macOS) can key a hit off a previous run's content and replay
	// diagnostics that still point at that older, by-now-deleted temp path -
	// producing "no such file" noise and false results unrelated to anything
	// this test is actually checking. Pointing the cache at a fresh dir per
	// run keeps every invocation a real, from-scratch analysis.
	cmd.Env = append(os.Environ(), "GOLANGCI_LINT_CACHE="+t.TempDir())

	out, err := cmd.CombinedOutput()
	if err != nil {
		// golangci-lint exits non-zero when it finds issues (or hits a config
		// error). Its own deprecation/config warnings are expected and go to
		// the same combined output on a clean run, so the exit code - not
		// output emptiness - is what actually signals a real finding here.
		t.Fatalf("golangci-lint reported issues in generated project (marker stripped):\n%s", out)
	}
}

// stripGeneratedMarker removes any line matching generatedFileMarker from
// src, line by line, preserving every other line unchanged. The blank line
// immediately following a stripped marker is dropped too, so removing the
// marker never leaves a leading blank line behind - that's an artifact of
// this stripping process, not something gofmt would ever consider clean, and
// would otherwise produce a false-positive "not gofmt'd" finding that has
// nothing to do with the actual generator output being linted.
func stripGeneratedMarker(src []byte) []byte {
	lines := bytes.Split(src, []byte("\n"))

	kept := lines[:0]

	skipNextBlank := false

	for _, line := range lines {
		if generatedFileMarker.Match(line) {
			skipNextBlank = true

			continue
		}

		if skipNextBlank {
			skipNextBlank = false

			if len(bytes.TrimSpace(line)) == 0 {
				continue
			}
		}

		kept = append(kept, line)
	}

	return bytes.Join(kept, []byte("\n"))
}
