package pipeline

import (
	"context"

	"go.redsock.ru/rerrors"
)

func Probe(ctx context.Context, node Node) Status {
	var status Status

	platform, err := node.Version(ctx)
	if err != nil {
		status.Err = rerrors.Wrap(err, "error reading node version")
		return status
	}

	status.Platform = platform

	names, err := node.ListRegistryNames(ctx)
	if err != nil {
		status.Err = rerrors.Wrap(err, "error listing registries")
		return status
	}

	status.IsRegistryEnabled = len(names) > 0

	return status
}
