package deploy

import (
	"context"
	"fmt"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/verv/internal/io"
	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
	"go.vervstack.ru/verv/plugins/deployment/nodes"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

func (p *projectDeploy) runInteractive(ctx context.Context, params pipeline.Params) error {
	store, err := nodes.NewStore(p.wd)
	if err != nil {
		return rerrors.Wrap(err, "error opening nodes store")
	}

	entry, aborted, err := p.chooseNode(ctx, store)
	if err != nil {
		return rerrors.Wrap(err, "error choosing node")
	}

	if aborted {
		return nil
	}

	params.Environment, aborted, err = p.resolveEnvironment(params.Environment)
	if err != nil {
		return rerrors.Wrap(err, "error resolving environment")
	}

	if aborted {
		return nil
	}

	deps := p.newDeps(nodeclient.New(entry.Url, entry.Secret))

	err = pipeline.Run(ctx, deps, params)
	if err != nil {
		return rerrors.Wrap(err, "error running pipeline")
	}

	return nil
}

func (p *projectDeploy) chooseNode(ctx context.Context, store *nodes.Store) (nodes.Entry, bool, error) {
	for {
		entries, err := store.List()
		if err != nil {
			return nodes.Entry{}, false, rerrors.Wrap(err, "error listing nodes")
		}

		hasNodes := len(entries) > 0

		if hasNodes {
			statuses := p.probeWithSpinner(ctx, entries)

			index, isPickAborted, pickErr := p.pickNode(entries, statuses)
			if pickErr != nil {
				return nodes.Entry{}, false, rerrors.Wrap(pickErr, "error picking node")
			}

			if isPickAborted {
				return nodes.Entry{}, true, nil
			}

			if index != addNodeIndex {
				return entries[index], false, nil
			}
		}

		isAddAborted, err := p.addNode(ctx, store)
		if err != nil {
			return nodes.Entry{}, false, rerrors.Wrap(err, "error adding node")
		}

		if isAddAborted {
			return nodes.Entry{}, true, nil
		}
	}
}

func (p *projectDeploy) probeWithSpinner(ctx context.Context, entries []nodes.Entry) []pipeline.Status {
	spinner := io.NewSpinner(p.io)
	spinner.Start("Checking nodes")

	statuses := probeNodes(ctx, entries)

	spinner.Stop(true, fmt.Sprintf("Checked %d node(s)", len(entries)))

	return statuses
}
