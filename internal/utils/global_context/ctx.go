package global_context

import (
	"context"
	"time"

	"go.redsock.ru/toolbox/closer"
)

const (
	defaultContextTimeout = time.Second * 5
)

func Background() context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)

	closer.Add(func() error {
		cancel()

		return nil
	})

	return ctx
}
