package processor

import (
	"go.vervstack.ru/verv/internal/config"
	"go.vervstack.ru/verv/internal/io"
)

func WithIo(io io.IO) opt {
	return func(p *Processor) {
		p.IO = io
	}
}

func WithWd(wd string) opt {
	return func(p *Processor) {
		p.WD = wd
	}
}

func WithConfig(cfg *config.VervConfig) opt {
	return func(p *Processor) {
		p.RscliConfig = cfg
	}
}
