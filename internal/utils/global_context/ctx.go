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
	return WithTimeout(defaultContextTimeout)
}

func WithTimeout(timeout time.Duration) context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)

	closer.Add(func() error {
		cancel()

		return nil
	})

	return ctx
}
