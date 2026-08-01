package processor

import (
	"os"

	"github.com/spf13/cobra"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/plugins/project"
)

const (
	PathFlag = "path"
	FastFlag = "fast"
)

// Processor - represents a single process of execution.
// e.g. verv tidy - calls a cmd/project/tidy Processor and executes it
// Contains all basic necessary information and primitives for CLI utility
type Processor struct {
	IO         io.IO
	WD         string
	VervConfig *config.VervConfig
}

type opt func(p *Processor)

func New(opts ...opt) Processor {
	p := Processor{}

	for _, o := range opts {
		o(&p)
	}

	if p.IO == nil {
		p.IO = io.StdIO{}
	}

	if p.WD == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		p.WD = wd
	}

	if p.VervConfig == nil {
		p.VervConfig = config.GetConfig()
	}

	return p
}

func (p *Processor) LoadProject(cmd *cobra.Command) (proj *project.Project, err error) {
	var pathToProject string

	if cmd != nil {
		pathToProject = cmd.Flag(PathFlag).Value.String()
	}

	pathToProject = toolbox.Coalesce(pathToProject, p.WD)

	proj, err = project.LoadProject(pathToProject, p.VervConfig)
	if err != nil {
		return nil, rerrors.Wrap(err, "error loading project")
	}

	return proj, nil
}
