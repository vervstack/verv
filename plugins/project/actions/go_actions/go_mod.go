package go_actions

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/cmd"
	vervconfig "go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/internal/utils/bins/makefile"
	"go.vervstack.ru/verv/plugins/project"
)

const (
	goBin = "go"

	formattedFileMode fs.FileMode = 0o644
)

type GoFmt struct{}

// Do walks every *.go file under the project path and runs it through
// golang.org/x/tools/imports (the goimports library, in-process), which is a
// strict superset of gofmt: it fixes formatting like gofmt AND fixes import
// grouping/ordering, which bare `go fmt` never touches.
func (a GoFmt) Do(p project.IProject) error {
	root := p.GetProjectPath()

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

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		src, readErr := os.ReadFile(path)
		if readErr != nil {
			return rerrors.Wrap(readErr, "error reading go file")
		}

		formatted, procErr := imports.Process(path, src, nil)
		if procErr != nil {
			return rerrors.Wrap(procErr, "error running goimports on "+path)
		}

		if bytes.Equal(src, formatted) {
			return nil
		}

		return os.WriteFile(path, formatted, formattedFileMode)
	})
	if err != nil {
		return rerrors.Wrap(err, "error formatting project")
	}

	return nil
}
func (a GoFmt) NameInAction() string {
	return "Performing project fix up"
}

type RunGoTidyAction struct{}

func (a RunGoTidyAction) Do(p project.IProject) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"mod", "tidy"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return rerrors.Wrap(err, "error executing go mod tidy")
	}

	err = GoFmt{}.Do(p)
	if err != nil {
		return rerrors.Wrap(err, "error formatting project")
	}

	return nil
}
func (a RunGoTidyAction) NameInAction() string {
	return "Cleaning up the project"
}

type RunMakeGenAction struct {
	C  *vervconfig.VervConfig
	IO io.IO
}

func (a RunMakeGenAction) Do(p project.IProject) error {
	if len(p.GetConfig().Servers) == 0 {
		return nil
	}

	err := makefile.Install()
	if err != nil {
		return rerrors.Wrap(err, "error installing makefile")
	}

	return nil
}
func (a RunMakeGenAction) NameInAction() string {
	return "Running `make gen`"
}

type UpdateAllPackages struct{}

func (a UpdateAllPackages) Do(p project.IProject) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    goBin,
		Args:    []string{"get", "-u", "all"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}

func (a UpdateAllPackages) NameInAction() string {
	return "Updating packages to latest version"
}
